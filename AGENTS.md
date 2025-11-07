# Repository Guidelines

このリポジトリは「赤煉瓦文化館 〜こっちにおいで〜」体験型Webホラーゲームのための開発ガイドラインです。

## プロジェクト概要

福岡市赤煉瓦文化館を舞台とした、時代を超える体験型ホラーゲーム。プレイヤーは1942年→1972年→2002年→2025年と時代を行き来し、自己内省を通じて「この建物が受け継がれてきた意味」を体感します。

## Project Structure & Module Organization

- `frontend/` - Next.js 15（TypeScript）。`src/app` にルーティング、`public/` にアセット
  - `/game` 以下に各シーン（S0〜S9）のコンポーネント
  - カメラフィルター機能（`/camera-filters`）
  - ホラー演出用のオーバーレイコンポーネント
- `server/` - Go（クリーンアーキテクチャ + DDD）
  - `internal/` ドメイン/アプリケーション/インフラ層
  - `cmd/server` エントリポイント
  - `Makefile` で開発タスク
  - ゲームセッション管理、クイズ生成、画像類似度判定
  - VOICEVOX統合（EC2へのプロキシ）
  - Web Push対応（`schema.sql`）
- `infra/` - Terraform IaC
  - VPC/ALB/ECS/ECR/RDS 等のコード
  - EC2（VOICEVOX server）
  - `terraform.tfvars`
- `.github/workflows/` - CI/CD（ECR ビルド→ECS 更新）
  - フロントエンド/バックエンド別々のワークフロー
- `docs/` - 企画書と仕様書
  - `proposal.md` - ゲームの企画書
  - `specification.md` - 開発チーム向け詳細仕様
- `TODO.md` - Web Push機能の詳細仕様書（RFC準拠）

## Build, Test, and Development Commands
Frontend（要 Corepack/Pnpm）
- `cd frontend && corepack enable && pnpm i`
- `pnpm dev` 開発サーバ（http://localhost:3000）
- `pnpm build && pnpm start` 本番ビルド/起動
- `pnpm lint` ESLint（next/core-web-vitals）
- カメラフィルター: `http://localhost:3000/camera-filters`

Server（Go 1.25 系）
- `cd server && make deps` 依存取得 / 整理
- `make run` もしくは `make run-cmd` 起動
- `make test` ユニットテスト実行（`go test ./...`）
- `make fmt && make lint` フォーマット / 静的解析（要 `golangci-lint`）
- DB初期化: `psql -f schema.sql` （PostgreSQL用）

Infra（要 AWS 資格情報）
- `cd infra && terraform init`
- `terraform plan -var-file=terraform.tfvars`
- `terraform apply` 変更適用（要レビュー）

Docker（本番環境）
- `docker build -t api:latest ./server` バックエンドビルド
- `docker build -t web:latest ./frontend` フロントエンドビルド

## Coding Style & Naming Conventions
- Frontend: 2 スペース、TypeScript。コンポーネントは `PascalCase.tsx`、ルートは小文字ケバブ。CSS は `*.module.css`。**ESLint を必ず通すこと（`pnpm lint`）**。
- Server: `gofmt` 準拠。公開識別子は `CamelCase`、パッケージ名は小文字。`panic` 禁止、エラーは戻り値で連鎖させる。**必ず `make fmt && make lint` を実行すること**。

## Code Quality Rules (必須)
**すべてのコード変更時に以下を実行し、エラーがないことを確認する：**
- フロントエンド: `cd frontend && pnpm lint`
- バックエンド: `cd server && make fmt && make lint`
- **lintエラーが残っている状態でのコミット・プッシュは禁止**

## Testing Guidelines
- Server: `_test.go` にテーブル駆動で `TestXxx` を作成。例: `go test -cover ./...`。副作用はモック/インメモリを使用。
- Frontend: クリティカルな UI/ロジックは追加テストを推奨（例: Vitest/RTL）。E2E は Playwright を想定し `frontend/tests` 配下に配置。

## Commit & Pull Request Guidelines
- コミットは Conventional Commits を推奨：`feat(server): add user handler`、`fix(frontend): layout`、`chore(infra): ecs service`
- PR 要件：目的/背景、変更範囲（frontend/server/infra）、関連 Issue、UI 変更はスクリーンショット、Infra は `terraform plan` の抜粋。CI Green を必須（`pnpm lint`/`make test`）。

