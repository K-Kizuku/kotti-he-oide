# プロジェクト構造

## ルートレイアウト

```
├── frontend/          # Next.js Reactアプリケーション
├── server/           # Go バックエンドAPIサーバー
├── services/         # マイクロサービス（画像認識等）
├── infra/            # AWS Terraformインフラストラクチャ設定
├── docs/             # 企画書・仕様書
├── .github/          # GitHub Actions CI/CDワークフロー
├── .kiro/            # Kiro AIアシスタント設定
├── .claude/          # Claude設定
├── .vscode/          # VSCode設定
├── .mcp.json         # Model Context Protocol設定
├── AGENTS.md         # エージェント向けガイドライン
└── README.md         # プロジェクト説明
```

## バックエンド構造 (server/)

Go バックエンドは**ドメイン駆動設計（DDD）**を用いた**レイヤードアーキテクチャ**に従っています：

```
server/
├── main.go                    # アプリケーションエントリーポイント
├── go.mod                     # Goモジュール定義
├── go.sum                     # 依存関係チェックサム
├── Makefile                   # ビルドと開発コマンド
├── Dockerfile                 # コンテナビルド定義
├── schema.sql                 # MySQL用データベーススキーマ
├── bin/                       # ビルド成果物
├── cmd/server/                # コマンドラインインターフェース
├── internal/
│   ├── interfaces/            # インターフェース層（外部インターフェース層）
│   │   └── http/
│   │       ├── handler/       # HTTPハンドラー（ゲームセッション、クイズ、メッセージ等）
│   │       ├── dto/           # データ転送オブジェクト
│   │       └── middleware/    # HTTPミドルウェア
│   ├── application/           # アプリケーション層
│   │   ├── service/          # アプリケーションサービス（音声生成、画像認識連携）
│   │   └── usecase/          # ユースケース/インタラクター（ゲームフロー制御）
│   ├── domain/               # ドメイン層（ビジネスロジック）
│   │   ├── model/            # ドメインエンティティ（Session、PlayerAnswer、Message、Quiz）
│   │   ├── repository/       # リポジトリインターフェース
│   │   ├── service/          # ドメインサービス（クイズ生成、タイマー管理）
│   │   └── valueobject/      # 値オブジェクト（SessionID、PlaceID、QuizID）
│   ├── infrastructure/       # インフラストラクチャ層
│   │   ├── config/           # 設定
│   │   └── persistence/      # データベース実装
│   └── gen/                  # Protocol Buffers生成コード
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
├── rds.tf                  # RDS MySQLデータベース
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
│   ├── game/                # ゲームメインフロー
│   │   ├── s0/             # 起動・注意書き・許可
│   │   ├── s1/             # 1942年パート（担当者との会話）
│   │   ├── s2/             # お気に入りの場所説明
│   │   ├── s3/             # 移動指示＋ホラー演出
│   │   ├── s4/             # 診査室：内省質問
│   │   ├── s5/             # 死亡届受理＋7分制限開始
│   │   ├── s6/             # 存在証明書探索パート
│   │   ├── s7/             # 2002年パート
│   │   ├── s8/             # メインホール探索
│   │   └── s9/             # 2025年・メッセージ刻み
│   ├── layout.tsx           # ルートレイアウト
│   ├── page.tsx             # ホームページ（QRコード案内）
│   ├── globals.css          # グローバルスタイル
│   └── page.module.css      # ページ固有スタイル
├── src/components/           # 共通コンポーネント
│   ├── camera/              # カメラ関連
│   ├── audio/               # 音声再生
│   ├── horror/              # ホラー演出オーバーレイ
│   └── quiz/                # クイズUI
├── src/lib/                 # ユーティリティ
│   ├── api/                # APIクライアント
│   └── utils/              # 汎用関数
├── public/                   # 静的アセット
│   ├── fonts/               # フォント（onryou.TTF等）
│   └── icons/               # アイコン
├── .next/                    # Next.jsビルド成果物
├── node_modules/             # 依存関係
├── Dockerfile                # コンテナビルド定義
├── package.json             # 依存関係とスクリプト
├── pnpm-lock.yaml           # パッケージロックファイル
├── tsconfig.json            # TypeScript設定
├── next.config.ts           # Next.js設定
└── eslint.config.mjs        # ESLint設定
```

## マイクロサービス構造 (services/)

```
services/
└── image_recognition/        # 画像認識サービス（Python + gRPC）
    ├── app/
    │   ├── server.py        # gRPCサーバー
    │   ├── services/        # 画像認識ロジック
    │   ├── models/          # データモデル
    │   └── utils/           # ユーティリティ
    ├── Dockerfile           # コンテナビルド定義
    ├── pyproject.toml       # Python依存関係
    └── README.md            # サービス説明
```

## アーキテクチャパターン

### ドメイン層の規約

- **エンティティ**: 振る舞いを持つリッチなドメインオブジェクト（Session、PlayerAnswer、Message、Quiz）
- **値オブジェクト**: バリデーション付きの不変オブジェクト（SessionID、PlaceID、QuizID）
- **リポジトリ**: データアクセス用インターフェース（ドメインで定義、インフラストラクチャで実装）
- **サービス**: エンティティに属さないドメインロジック（クイズ生成、タイマー管理、音声生成連携）

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

## テストポリシー

**重要: このプロジェクトはPoC（Proof of Concept）です。**

- **テストコードの作成は一切不要**です
- 開発速度とプロトタイピングを最優先します
- ユニットテスト、統合テスト、E2Eテストなど、あらゆる種類のテストコードの実装は求められません

## コード品質管理

**必須ルール: すべてのコード変更時にlintチェックを実行し、エラーがないことを確認する**

- フロントエンド: `cd frontend && pnpm lint`
- バックエンド: `cd server && make fmt && make lint`
- lintエラーが残っている状態でのコミット・プッシュは禁止
- CI/CDパイプラインでもlintチェックが自動実行される

## CI/CD構造 (.github/workflows/)

GitHub Actionsによる自動デプロイメント：

```
.github/workflows/
├── frontend-deploy.yml       # フロントエンドデプロイ（frontend/**変更時）
└── server-deploy.yml         # バックエンドデプロイ（server/**変更時）
```

### 現在の実装状況

- レイヤードアーキテクチャ + DDD の完全実装済み
- MySQLデータベース対応（ゲームセッション、プレイヤー回答、メッセージ保存）
- Dockerコンテナ化とCI/CD自動デプロイ完備
- AWS ECS Fargate での本番環境運用中
- Terraform によるインフラストラクチャ as Code 実装済み
- 画像認識マイクロサービス（gRPC、類似度計算）
- カメラ機能（撮影、ホラー演出オーバーレイ）
- VOICEVOX音声生成連携（EC2上で動作）
