# TECH.md - 技術仕様書（現状反映版）

「kotti-he-oide」プロジェクトの技術的な詳細を、リポジトリの実装・設定に基づいて最新の状態へ更新したドキュメントです。以前の版に含まれていた未実装・計画中の記述は、実装状況に応じて明確に区別しました。

## 目次

1. 技術スタック（実装ベース）
2. アーキテクチャ
3. 開発環境・ツール
4. データベース設計（schema.sql 準拠）
5. API 仕様（実装準拠）
6. gRPC マイクロサービス
7. インフラストラクチャ（Terraform）
8. CI/CD（GitHub Actions）
9. セキュリティ（実装/計画）
10. パフォーマンス（現状）
11. 開発ガイドライン
12. 特殊機能詳細（Web Push / カメラ）
13. 運用・監視
14. 今後の技術的課題
15. よく使用するコマンド
16. デフォルトポート / メタ情報

---

## 1. 技術スタック（実装ベース）

### バックエンド（server/）
- 言語: Go 1.25.1
- HTTP: 標準ライブラリ（net/http）
- アーキテクチャ: DDD + レイヤード（`internal/{domain,application,infrastructure,interfaces}`）
- データ永続化: 現状はインメモリ実装（`internal/infrastructure/persistence/memory_*`）。PostgreSQL 用 `server/schema.sql` は用意済みだが、アプリからは未接続。
- gRPC クライアント: 画像認識マイクロサービスへ接続（`/api/ml/*` プロキシ）。
- 主要ライブラリ（実際に使用）:
  - `github.com/SherClockHolmes/webpush-go`（Web Push 送信 / VAPID キー生成）
  - `google.golang.org/grpc`（gRPC クライアント）
  - `google.golang.org/protobuf`（gRPC スタブ）
- ツール/依存（未使用または将来利用予定）:
  - `github.com/golang-jwt/jwt/v5`（go.mod に間接依存として存在。現状コード未使用）
  - `github.com/google/wire`（tool 指定のみ。現状コード未使用の DI）

### マイクロサービス（services/image_recognition/）
- 言語: Python 3.12+
- 通信: gRPC（`grpcio`）
- 画像処理: OpenCV（`opencv-python-headless`）
- 数値計算: NumPy
- S3 連携: boto3（参照画像を起動時に取得可能。環境変数で制御）
- パッケージ管理: uv（`uv sync` / `uv run`）

### フロントエンド（frontend/）
- フレームワーク: Next.js 15.5.2（App Router）
- ランタイム: React 19.1.0
- 言語: TypeScript 5.x
- パッケージマネージャ: pnpm
- ビルド: Turbopack
- 機能:
  - カメラフィルター（Canvas 2D + ノイズエンジン、`/camera-filters`）
  - Web Push 設定 UI（`/notifications`）
  - Service Worker（`public/sw.js`。通知受信に特化、オフラインキャッシュは未実装）

### インフラストラクチャ
- クラウド: AWS（ECS Fargate / ALB / ECR / S3 / RDS 等）
- IaC: Terraform（AWS Provider ~> 5.0, Terraform >= 1.6）
- ログ: CloudWatch Logs
- コード生成: Buf v2（Go/Python スタブを生成）

---

## 2. アーキテクチャ

- バックエンドは DDD + レイヤード。依存方向は Interfaces → Application → Domain ← Infrastructure。
- 現状の永続化はメモリ実装。将来的に PostgreSQL 実装（`schema.sql`）へ差し替え予定。
- 画像認識は Python gRPC サービスへ委譲し、Go 側で HTTP→gRPC のプロキシを提供。

---

## 3. 開発環境・ツール

### 必須ツール
- Go 1.25.1+
- Node.js 20.x+ / pnpm（Corepack 推奨）
- Python 3.12+ / uv
- Docker / Terraform >= 1.6 / AWS CLI
- Buf CLI v2

### 開発支援
- golangci-lint（Go） / ESLint（Next.js） / ruff・mypy・pytest（Python）

### コード品質チェック（必須）

フロントエンド
```bash
cd frontend
pnpm lint
```

バックエンド
```bash
cd server
make fmt && make lint
```

マイクロサービス（Python）
```bash
cd services/image_recognition
uv run ruff check
uv run ruff format
uv run mypy .
```

---

## 4. データベース設計（server/schema.sql 準拠、現状は未接続）

