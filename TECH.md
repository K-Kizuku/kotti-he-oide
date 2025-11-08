# TECH.md - 技術仕様書

「赤煉瓦文化館 〜こっちにおいで〜」プロジェクトの技術的な詳細を、実装状況に基づいて記載したドキュメントです。

## 目次

1. 技術スタック
2. アーキテクチャ
3. 開発環境・ツール
4. データベース設計
5. API 仕様
6. マイクロサービス
7. インフラストラクチャ
8. CI/CD
9. セキュリティ
10. パフォーマンス
11. 開発ガイドライン
12. 特殊機能詳細
13. 運用・監視
14. よく使用するコマンド
15. デフォルトポート

---

## 1. 技術スタック

### バックエンド（server/）
- **言語**: Go 1.25.1
- **HTTP**: 標準ライブラリ（net/http）
- **アーキテクチャ**: DDD + レイヤード（`internal/{domain,application,infrastructure,interfaces}`）
- **データ永続化**: MySQL（RDS）
  - リポジトリ実装: `internal/infrastructure/persistence/mysql_*_repository.go`
  - マイグレーション: `migrations/000001_create_game_tables.up.sql`、`migrations/000002_create_push_tables.up.sql`
- **主要ライブラリ**:
  - `github.com/go-sql-driver/mysql` v1.8.1（MySQL ドライバ）
  - `github.com/google/uuid` v1.6.0（セッションID生成）
  - `github.com/SherClockHolmes/webpush-go` v1.4.0（Web Push 送信 / VAPID）
  - `google.golang.org/grpc` v1.65.0（gRPC クライアント）
  - `google.golang.org/protobuf` v1.34.2（Protocol Buffers）
- **間接依存**:
  - `github.com/golang-jwt/jwt/v5` v5.2.1（webpush-go の依存）
  - `golang.org/x/crypto` v0.31.0（暗号化処理）

### マイクロサービス

#### 画像認識サービス（services/image_recognition/）
- **言語**: Python 3.12+
- **通信**: gRPC（`grpcio` >= 1.62）
- **画像処理**: OpenCV（`opencv-python-headless` >= 4.9）
- **数値計算**: NumPy >= 1.26
- **S3 連携**: boto3 >= 1.34（参照画像の取得）
- **パッケージ管理**: uv

#### AI エージェントサービス（services/ai_agent/）
- **言語**: Python 3.12+
- **フレームワーク**: FastAPI >= 0.109.0
- **AI エージェント**: Strands Agents >= 0.1.0、Bedrock AgentCore >= 0.1.0
- **LLM**: Amazon Bedrock Claude 4 Sonnet
- **AWS SDK**: boto3 >= 1.34.0
- **パッケージ管理**: uv
- **機能**:
  - クイズ生成（プレイヤー回答に基づく個別化）
  - 回答バリデーション
  - 対話生成（1942年設定維持）
  - メッセージ検証
  - AgentCore Memory による過去回答の保存・検索

### フロントエンド（frontend/）
- **フレームワーク**: Next.js 15.5.2（App Router）
- **ランタイム**: React 19.1.0
- **言語**: TypeScript 5.x
- **パッケージマネージャ**: pnpm
- **ビルド**: Turbopack
- **主要機能**:
  - ゲームフロー（S0〜S9）
  - カメラ撮影とホラー演出オーバーレイ
  - 音声再生（VOICEVOX生成音声）
  - セッション管理
  - タイマー機能（S6の7分制限）

### インフラストラクチャ
- **クラウド**: AWS
- **コンピュート**: ECS Fargate
- **ロードバランサー**: Application Load Balancer（ALB）
- **データベース**: RDS MySQL
- **ストレージ**: S3
- **コンテナレジストリ**: ECR
- **ネットワーク**: VPC（パブリック/プライベートサブネット）
- **サービスディスカバリ**: AWS Cloud Map
- **IaC**: Terraform >= 1.6（AWS Provider ~> 5.0）
- **ログ**: CloudWatch Logs
- **コード生成**: Buf v2（Protocol Buffers）