## Security & Configuration Tips
- 機密情報（API Keys、`.env*`、資格情報）はコミットしない。Next.js の公開値は `NEXT_PUBLIC_` のみ。
- Terraform の状態は原則リモートステート（S3 等）を使用し、不要な状態ファイルのコミットを避ける。
- AWS 資格情報は GitHub Secrets を使用（Workflows 参照）。

## ゲーム機能の概要

### シーン構造（S0〜S9）
- **S0**: 起動・注意書き・カメラ/音声許可
- **S1**: 1942年パート - 担当者との会話（VOICEVOX音声）
- **S2**: お気に入りの場所5箇所の説明
- **S3**: 2階診査室への移動指示 + ホラー演出
- **S4**: 診査室 - 内省10問（逐次保存）
- **S5**: 死亡届受理通知 + 7分タイマー開始（1972年）
- **S6**: 存在証明書探索（5箇所撮影 + クイズ）**最重要**
- **S7**: 2002年パート - メッセージ受け取り
- **S8**: メインホール探索（3分）
- **S9**: 2025年 - 自分のメッセージを刻む

### 主要機能
- **VOICEVOX統合**: EC2上のVOICEVOXサーバーから動的に音声生成。青山龍星（しっとり）を使用。
- **画像類似度判定**: S6での場所到達判定。OpenCV等で基準画像との類似度計算（閾値0.5〜0.6）。
- **セッション管理**: UUID v4ベースのセッションID、60分有効期限、逐次保存でリロード対応。
- **クイズ生成**: S4の回答を元に4択クイズを生成。過去プレイヤー回答を匿名でダミー選択肢に使用。
- **匿名継承システム**: プレイヤーメッセージを匿名保存し、次のプレイヤーに継承。
- **カメラフィルター**: Canvas 2D APIによるリアルタイム画像処理。ホラー演出用のフィルター（色味+ノイズ+影）。
- **タイマー機能**: S6は7分、S8は3分のカウントダウン。クライアント・サーバー両方で管理。
- **ローカライズ**: 日本語/英語対応（キー化）。デフォルトは日本語。

### 既存機能
- **Web Push通知**: RFC 8030/8291/8292準拠。VAPID認証、メッセージ暗号化、PostgreSQL管理。
- **CI/CD**: GitHub Actionsによる自動デプロイ。`frontend/**` または `server/**` の変更で対応するサービスを自動更新。
- **コンテナ化**: Docker マルチステージビルド対応。ECRへの自動プッシュ。

## Communication Rules
- **全てのコミュニケーションは日本語で行う**
- コード内のコメントも可能な限り日本語で記述する
- 変数名や関数名は英語でも構わないが、説明やドキュメントは日本語とする

## Agent-Specific Notes

### 参照すべきドキュメント
- **ゲーム全体の理解**: `docs/proposal.md` と `docs/specification.md` を必ず参照
- **シーン仕様**: `docs/specification.md` の各シーン詳細を確認
- **Web Push機能**: `TODO.md` の仕様書を参照
- **カメラ機能**: `frontend/CAMERA_memo.md` の実装メモを確認（存在する場合）

### 開発時の注意事項
- 自動化エージェントは本ガイドのコマンドのみ実行し、設定ファイルの不要変更を避ける
- ネストした `AGENTS.md` がある場合は最も深い階層を優先
- **コード変更後は必ずlintチェックを実行し、エラーがないことを確認してからコミットすること**
- セッション管理は必ず逐次保存を実装（S4の回答など）
- 画像類似度判定の閾値は調整可能にする
- タイマー処理はクライアント・サーバー両方で実装し、サーバー側を正とする
- ホラー演出は過度にならないよう、ジャンプスケアは最大2回まで
- プレイヤーデータの匿名性を厳守（IP・端末ID等を保存しない）

### API実装時の重要ポイント
- セッションIDは必ずUUID v4を使用
- S6のクイズは入室時に5問すべてを事前生成
- 画像類似度判定失敗時は必ず代替手段（Web選択）を提供
- VOICEVOXサーバーへの接続エラー時のフォールバック処理を実装

