# TECH.md - 技術仕様書

「kotti-he-oide」プロジェクトの技術的な詳細をまとめたドキュメントです。

## 目次

1. [技術スタック](#技術スタック)
2. [アーキテクチャ](#アーキテクチャ)
3. [開発環境・ツール](#開発環境ツール)
4. [データベース設計](#データベース設計)
5. [API仕様](#api仕様)
6. [インフラストラクチャ](#インフラストラクチャ)
7. [CI/CD](#cicd)
8. [セキュリティ](#セキュリティ)
9. [パフォーマンス](#パフォーマンス)
10. [開発ガイドライン](#開発ガイドライン)

---

## 技術スタック

### バックエンド
- **言語**: Go 1.25.1
- **フレームワーク**: 標準ライブラリHTTPサーバー (net/http)
- **アーキテクチャ**: ドメイン駆動設計（DDD）+ レイヤードアーキテクチャ
- **データベース**: PostgreSQL (RDS)
- **主要ライブラリ**:
  - `github.com/SherClockHolmes/webpush-go` v1.4.0 (Web Push通知)
  - `github.com/golang-jwt/jwt/v5` v5.2.1 (JWT認証)
  - `github.com/jackc/pgx/v5` v5.7.5 (PostgreSQLドライバー)
  - `github.com/google/wire` v0.7.0 (依存性注入)

### フロントエンド
- **フレームワーク**: Next.js 15.5.2 (App Router)
- **ランタイム**: React 19.1.0
- **言語**: TypeScript 5.x
- **パッケージマネージャー**: pnpm
- **ビルドツール**: Turbopack
- **特殊機能**:
  - Canvas 2D API (カメラフィルター)
  - getUserMedia API (カメラアクセス)
  - Web Push API (プッシュ通知)
  - Service Worker (バックグラウンド処理)

### インフラストラクチャ
- **クラウド**: AWS
- **コンテナ**: Docker (マルチステージビルド)
- **オーケストレーション**: ECS Fargate
- **ロードバランサー**: Application Load Balancer (ALB)
- **データベース**: RDS PostgreSQL
- **ストレージ**: S3
- **コンテナレジストリ**: ECR
- **IaC**: Terraform ~> 5.0
- **ログ**: CloudWatch Logs
- **CI/CD**: GitHub Actions

---

## アーキテクチャ

### システム全体構成

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   Database      │
│   (Next.js)     │◄──►│   (Go API)      │◄──►│   (PostgreSQL)  │
│   Port: 3000    │    │   Port: 8080    │    │   Port: 5432    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   ALB           │    │   ECS Fargate   │    │   RDS           │
│   (80/443)      │    │   (Containers)  │    │   (Multi-AZ)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### バックエンドアーキテクチャ (DDD + レイヤード)

```
┌─────────────────────────────────────────────────────────────┐
│                    Interfaces Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   HTTP      │  │    DTO      │  │    Middleware       │ │
│  │  Handlers   │  │  Objects    │  │                     │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Application Layer                         │
│  ┌─────────────┐  ┌─────────────────────────────────────┐   │
│  │  Use Cases  │  │        Application Services         │   │
│  │             │  │                                     │   │
│  └─────────────┘  └─────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Domain Layer                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Models    │  │ Value       │  │    Domain           │ │
│  │ (Entities)  │  │ Objects     │  │   Services          │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
│  ┌─────────────────────────────────────────────────────┐   │
│  │            Repository Interfaces                   │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              ▲
                              │
┌─────────────────────────────────────────────────────────────┐
│                Infrastructure Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │ Persistence │  │   Config    │  │   External APIs     │ │
│  │             │  │             │  │                     │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 依存関係フロー
- **Interfaces** → **Application** → **Domain** ← **Infrastructure**
- 依存関係は内側に向かう（クリーンアーキテクチャの原則）
- Infrastructure層がDomain層のインターフェースを実装

---

## 開発環境・ツール

### 必須ツール
- **Go**: 1.25.1以上
- **Node.js**: 20.x以上
- **pnpm**: 最新版 (corepack enable)
- **Docker**: 最新版
- **Terraform**: ~> 5.0
- **AWS CLI**: 最新版

### 開発支援ツール
- **golangci-lint**: Goコード静的解析
- **air**: Goホットリロード
- **ESLint**: TypeScript/React静的解析
- **Turbopack**: Next.jsビルドツール

### エディタ設定
- **VSCode**: 推奨エディタ
- **Go拡張**: Go言語サポート
- **TypeScript拡張**: TypeScript言語サポート
- **Terraform拡張**: HCL構文サポート

### コード品質チェック（必須）

**すべてのコード変更時に以下を実行：**

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

**重要**: lintエラーが残っている状態でのコミット・プッシュは禁止

---

## データベース設計

### PostgreSQL スキーマ

#### 主要テーブル

**users** - ユーザー管理
```sql
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  email TEXT UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

**push_subscriptions** - Web Push購読管理
```sql
CREATE TABLE push_subscriptions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  endpoint TEXT NOT NULL UNIQUE,
  p256dh TEXT NOT NULL,
  auth TEXT NOT NULL,
  ua TEXT,
  expiration_time TIMESTAMPTZ,
  is_valid BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

**push_jobs** - 非同期通知ジョブ
```sql
CREATE TYPE job_status AS ENUM ('pending','sending','succeeded','failed','cancelled');

CREATE TABLE push_jobs (
  id BIGSERIAL PRIMARY KEY,
  idempotency_key TEXT UNIQUE,
  template_key TEXT REFERENCES notification_templates(key),
  user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  topic TEXT,
  urgency TEXT CHECK (urgency IN ('very-low','low','normal','high')),
  ttl_seconds INT CHECK (ttl_seconds >= 0),
  payload JSONB NOT NULL,
  schedule_at TIMESTAMPTZ,
  status job_status NOT NULL DEFAULT 'pending',
  retry_count INT NOT NULL DEFAULT 0,
  last_error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### インデックス戦略
- `push_subscriptions.endpoint` (UNIQUE)
- `push_subscriptions.user_id`
- `push_subscriptions.is_valid` (WHERE is_valid = true)
- `push_jobs.status`
- `push_jobs.schedule_at`

---

## API仕様

### 基本エンドポイント

#### ヘルスチェック
```
GET /health
Response: 200 OK
```

#### ユーザー管理
```
GET    /api/users           # 全ユーザー取得
POST   /api/users           # ユーザー作成
GET    /api/users/{id}      # 特定ユーザー取得
DELETE /api/users/{id}      # ユーザー削除
```

#### Web Push通知 (計画中)
```
GET    /api/push/vapid-public-key      # VAPID公開鍵取得
POST   /api/push/subscribe             # プッシュ購読登録
DELETE /api/push/subscriptions/{id}    # 購読解除
POST   /api/push/send                  # 通知送信（管理用）
POST   /api/push/send/batch            # バッチ通知送信
POST   /api/push/click                 # 通知クリック計測
```

### フロントエンドルート
```
/                    # ホームページ
/camera-filters      # カメラフィルターデモ
/notifications       # 通知設定ページ（計画中）
```

### レスポンス形式
```json
{
  "success": true,
  "data": {},
  "message": "Success",
  "timestamp": "2025-01-06T12:00:00Z"
}
```

---

## インフラストラクチャ

### AWS構成

#### ネットワーク
- **VPC**: カスタムVPC
- **サブネット**: パブリック/プライベート構成
- **セキュリティグループ**: 最小権限の原則

#### コンピュート
- **ECS Cluster**: Fargateクラスター
- **ECS Services**: 
  - `poc-api` (バックエンド)
  - `poc-web` (フロントエンド)
- **ALB**: パスベースルーティング
  - `/` → Web Service
  - `/api/*` → API Service

#### データ・ストレージ
- **RDS**: PostgreSQL (Multi-AZ)
- **S3**: 静的アセット用
- **ECR**: コンテナイメージ保存

#### 監視・ログ
- **CloudWatch Logs**: アプリケーションログ
- **CloudWatch Metrics**: システムメトリクス

### Terraformファイル構成
```
infra/
├── providers.tf      # AWSプロバイダー設定
├── variables.tf      # 入力変数定義
├── outputs.tf        # 出力値定義
├── vpc.tf           # VPC・ネットワーク
├── security.tf      # セキュリティグループ・IAM
├── alb.tf           # Application Load Balancer
├── ecs_cluster.tf   # ECSクラスター
├── ecs_services_*.tf # ECSサービス定義
├── rds.tf           # RDSデータベース
├── s3.tf            # S3バケット
└── ecr.tf           # ECRリポジトリ
```

---

## CI/CD

### GitHub Actions ワークフロー

#### フロントエンドデプロイ (`.github/workflows/frontend-deploy.yml`)
```yaml
trigger: push to main, changes in frontend/**
steps:
  1. Node.js 20 setup
  2. pnpm install
  3. pnpm lint (必須チェック)
  4. Docker build
  5. ECR push
  6. ECS service update
```

#### バックエンドデプロイ (`.github/workflows/server-deploy.yml`)
```yaml
trigger: push to main, changes in server/**
steps:
  1. Go 1.25.1 setup
  2. go mod download
  3. make test
  4. make fmt (必須チェック)
  5. Docker build
  6. ECR push
  7. ECS service update
```

### デプロイフロー
1. **コード変更** → GitHub push
2. **CI実行** → lint/test/build
3. **イメージビルド** → Docker multi-stage build
4. **ECRプッシュ** → コンテナレジストリ更新
5. **ECS更新** → サービス自動デプロイ

---

## セキュリティ

### 認証・認可
- **VAPID**: Web Push認証
- **JWT**: API認証（計画中）
- **HTTPS**: 全通信暗号化

### データ保護
- **機密情報**: GitHub Secrets管理
- **データベース**: RDS暗号化
- **通信**: TLS 1.2以上

### セキュリティヘッダー
- Content Security Policy (CSP)
- X-Frame-Options
- X-Content-Type-Options

### アクセス制御
- **IAM**: 最小権限の原則
- **セキュリティグループ**: 必要最小限のポート開放
- **VPC**: プライベートサブネット活用

---

## パフォーマンス

### フロントエンド最適化
- **Turbopack**: 高速ビルド
- **React 19**: 最新パフォーマンス改善
- **Canvas 2D**: リアルタイム画像処理最適化
- **Service Worker**: バックグラウンド処理

### バックエンド最適化
- **Go**: 高性能・低メモリ使用量
- **Connection Pooling**: pgx/v5プール機能
- **非同期処理**: Web Push送信ワーカー

### インフラ最適化
- **ECS Fargate**: オートスケーリング
- **ALB**: 負荷分散
- **RDS**: Multi-AZ高可用性
- **CloudWatch**: 監視・アラート

---

## 開発ガイドライン

### コーディング規約

#### Go
- `gofmt`準拠フォーマット
- 公開識別子: `PascalCase`
- プライベート識別子: `camelCase`
- パッケージ名: 小文字単語
- `panic`禁止、エラーは戻り値で処理

#### TypeScript/React
- 2スペースインデント
- コンポーネント: `PascalCase.tsx`
- ファイル名: kebab-case
- CSS: `*.module.css`

### アーキテクチャ原則
- **DDD**: ビジネスロジックはドメイン層
- **値オブジェクト**: プリミティブ型検証
- **リポジトリパターン**: データアクセス抽象化
- **依存性注入**: 疎結合設計

### テスト戦略
- **単体テスト**: `_test.go`, テーブル駆動テスト
- **統合テスト**: API エンドポイント
- **E2E テスト**: Playwright（計画中）

### コミット規約
- **Conventional Commits**: `feat(server): add user handler`
- **スコープ**: `frontend`, `server`, `infra`
- **必須チェック**: lint/test通過

---

## 特殊機能詳細

### Web Push通知システム
- **RFC準拠**: 8030/8291/8292
- **VAPID認証**: サーバー識別
- **メッセージ暗号化**: エンドツーエンド
- **配信制御**: TTL/Urgency/Topic対応

### カメラフィルター機能
- **リアルタイム処理**: Canvas 2D API
- **フィルター種類**: 5種類（レトロ/ホラー/シリアス/VHS/コミック）
- **画像処理**: Sobel、ポスタライズ、ビネット等
- **カメラ制御**: フロント/バック切り替え

### PWA対応
- **Service Worker**: バックグラウンド処理
- **Web App Manifest**: インストール可能
- **オフライン対応**: キャッシュ戦略

---

## 運用・監視

### メトリクス
- **アプリケーション**: レスポンス時間、エラー率
- **インフラ**: CPU/メモリ使用率、ネットワーク
- **ビジネス**: 通知配信数、クリック率

### ログ管理
- **構造化ログ**: JSON形式
- **ログレベル**: ERROR/WARN/INFO/DEBUG
- **機密情報**: マスキング処理

### アラート
- **システム**: 高CPU/メモリ使用率
- **アプリケーション**: エラー率閾値超過
- **データベース**: 接続数・レスポンス時間

---

## 今後の技術的課題

### スケーラビリティ
- **水平スケーリング**: ECS Auto Scaling
- **データベース**: Read Replica対応
- **キャッシュ**: Redis導入検討

### 機能拡張
- **認証システム**: OAuth2/OIDC対応
- **リアルタイム通信**: WebSocket対応
- **ファイルアップロード**: S3直接アップロード

### 開発効率化
- **コード生成**: sqlc導入
- **API文書**: OpenAPI/Swagger
- **テスト自動化**: カバレッジ向上

---

**最終更新**: 2025-01-06  
**バージョン**: 1.0.0