---

## 2. アーキテクチャ

### 全体構成
```
[フロントエンド (Next.js)]
         ↓ HTTP
    [ALB (AWS)]
    /          \
   /            \
[API Server]  [Web Server]
  (Go)         (Next.js)
   |
   |-- MySQL (RDS)
   |-- S3
   |-- 画像認識サービス (gRPC)
   |-- AI エージェントサービス (HTTP)
   |-- VOICEVOX (EC2)
```

### バックエンドアーキテクチャ（DDD + レイヤード）

**依存方向**: Interfaces → Application → Domain ← Infrastructure

#### レイヤー構成
- **Interfaces 層** (`internal/interfaces/http/`)
  - HTTPハンドラー、DTO、ミドルウェア
  - 外部からのリクエストを受け付け、Application層へ委譲

- **Application 層** (`internal/application/usecase/`)
  - ユースケース実装（ビジネスフロー制御）
  - トランザクション境界の管理

- **Domain 層** (`internal/domain/`)
  - エンティティ、値オブジェクト、リポジトリインターフェース
  - ビジネスロジックの中核
  - インフラストラクチャに依存しない

- **Infrastructure 層** (`internal/infrastructure/`)
  - リポジトリ実装（MySQL）
  - 外部サービス連携（VOICEVOX、gRPC）
  - データベース接続管理

### マイクロサービス連携
- **画像認識**: Go → gRPC → Python（類似度計算）
- **AI エージェント**: Go → HTTP → Python/FastAPI（クイズ生成等）
- **音声生成**: Go → HTTP → VOICEVOX（EC2）

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

## 4. データベース設計

### ゲームテーブル（migrations/000001_create_game_tables.up.sql）

#### sessions
セッション管理テーブル
- `session_id` CHAR(36) PRIMARY KEY（UUID v4）
- `current_scene` VARCHAR(10)（S0〜S9）
- `s6_started_at` DATETIME（S6開始時刻、7分制限用）
- `created_at` DATETIME
- `expires_at` DATETIME（60分TTL）
- インデックス: `expires_at`, `created_at`

#### session_answers
S4の内省質問回答テーブル
- `id` INT AUTO_INCREMENT PRIMARY KEY
- `session_id` CHAR(36)（外部キー）
- `question_id` INT（質問番号）
- `answer_text` TEXT
- `answered_at` DATETIME
- ユニーク制約: `(session_id, question_id)`

#### session_s6_progress
S6の場所探索進捗テーブル
- `id` INT AUTO_INCREMENT PRIMARY KEY
- `session_id` CHAR(36)（外部キー）
- `place_id` VARCHAR(50)（場所ID）
- `verified` BOOLEAN（場所到達確認）
- `verified_by` VARCHAR(20)（確認方法: camera/manual）
- `quiz_id` CHAR(36)（クイズID）
- `answered` BOOLEAN（クイズ回答済み）
- `correct` BOOLEAN（正解フラグ）
- `verified_at` DATETIME
- ユニーク制約: `(session_id, place_id)`

#### quiz_questions
生成されたクイズテーブル
- `quiz_id` CHAR(36) PRIMARY KEY
- `session_id` CHAR(36)（外部キー）
- `place_id` VARCHAR(50)
- `question_text` TEXT
- `option_1` 〜 `option_4` TEXT（4択）
- `answer_index` INT（正解の選択肢番号）
- `created_at` DATETIME

#### player_messages
S9で刻まれたメッセージテーブル
- `id` INT AUTO_INCREMENT PRIMARY KEY
- `session_id` CHAR(36)
- `place_id` VARCHAR(50)（メッセージを刻んだ場所）
- `message_text` TEXT
- `created_at` DATETIME

### Web Pushテーブル（migrations/000002_create_push_tables.up.sql）

