# プロジェクト構造

## ルートレイアウト

```
├── frontend/          # Next.js Reactアプリケーション
├── server/           # Go バックエンドAPIサーバー
├── infra/            # AWS Terraformインフラストラクチャ設定
├── .kiro/            # Kiro AIアシスタント設定
├── .claude/          # Claude設定
├── .vscode/          # VSCode設定
├── CLAUDE.md         # Claude向けプロジェクト説明
└── README.md         # プロジェクト説明
```

## バックエンド構造 (server/)

Go バックエンドは**ドメイン駆動設計（DDD）**を用いた**レイヤードアーキテクチャ**に従っています：

```
server/
├── main.go                    # アプリケーションエントリーポイント
├── go.mod                     # Goモジュール定義
├── Makefile                   # ビルドと開発コマンド
├── cmd/server/                # コマンドラインインターフェース
├── internal/
│   ├── interfaces/            # インターフェース層（外部インターフェース層）
│   │   └── http/
│   │       ├── handler/       # HTTPハンドラー
│   │       ├── dto/           # データ転送オブジェクト
│   │       └── middleware/    # HTTPミドルウェア
│   ├── application/           # アプリケーション層
│   │   ├── service/          # アプリケーションサービス
│   │   └── usecase/          # ユースケース/インタラクター
│   ├── domain/               # ドメイン層（ビジネスロジック）
│   │   ├── model/            # ドメインエンティティ（User）
│   │   ├── repository/       # リポジトリインターフェース
│   │   ├── service/          # ドメインサービス
│   │   └── valueobject/      # 値オブジェクト（Email、UserID）
│   └── infrastructure/       # インフラストラクチャ層
│       ├── config/           # 設定
│       └── persistence/      # データベース実装
└── pkg/                      # 共有パッケージ
    └── errors/               # エラーハンドリングユーティリティ
```

## インフラストラクチャ構造 (infra/)

AWS Terraform によるインフラストラクチャ定義：

```
infra/
├── alb.tf                  # Application Load Balancer設定
├── ecr.tf                  # Elastic Container Registry
├── ecs_cluster.tf          # ECSクラスター設定
├── ecs_services_api.tf     # APIサービス設定
├── ecs_services_web.tf     # Webサービス設定
├── outputs.tf              # Terraform出力値
├── providers.tf            # AWSとRandomプロバイダー
├── rds.tf                  # RDS PostgreSQLデータベース
├── s3.tf                   # S3バケット設定
├── security.tf             # セキュリティグループとIAMロール
├── variables.tf            # 入力変数
├── versions.tf             # Terraformバージョン制約
├── vpc.tf                  # VPCとネットワーキング
├── terraform.tfvars        # 変数値（機密情報含む）
├── terraform.tfstate       # Terraform状態ファイル
└── README.md               # インフラストラクチャ説明
```

## フロントエンド構造 (frontend/)

標準的な Next.js App Router 構造：

```
frontend/
├── src/app/                  # App Routerページとレイアウト
│   ├── layout.tsx           # ルートレイアウト
│   ├── page.tsx             # ホームページ
│   ├── globals.css          # グローバルスタイル
│   └── page.module.css      # ページ固有スタイル
├── public/                   # 静的アセット
├── package.json             # 依存関係とスクリプト
├── tsconfig.json            # TypeScript設定
├── next.config.ts           # Next.js設定
└── eslint.config.mjs        # ESLint設定
```

## アーキテクチャパターン

### ドメイン層の規約

- **エンティティ**: 振る舞いを持つリッチなドメインオブジェクト（User）
- **値オブジェクト**: バリデーション付きの不変オブジェクト（Email、UserID）
- **リポジトリ**: データアクセス用インターフェース（ドメインで定義、インフラストラクチャで実装）
- **サービス**: エンティティに属さないドメインロジック

### 命名規約

- **Go**: エクスポート関数/型は PascalCase、プライベートは camelCase
- **ファイル**: Go ファイルは snake_case、ディレクトリは kebab-case
- **パッケージ**: 小文字、可能な限り単語 1 つ
- **React**: コンポーネントは PascalCase、関数/変数は camelCase

## 依存関係フロー

レイヤードアーキテクチャの依存関係：

- **Interfaces** → **Application** → **Domain** ← **Infrastructure**
- 依存関係は内側に向かう（クリーンアーキテクチャ）
- Infrastructure は Domain インターフェースを実装

## 開発ガイドライン

- DDD の原則に従う：ビジネスロジックはドメイン層に配置
- 値オブジェクトを使用してプリミティブ型の検証を行う（UserID、Email）
- データアクセスの抽象化にはリポジトリパターンを実装
- カスタムドメインエラー型を通じてエラーを処理
- 疎結合のために依存性注入を使用

### 現在の実装状況

- レイヤードアーキテクチャ + DDD の完全実装済み
- ユーザー管理はインメモリストレージで実装（本格的な DB 対応準備中）
- AWS ECS Fargate での本番環境デプロイメント対応済み
- Terraform によるインフラストラクチャ as Code 実装済み