アプリはインメモリ実装で動作中ですが、PostgreSQL 用のスキーマは以下を含みます。

- `push_subscriptions`：購読情報（`endpoint` UNIQUE、`is_valid` 部分インデックス）
- `notification_templates`：通知テンプレート
- `notification_prefs`：ユーザー別通知設定
- `push_jobs`：非同期ジョブ（`job_status` enum: pending/sending/succeeded/failed/cancelled）
- `push_logs`：配信ログ（HTTP ステータス/ヘッダ/エラー）

代表的なインデックス:
- `push_jobs(status)`, `push_jobs(schedule_at)`, `push_jobs(idempotency_key)` など
- `push_logs(job_id)`, `push_logs(subscription_id)`, `push_logs(created_at)`

---

## 5. API 仕様（実装準拠）

### ヘルスチェック
```
GET /api/healthz
200 OK
{
  "status": "ok",
  "message": "Server is running"
}
```

### ユーザー管理
```
GET    /api/users
POST   /api/users
GET    /api/users/{id}
DELETE /api/users/{id}
```

戻り値例：`GET /api/users`
```json
{
  "users": [
    { "id": 1, "name": "Alice", "email": "alice@example.com", "created_at": "...", "updated_at": "..." }
  ],
  "count": 1
}
```

### Web Push
```
GET    /api/push/vapid-public-key      # VAPID 公開鍵取得（ランタイム生成）
POST   /api/push/subscribe             # 購読登録（メモリ保存）
DELETE /api/push/subscriptions/{id}    # 購読解除
POST   /api/push/send                  # 通知送信ジョブ作成（201 Created）
POST   /api/push/send/batch            # バッチ送信ジョブ作成（201 Created）
# POST /api/push/click                 # クリック計測（未実装）
```

戻り値例：`POST /api/push/send`
```json
{ "jobId": "...", "success": true, "message": "queued" }
```

### 機械学習 API（gRPC プロキシ）
```
GET  /api/ml/hello?name=world
POST /api/ml/recognize   # multipart/form-data（image または file）/ 生バイナリ
```

戻り値例：`POST /api/ml/recognize`
```json
{ "is_match": false, "similarity_score": 0.42, "error_message": "", "backend": "127.0.0.1:50051" }
```

---

## 6. gRPC マイクロサービス（services/image_recognition）

### 技術スタック
- Python 3.12+ / gRPC（`grpcio`）/ OpenCV / NumPy / boto3 / uv

### サービス定義（Protocol Buffers）
`schema/proto/image_recognition/v1/image_recognition.proto`
```protobuf
service ImageRecognitionService {
  rpc Hello (HelloRequest) returns (HelloReply);
  rpc RecognizeImage(RecognizeImageRequest) returns (RecognizeImageResponse);
  rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
}
```

### 主要機能
- 画像類似度判定（しきい値指定可。未指定は `DEFAULT_SIMILARITY_THRESHOLD`）
- 複数フォーマット対応（JPEG/PNG/WebP/BMP）
- S3 から参照画像を起動時ロード（`S3_BUCKET_NAME` 等）

### 開発コマンド
```bash
cd services/image_recognition
uv sync                         # 依存
uv run python -m app.server     # 起動
uv run pytest                   # テスト
uv run ruff check && uv run mypy app
```

### Buf 設定とコード生成
`buf.gen.yaml` により Go/Python のスタブを生成：
```bash
buf generate   # ルートで実行（Go: server/internal/gen, Python: services/.../app/gen）
```

---

## 7. インフラストラクチャ（infra/）

### 構成（主なファイル）
- `vpc.tf` / `alb.tf` / `ecs_cluster.tf`
- `ecs_services_web.tf` / `ecs_services_api.tf` / `ecs_services_microservice.tf`
- `rds.tf` / `s3.tf` / `ecr.tf` / `security.tf`
- `providers.tf` / `versions.tf` / `variables.tf` / `outputs.tf`

### コンピュート
- ECS Fargate：
  - `${var.name_prefix}-web`（フロント）
  - `${var.name_prefix}-api`（バックエンド）
  - `${var.name_prefix}-microservice`（gRPC、Cloud Map に登録）
- ALB：`/` → Web、`/api/*` → API へルーティング
- API → マイクロサービスは Cloud Map のプライベート DNS（例: `microservice.${var.name_prefix}.local`）で接続