#### push_subscriptions
購読情報テーブル
- `id` BIGSERIAL PRIMARY KEY
- `user_id` BIGINT（外部キー、オプション）
- `endpoint` TEXT UNIQUE（Push Service エンドポイント）
- `p256dh` TEXT（ブラウザ公開鍵）
- `auth` TEXT（認証シークレット）
- `ua` TEXT（User-Agent）
- `expiration_time` TIMESTAMPTZ
- `is_valid` BOOLEAN
- `created_at`, `updated_at` TIMESTAMPTZ

#### push_jobs
非同期送信ジョブテーブル
- `id` BIGSERIAL PRIMARY KEY
- `idempotency_key` TEXT UNIQUE
- `template_key` TEXT
- `user_id` BIGINT
- `topic` TEXT（Web Push Topic ヘッダ）
- `urgency` TEXT（very-low/low/normal/high）
- `ttl_seconds` INT
- `payload` JSONB
- `schedule_at` TIMESTAMPTZ
- `status` job_status ENUM（pending/sending/succeeded/failed/cancelled）
- `retry_count` INT
- `last_error` TEXT
- `created_at`, `updated_at` TIMESTAMPTZ

#### push_logs
配信ログテーブル
- `id` BIGSERIAL PRIMARY KEY
- `job_id` BIGINT（外部キー）
- `subscription_id` BIGINT（外部キー）
- `response_status` INT（HTTP ステータス）
- `response_headers` JSONB
- `error` TEXT
- `created_at` TIMESTAMPTZ

#### notification_templates
通知テンプレートテーブル
- `id` BIGSERIAL PRIMARY KEY
- `key` TEXT UNIQUE
- `title` TEXT
- `body` TEXT
- `url` TEXT
- `icon` TEXT
- `data` JSONB
- `created_at` TIMESTAMPTZ

#### notification_prefs
ユーザー通知設定テーブル
- `user_id` BIGINT PRIMARY KEY
- `enabled` BOOLEAN
- `topics` JSONB
- `quiet_hours` JSONB

---

## 5. API 仕様

### ヘルスチェック
```
GET /api/healthz
```
レスポンス:
```json
{
  "status": "ok",
  "message": "Server is running"
}
```

### セッション管理
```
POST /api/session
```
セッション作成。60分TTLのセッションIDを発行。

```
GET /api/session/{session_id}
```
セッション情報取得。

```
POST /api/session/{session_id}/s6/start
```
S6（7分制限探索パート）開始。`s6_started_at`を記録。

### S4 回答（内省質問）
```
POST /api/session/{session_id}/answers
```
リクエスト:
```json
{
  "question_id": 1,
  "answer_text": "子どもの頃の夢は..."
}
```

```
GET /api/session/{session_id}/answers
```
セッションの全回答を取得。

### S6 進捗管理（場所探索）
```
POST /api/session/{session_id}/s6/initialize
```
5つの場所の進捗レコードを初期化。

```
POST /api/session/{session_id}/s6/verify-location
```
リクエスト:
```json
{
  "place_id": "main_hall",
  "image": "base64_encoded_image"
}
```
画像認識サービスで場所到達を確認。

```
GET /api/session/{session_id}/s6/quiz/{place_id}
```
その場所のクイズを取得（AI エージェントで生成）。

```
POST /api/session/{session_id}/s6/answer
```
リクエスト:
```json
{
  "place_id": "main_hall",
  "quiz_id": "uuid",
  "selected_index": 2
}
```
クイズ回答を送信。正解判定を返す。

```
GET /api/session/{session_id}/s6/progress
```
5つの場所の進捗状況を取得。

### メッセージ管理
```
POST /api/session/{session_id}/message
```
リクエスト:
```json
{
  "place_id": "main_hall",
  "message_text": "私の人生の目標は..."
}
```
S9で刻むメッセージを保存。

```
GET /api/messages
```
過去のプレイヤーが刻んだメッセージを取得（匿名化）。

### Web Push 通知
```
GET /api/push/vapid-public-key
```
VAPID公開鍵を取得（フロントエンドの購読登録に使用）。

```
POST /api/push/subscribe
```
リクエスト:
```json
{
  "endpoint": "https://...",
  "keys": {
    "p256dh": "...",
    "auth": "..."
  }
}
```
購読情報を登録。

