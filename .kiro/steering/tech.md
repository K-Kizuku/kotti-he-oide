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

## 開発ツール
- **リンティング**: Next.js設定でのESLint
- **コードフォーマット**: バックエンド用Go fmt
- **オプションツール**: air (Goホットリロード用), golangci-lint (Goリンティング用)

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

## デフォルトポート
- バックエンド: 8080 (PORT環境変数で設定可能)
- フロントエンド: 3000 (Next.jsデフォルト)