### データ/ログ
- RDS（PostgreSQL）定義あり（現状アプリ未接続）
- S3（参照画像など）
- CloudWatch Logs（各サービス）

---

## 8. CI/CD（.github/workflows/）

- `frontend-deploy.yml`：`frontend/**` の変更でビルド→ECR→ECS 更新
- `server-deploy.yml`：`server/**` または `schema/proto/**` 変更で Go / Buf 生成→ECR→ECS 更新
- `microservice-deploy.yml`：`services/**` または `schema/proto/**` 変更で Python / Buf 生成→ECR→ECS 更新
- `infra-plan.yml`：`infra/**` 変更で Terraform fmt/validate/plan

---

## 9. セキュリティ（実装/計画）

### 認証・認可
- VAPID（実装済）: Web Push の署名・認証（鍵はランタイム生成・未保存）
- JWT（計画中）: コード内参照は無し
- HTTPS: 本番は ALB で TLS 終端

### データ保護
- 機密情報: GitHub Secrets / 環境変数
- RDS: 暗号化（Terraform 定義）
- 通信: TLS 1.2+

### セキュリティヘッダー
- 必要に応じてフロント/ALB 側で付与（現状サーバー側での明示付与は無し）

---

## 10. パフォーマンス（現状）

### フロントエンド
- Turbopack、React 19、Canvas 2D 最適化（解像度/フレーム制御）

### バックエンド
- Go（net/http）
- 非同期 Push 送信ワーカー（メモリキュー）
- DB 接続プール（pgx 等）は未使用

---

## 11. 開発ガイドライン

### コーディング規約
- Go: `gofmt`、公開識別子は `PascalCase`、パッケージ名は小文字、`panic` 禁止
- TypeScript/React: 2 スペース、`PascalCase.tsx`、kebab-case ルート、CSS Modules

### テスト
- Go: テーブル駆動の `_test.go`
- Python: pytest / mypy / ruff
- フロント: 重要機能にテスト推奨（E2E は将来 Playwright）

### コミット/PR
- Conventional Commits 推奨、CI Green 必須（`pnpm lint` / `make test`）

---

## 12. 特殊機能詳細

### Web Push 通知
- RFC 8030/8291/8292 準拠、VAPID 認証、メッセージ暗号化
- TTL / Urgency / Topic に対応（Web Push ヘッダ）
- クリック計測エンドポイントは未実装（Service Worker からは呼出し想定あり）

### カメラフィルター
- Canvas 2D によるリアルタイム画像処理（5 種類 + ノイズ合成）
- `src/app/camera-filters` に実装

### PWA
- Web App Manifest あり
- Service Worker は通知用途のみ（オフラインキャッシュ無し）

---

## 13. 運用・監視

- メトリクス例: 応答時間/エラー率、CPU/メモリ、通知配信/クリック
- ログ: CloudWatch Logs（JSON 構造化を推奨）
- アラート: 高負荷/エラー率等（CloudWatch）

---

## 14. 今後の技術的課題

- RDS への接続・リポジトリ移行（メモリ→PostgreSQL）
- 認証（JWT/OIDC）
- クリック計測 API の実装
- キャッシュ（例: Redis）/ 分散トレーシング
- E2E テスト導入（Playwright）

---

## 15. よく使用するコマンド

### バックエンド（server/）
```bash
make run         # ローカル起動
make run-cmd     # cmd/server 版で起動
make build       # ビルド
make test        # テスト
make deps        # 依存取得
make fmt && make lint
make proto       # buf generate（Go/Python 両方のスタブ生成）
```

### フロントエンド（frontend/）
```bash
pnpm dev         # 開発（http://localhost:3000）
pnpm build && pnpm start
pnpm lint
```

### マイクロサービス（services/image_recognition/）
```bash
uv sync
uv run python -m app.server
uv run pytest && uv run ruff check && uv run mypy app
```

### Protocol Buffers（ルート）
```bash
buf generate
buf lint
buf format
```

### インフラ（infra/）
```bash
terraform init
terraform plan -var-file=terraform.tfvars
terraform apply
```

---

## 16. デフォルトポート / メタ情報

- バックエンド: 8080（`PORT` で変更可）
- フロントエンド: 3000（Next.js デフォルト）
- 画像認識 gRPC: 50051（`GRPC_PORT`）
- ALB: 80/443

最終更新: 2025-09-07  
バージョン: 1.2.0