```
DELETE /api/push/subscriptions/{subscription_id}
```
購読解除。

```
POST /api/push/send/{session_id}
```
リクエスト:
```json
{
  "title": "通知タイトル",
  "body": "通知本文",
  "url": "/game/s5"
}
```
指定セッションに通知を送信。

### 音声生成
```
POST /api/voice/generate
```
リクエスト:
```json
{
  "text": "こんにちは",
  "speaker_id": 46
}
```
VOICEVOX で音声を生成し、S3にアップロード。音声URLを返す。

---

## 6. マイクロサービス

### 画像認識サービス（services/image_recognition/）

#### 技術スタック
- Python 3.12+ / gRPC / OpenCV / NumPy / boto3 / uv

#### サービス定義（Protocol Buffers）
`schema/proto/image_recognition/v1/image_recognition.proto`
```protobuf
service ImageRecognitionService {
  rpc Hello (HelloRequest) returns (HelloReply);
  rpc RecognizeImage(RecognizeImageRequest) returns (RecognizeImageResponse);
  rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
}
```

#### 主要機能
- 画像類似度判定（OpenCV特徴量マッチング）
- しきい値指定可能（デフォルト: `DEFAULT_SIMILARITY_THRESHOLD`）
- 複数フォーマット対応（JPEG/PNG/WebP/BMP）
- S3から参照画像を起動時ロード

#### 開発コマンド
```bash
cd services/image_recognition
uv sync                         # 依存関係インストール
uv run python -m app.server     # gRPCサーバー起動（ポート50051）
uv run pytest                   # テスト実行
uv run ruff check && uv run mypy app  # Lint & 型チェック
```

### AI エージェントサービス（services/ai_agent/）

#### 技術スタック
- Python 3.12+ / FastAPI / Strands Agents / AWS Bedrock AgentCore / boto3 / uv

#### エンドポイント
```
GET  /ping                      # ヘルスチェック
GET  /health                    # 詳細ヘルスチェック
POST /invocations               # メインエントリーポイント
```

#### 操作（operation）
`POST /invocations` のリクエストボディで指定：

1. **generate_quiz**: クイズ生成
   - プレイヤーのS4回答を元に個別化された4択クイズを生成
   - 過去プレイヤーの回答をダミー選択肢に活用

2. **validate_response**: 回答バリデーション
   - 無効回答の検出とフィードバック生成

3. **generate_dialogue**: 対話生成
   - 1942年設定を維持した担当者との対話テキスト生成

4. **verify_message**: メッセージ検証
   - S7での回答再入力の正確性判定

5. **store_answer**: 回答保存
   - AgentCore Memory への回答保存

#### 開発コマンド
```bash
cd services/ai_agent
uv sync                         # 依存関係インストール
make dev                        # 開発モード起動（ホットリロード）
make run                        # 本番モード起動
make lint                       # Lint実行
```

### Protocol Buffers コード生成
`buf.gen.yaml` により Go/Python のスタブを生成：
```bash
buf generate   # ルートで実行
# Go: server/internal/gen
# Python: services/image_recognition/app/gen
```

---

## 7. インフラストラクチャ（infra/）

### Terraform 構成
- `vpc.tf`: VPC、サブネット、インターネットゲートウェイ
- `alb.tf`: Application Load Balancer、ターゲットグループ、リスナールール
- `ecs_cluster.tf`: ECS クラスター
- `ecs_services_web.tf`: フロントエンドサービス（Next.js）
- `ecs_services_api.tf`: バックエンドサービス（Go）
- `ecs_services_microservice.tf`: マイクロサービス（画像認識、AI エージェント）
- `rds.tf`: RDS MySQL インスタンス
- `s3.tf`: S3 バケット（画像、音声ファイル）
- `ecr.tf`: ECR リポジトリ（web、api、microservice）
- `security.tf`: セキュリティグループ、IAM ロール
- `providers.tf`: AWS プロバイダー設定
- `versions.tf`: Terraform バージョン制約
- `variables.tf`: 入力変数
- `outputs.tf`: 出力値

