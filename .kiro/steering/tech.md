# 技術スタック

## バックエンド (Go)
- **言語**: Go 1.25.1
- **フレームワーク**: 標準ライブラリHTTPサーバー (net/http)
- **アーキテクチャ**: ドメイン駆動設計（DDD）を用いたクリーンアーキテクチャ
- **モジュール**: `github.com/K-Kizuku/kotti-he-oide`

## フロントエンド (Next.js)
- **フレームワーク**: Next.js 15.5.2 with App Router
- **ランタイム**: React 19.1.0
- **言語**: TypeScript 5.x
- **パッケージマネージャー**: pnpm (pnpm-lock.yamlに基づく)
- **ビルドツール**: Turbopack (--turbopackフラグを使用)

## インフラストラクチャ (AWS)
- **コンテナオーケストレーション**: ECS Fargate
- **ロードバランサー**: Application Load Balancer (ALB)
- **データベース**: RDS PostgreSQL（Web Push用スキーマ含む）
- **コンテナレジストリ**: Elastic Container Registry (ECR)
- **ストレージ**: S3
- **ネットワーク**: VPC、パブリック/プライベートサブネット
- **IaC**: Terraform ~> 5.0
- **ログ**: CloudWatch Logs
- **CI/CD**: GitHub Actions（ECR/ECS自動デプロイ）

## 開発ツール
- **リンティング**: Next.js設定でのESLint
- **コードフォーマット**: バックエンド用Go fmt
- **オプションツール**: air (Goホットリロード用), golangci-lint (Goリンティング用)
- **コンテナ**: Docker（マルチステージビルド対応）
- **MCP**: Model Context Protocol設定（context7サーバー）

## 主要ライブラリ
### バックエンド
- **Web Push**: `github.com/SherClockHolmes/webpush-go` v1.4.0
- **JWT**: `github.com/golang-jwt/jwt/v5` v5.2.1
- **PostgreSQL**: `github.com/jackc/pgx/v5` v5.7.5
- **DI**: `github.com/google/wire` v0.7.0

### フロントエンド
- **カメラ処理**: Canvas 2D API、getUserMedia
- **画像フィルター**: 自作ImageData処理（Sobel、ポスタライズ等）

## よく使用するコマンド

### バックエンド (server/)
```bash
# 開発
make run          # 開発モードでサーバーを実行
make dev          # 自動リロードで実行 (airが必要)

# ビルド & テスト
make build        # bin/serverにバイナリをビルド
make test         # Goテストを実行
make deps         # 依存関係をインストール/更新

# コード品質
make fmt          # Goコードをフォーマット
make lint         # リンターを実行 (golangci-lintが必要)
make clean        # ビルド成果物をクリーン
```

### フロントエンド (frontend/)
```bash
# 開発
pnpm dev          # Turbopackで開発サーバーを開始
pnpm build        # Turbopackで本番用ビルド
pnpm start        # 本番サーバーを開始
pnpm lint         # ESLintを実行
```

### インフラストラクチャ (infra/)
```bash
# Terraform操作
terraform init      # Terraformを初期化
terraform plan      # インフラストラクチャ変更をプレビュー
terraform apply     # インフラストラクチャ変更を適用
terraform destroy   # インフラストラクチャリソースを削除

# ECRへのイメージプッシュ（例：東京リージョン）
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin <ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com
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
```

## デフォルトポート
- バックエンド: 8080 (PORT環境変数で設定可能)
- フロントエンド: 3000 (Next.jsデフォルト)
- ALB: 80/443 (HTTP/HTTPS)