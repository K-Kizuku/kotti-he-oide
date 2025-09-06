# プロジェクト構造

## ルートレイアウト
```
├── frontend/          # Next.js Reactアプリケーション
├── server/           # Go バックエンドAPIサーバー
├── infra/            # インフラストラクチャ設定
└── .kiro/            # Kiro AIアシスタント設定
```

## バックエンド構造 (server/)

Goバックエンドは**ドメイン駆動設計（DDD）**を用いた**クリーンアーキテクチャ**に従っています：

```
server/
├── main.go                    # アプリケーションエントリーポイント
├── go.mod                     # Goモジュール定義
├── Makefile                   # ビルドと開発コマンド
├── cmd/server/                # コマンドラインインターフェース（空）
├── internal/
│   ├── application/           # アプリケーション層
│   │   ├── service/          # アプリケーションサービス
│   │   └── usecase/          # ユースケース/インタラクター
│   ├── domain/               # ドメイン層（ビジネスロジック）
│   │   ├── model/            # ドメインエンティティ（User）
│   │   ├── repository/       # リポジトリインターフェース
│   │   ├── service/          # ドメインサービス
│   │   └── valueobject/      # 値オブジェクト（Email、UserID）
│   ├── handler/              # HTTPハンドラー（現在の実装）
│   ├── infrastructure/       # インフラストラクチャ層
│   │   ├── config/           # 設定
│   │   └── persistence/      # データベース実装
│   ├── interfaces/           # インターフェースアダプター
│   │   ├── http/             # HTTPインターフェース層
│   │   │   ├── dto/          # データ転送オブジェクト
│   │   │   ├── handler/      # HTTPハンドラー
│   │   │   └── middleware/   # HTTPミドルウェア
│   │   └── repository/       # リポジトリ実装
│   ├── model/                # データモデル（User、リクエスト）
│   └── service/              # サービス
└── pkg/                      # 共有パッケージ
    └── errors/               # エラーハンドリングユーティリティ
```

## フロントエンド構造 (frontend/)

標準的なNext.js App Router構造：

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
- **Go**: エクスポート関数/型はPascalCase、プライベートはcamelCase
- **ファイル**: Goファイルはsnake_case、ディレクトリはkebab-case
- **パッケージ**: 小文字、可能な限り単語1つ
- **React**: コンポーネントはPascalCase、関数/変数はcamelCase

### 現在の実装に関する注意事項
- ユーザー管理は現在インメモリストレージで実装されています
- HTTPハンドラーは `/internal/handler/` にあります（移行期の構造）
- ドメインモデルは `/internal/domain/model/`（DDD）と `/internal/model/`（シンプル）の両方に存在します
- クリーンアーキテクチャ層は部分的に実装されています