### コンピュート
- **ECS Fargate サービス**:
  - `${var.name_prefix}-web`: フロントエンド（Next.js）
  - `${var.name_prefix}-api`: バックエンド（Go）
  - `${var.name_prefix}-image-recognition`: 画像認識サービス（Python gRPC）
  - `${var.name_prefix}-ai-agent`: AI エージェントサービス（Python FastAPI）

- **ALB ルーティング**:
  - `/` → Web サービス
  - `/api/*` → API サービス

- **サービスディスカバリ**:
  - AWS Cloud Map でマイクロサービスをプライベート DNS に登録
  - 例: `image-recognition.${var.name_prefix}.local:50051`

### データストア
- **RDS MySQL**: ゲームデータ、Web Push データ
- **S3**: 参照画像、撮影画像、生成音声ファイル

### ログ・監視
- **CloudWatch Logs**: 各サービスのログ
- **CloudWatch Metrics**: CPU、メモリ、リクエスト数等

### デプロイコマンド
```bash
cd infra
terraform init
terraform plan -var-file=terraform.tfvars
terraform apply -var-file=terraform.tfvars
```

---

## 8. CI/CD

### GitHub Actions ワークフロー

#### frontend-deploy.yml
- **トリガー**: `frontend/**` の変更
- **処理**:
  1. Next.js ビルド（Turbopack）
  2. Docker イメージビルド
  3. ECR へプッシュ
  4. ECS サービス更新

#### server-deploy.yml
- **トリガー**: `server/**` または `schema/proto/**` の変更
- **処理**:
  1. Buf による Protocol Buffers コード生成
  2. Go ビルド
  3. Docker イメージビルド
  4. ECR へプッシュ
  5. ECS サービス更新

#### microservice-deploy.yml（画像認識、AI エージェント）
- **トリガー**: `services/**` または `schema/proto/**` の変更
- **処理**:
  1. Buf による Protocol Buffers コード生成（画像認識のみ）
  2. Python 依存関係インストール（uv）
  3. Docker イメージビルド
  4. ECR へプッシュ
  5. ECS サービス更新

#### infra-plan.yml
- **トリガー**: `infra/**` の変更
- **処理**:
  1. Terraform fmt チェック
  2. Terraform validate
  3. Terraform plan（PR コメントに結果を投稿）

### デプロイフロー
```
コード変更 → GitHub Push → GitHub Actions
  ↓
Docker Build → ECR Push → ECS Service Update
  ↓
新しいタスク起動 → ヘルスチェック → 旧タスク停止
```

---

## 9. セキュリティ

### 認証・認可
- **VAPID**: Web Push の署名・認証（RFC 8292）
  - 鍵はランタイム生成（`service.NewVAPIDService()`）
  - 公開鍵はフロントエンドに配信
- **セッション管理**: UUID v4 によるセッションID
  - 60分TTL、有効期限切れセッションは自動削除
- **HTTPS**: ALB で TLS 終端

### データ保護
- **機密情報管理**: GitHub Secrets、環境変数
- **RDS 暗号化**: Terraform で有効化
- **通信暗号化**: TLS 1.2+
- **S3**: プライベートバケット、署名付きURL

### CORS
- 開発環境: `Access-Control-Allow-Origin: *`
- 本番環境: 特定オリジンに制限推奨

### 入力バリデーション
- AI エージェントサービス: `InputValidator` による入力検証
- コンテンツフィルタリング: `ContentFilter` による不適切コンテンツ検出

---

## 10. パフォーマンス

### フロントエンド
- **ビルド**: Turbopack による高速ビルド
- **レンダリング**: React 19 の最適化機能
- **カメラ処理**: Canvas 2D、解像度・フレームレート制御
- **音声**: HTML5 Audio API、プリロード

### バックエンド
- **HTTP サーバー**: Go 標準ライブラリ（net/http）
- **DB 接続プール**: MySQL driver による接続プール管理
  - MaxOpenConns: 25
  - MaxIdleConns: 5
  - ConnMaxLifetime: 5分
- **非同期処理**: Web Push 送信ワーカー

### マイクロサービス
- **画像認識**: OpenCV による高速特徴量マッチング
- **AI エージェント**: 
  - Bedrock Claude 4 Sonnet（高速推論）
  - AgentCore Memory によるキャッシュ活用
  - タイムアウト: 15秒（AgentCore Runtime）

### スケーリング
- **ECS Fargate**: 水平スケーリング可能
- **RDS**: 垂直スケーリング、リードレプリカ追加可能
- **ALB**: 自動スケーリング

---

## 11. 開発ガイドライン

### コーディング規約

#### Go
- フォーマット: `gofmt`（`make fmt`）
- 命名:
  - エクスポート: `PascalCase`
  - プライベート: `camelCase`
  - パッケージ名: 小文字、単語1つ
- エラーハンドリング: `panic` 禁止、カスタムエラー型使用
- DDD 原則: ビジネスロジックはドメイン層に配置

#### TypeScript/React
- インデント: 2スペース
- コンポーネント: `PascalCase.tsx`
- ルート: kebab-case（`/game/s0`）
- スタイル: CSS Modules（`*.module.css`）

#### Python
- フォーマット: ruff（`uv run ruff format`）
- 型チェック: mypy（`uv run mypy app`）
- Lint: ruff（`uv run ruff check`）

### テストポリシー

**このプロジェクトはPoC（Proof of Concept）です。**

- **テストコードの作成は不要**
- 開発速度とプロトタイピングを最優先
- ユニットテスト、統合テスト、E2Eテストの実装は求められない

### コード品質チェック（必須）

**すべてのコード変更時に実行し、エラーがないことを確認：**

#### フロントエンド
```bash
cd frontend
pnpm lint    # ESLintチェック（エラーがあれば修正必須）
```

#### バックエンド
```bash
cd server
make fmt     # コードフォーマット
make lint    # golangci-lintチェック（エラーがあれば修正必須）
```

#### マイクロサービス（Python）
```bash
cd services/image_recognition  # または services/ai_agent
uv run ruff check
uv run ruff format
uv run mypy .
```

**注意**: lintエラーが残っている状態でのコミット・プッシュは禁止

### コミット/PR
- Conventional Commits 推奨
- CI Green 必須
- lint エラーゼロ必須

---

## 12. 特殊機能詳細

### Web Push 通知
- **標準準拠**: RFC 8030（Web Push Protocol）、RFC 8291（メッセージ暗号化）、RFC 8292（VAPID）
- **VAPID 認証**: 
  - 秘密鍵・公開鍵ペアをランタイム生成
  - JWT による署名
- **メッセージ暗号化**: ブラウザ公開鍵（p256dh）と認証シークレット（auth）を使用
- **ヘッダー対応**:
  - TTL（Time To Live）
  - Urgency（very-low/low/normal/high）
  - Topic（通知の置き換え）
- **非同期送信**: ジョブキューによるバッチ処理

### カメラ機能
- **撮影**: `getUserMedia` API
- **ホラー演出オーバーレイ**: Canvas 2D による画像合成
  - 影の追加
  - ノイズエフェクト
  - グリッチエフェクト
- **画像認識連携**: 撮影画像を gRPC サービスに送信して場所判定

### 音声生成
- **VOICEVOX**: EC2 上で動作
- **話者**: 青山龍星（しっとり）
- **フロー**:
  1. Go サーバーから VOICEVOX API にリクエスト
  2. 生成された音声を S3 にアップロード
  3. 署名付き URL をフロントエンドに返す
  4. フロントエンドで HTML5 Audio API で再生

### AI エージェント
- **クイズ生成**: プレイヤーの S4 回答を元に個別化
- **対話生成**: 1942年設定を維持した担当者との対話
- **メッセージ検証**: S7 での回答再入力の正確性判定
- **Memory 活用**: AgentCore Memory で過去回答を保存・検索

### PWA
- **Web App Manifest**: `frontend/public/manifest.json`
- **Service Worker**: `frontend/public/sw.js`（通知受信専用）
- **オフラインキャッシュ**: 未実装

---

## 13. 運用・監視

### メトリクス
- **アプリケーション**: 応答時間、エラー率、リクエスト数
- **インフラ**: CPU使用率、メモリ使用率、ネットワークトラフィック
- **Web Push**: 配信成功率、クリック率
- **AI エージェント**: 推論時間、トークン使用量

### ログ
- **CloudWatch Logs**: 各サービスのログを集約
- **構造化ログ**: JSON形式推奨
- **ログレベル**: ERROR、WARN、INFO、DEBUG

### アラート
- **高負荷**: CPU/メモリ使用率が閾値超過
- **エラー率**: 5xx エラーが増加
- **レイテンシ**: 応答時間が閾値超過
- **ヘルスチェック失敗**: サービスダウン検知

---

## 14. よく使用するコマンド

### バックエンド（server/）
```bash
# 開発
make run          # 開発モードでサーバーを実行
make dev          # 自動リロードで実行（airが必要）

# ビルド & テスト
make build        # bin/serverにバイナリをビルド
make test         # Goテストを実行
make deps         # 依存関係をインストール/更新

# コード品質
make fmt          # Goコードをフォーマット
make lint         # リンターを実行（golangci-lintが必要）
make clean        # ビルド成果物をクリーン

# Protocol Buffers
make proto        # buf generate（Go/Pythonスタブ生成）
```

### フロントエンド（frontend/）
```bash
# 開発
pnpm dev          # Turbopackで開発サーバーを開始（http://localhost:3000）
pnpm build        # Turbopackで本番用ビルド
pnpm start        # 本番サーバーを開始
pnpm lint         # ESLintを実行
```

### マイクロサービス

#### 画像認識（services/image_recognition/）
```bash
uv sync                         # 依存関係インストール
uv run python -m app.server     # gRPCサーバー起動（ポート50051）
uv run pytest                   # テスト実行
uv run ruff check               # Lint実行
uv run ruff format              # フォーマット
uv run mypy app                 # 型チェック
```

#### AI エージェント（services/ai_agent/）
```bash
make setup        # 依存関係インストール
make dev          # 開発モード起動（ホットリロード）
make run          # 本番モード起動
make lint         # Lint実行
make test         # テスト実行
```

### Protocol Buffers（ルート）
```bash
buf generate      # Go/Pythonスタブ生成
buf lint          # Lintチェック
buf format        # フォーマット
```

### インフラストラクチャ（infra/）
```bash
# Terraform操作
terraform init      # Terraformを初期化
terraform plan      # インフラストラクチャ変更をプレビュー
terraform apply     # インフラストラクチャ変更を適用
terraform destroy   # インフラストラクチャリソースを削除

# ECRへのイメージプッシュ（例：東京リージョン）
aws ecr get-login-password --region ap-northeast-1 | \
  docker login --username AWS --password-stdin <ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com

docker build -t api:latest ./server
docker tag api:latest <ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com/api:latest
docker push <ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com/api:latest
```

### CI/CD（GitHub Actions）
```bash
# 自動デプロイトリガー
git push origin main  # mainブランチへのpushで自動デプロイ

# フロントエンド: frontend/**の変更でfrontend-deploy.ymlが実行
# バックエンド: server/**の変更でserver-deploy.ymlが実行
# マイクロサービス: services/**の変更でmicroservice-deploy.ymlが実行
```

---

## 15. デフォルトポート

- **バックエンド**: 8080（`PORT`環境変数で設定可能）
- **フロントエンド**: 3000（Next.jsデフォルト）
- **画像認識 gRPC**: 50051（`GRPC_PORT`環境変数で設定可能）
- **AI エージェント**: 8080（FastAPIデフォルト）
- **ALB**: 80/443（HTTP/HTTPS）

---

最終更新: 2025-11-08  
バージョン: 2.0.0

