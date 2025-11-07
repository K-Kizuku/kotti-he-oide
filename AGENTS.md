# Repository Guidelines

このリポジトリは「赤煉瓦文化館 〜こっちにおいで〜」体験型Webホラーゲームのための開発ガイドラインです。

**重要**: このドキュメントの内容は `.kiro/steering/` ディレクトリ内のファイル（product.md、tech.md、structure.md）を常に参照し、最新の仕様に準拠してください。steeringファイルの内容が最も正確で最新の情報源です。

## プロジェクト概要

福岡市赤煉瓦文化館（現エンジニアカフェ）を舞台とした、体験型Webホラーゲームです。ハッカソンプロジェクトとして開発され、以下の主要コンポーネントで構成されています：

- **フロントエンド**: Next.js 15.5.2 + React 19.1.0 + TypeScript（App Router使用）
- **バックエンド**: Go 1.25.1 REST APIサーバー（レイヤードアーキテクチャ + DDD）
- **インフラストラクチャ**: AWS ECS Fargate + Terraform（本番環境対応）
- **VOICEVOX Server**: AWS EC2上で動作する音声生成サーバー

プレイヤーは1942年→1972年→2002年→2025年と時代を行き来し、自己内省を通じて「この建物が受け継がれてきた意味」を体感します。

## Project Structure & Module Organization

- `frontend/` - Next.js 15.5.2（TypeScript + Turbopack）。`src/app` にルーティング、`public/` にアセット
  - `/game` 以下に各シーン（S0〜S9）のコンポーネント
  - カメラフィルター機能（`/camera-filters`）
  - ホラー演出用のオーバーレイコンポーネント
- `server/` - Go 1.25.1（クリーンアーキテクチャ + DDD）
  - `internal/` ドメイン/アプリケーション/インフラ層
  - `cmd/server` エントリポイント
  - `Makefile` で開発タスク
  - ゲームセッション管理、クイズ生成、画像類似度判定
  - VOICEVOX統合（EC2へのプロキシ）
  - Web Push対応（`schema.sql`）
- `services/` - マイクロサービス（画像認識等）
  - `image_recognition/` - Python + gRPC画像認識サービス
- `infra/` - Terraform IaC
  - VPC/ALB/ECS/ECR/RDS 等のコード
  - EC2（VOICEVOX server）
  - `terraform.tfvars`
- `.github/workflows/` - CI/CD（ECR ビルド→ECS 更新）
  - フロントエンド/バックエンド別々のワークフロー
- `docs/` - 企画書と仕様書
  - `proposal.md` - ゲームの企画書
  - `specification.md` - 開発チーム向け詳細仕様

## Build, Test, and Development Commands

### Frontend（pnpm使用）
```bash
cd frontend
pnpm dev          # Turbopackで開発サーバーを開始（http://localhost:3000）
pnpm build        # Turbopackで本番用ビルド
pnpm start        # 本番サーバーを開始
pnpm lint         # ESLintを実行（必須）
```

### Server（Go 1.25.1）
```bash
cd server
make run          # 開発モードでサーバーを実行（port 8080）
make dev          # 自動リロードで実行（airが必要）
make build        # bin/serverにバイナリをビルド
make test         # Goテストを実行（PoCのため新規テスト作成は不要）
make deps         # 依存関係をインストール/更新
make fmt          # Goコードをフォーマット（必須）
make lint         # リンターを実行（golangci-lintが必要、必須）
make clean        # ビルド成果物をクリーン
```

### Infrastructure（Terraform）
```bash
cd infra
terraform init      # Terraformを初期化
terraform plan      # インフラストラクチャ変更をプレビュー
terraform apply     # インフラストラクチャ変更を適用
terraform destroy   # インフラストラクチャリソースを削除
```

### Docker（本番環境）
```bash
docker build -t api:latest ./server    # バックエンドビルド
docker build -t web:latest ./frontend  # フロントエンドビルド
```

## Coding Style & Naming Conventions

### Frontend
- 2スペースインデント、TypeScript使用
- コンポーネントは `PascalCase.tsx`
- ルートは小文字ケバブケース
- CSS は `*.module.css`
- **ESLint を必ず通すこと（`pnpm lint`）**

### Backend（Go）
- `gofmt` 準拠
- エクスポート関数/型は `PascalCase`、プライベートは `camelCase`
- Goファイルは `snake_case`、ディレクトリは `kebab-case`
- パッケージ名は小文字、可能な限り単語1つ
- `panic` 禁止、エラーは戻り値で連鎖させる
- **必ず `make fmt && make lint` を実行すること**

### React
- コンポーネントは `PascalCase`
- 関数/変数は `camelCase`

## Code Quality Rules (必須)
**すべてのコード変更時に以下を実行し、エラーがないことを確認する：**
- フロントエンド: `cd frontend && pnpm lint`
- バックエンド: `cd server && make fmt && make lint`
- **lintエラーが残っている状態でのコミット・プッシュは禁止**

## テストポリシー

**重要: このプロジェクトはPoC（Proof of Concept）です。**

- **テストコードの作成は一切不要**です
- 開発速度とプロトタイピングを最優先します
- ユニットテスト、統合テスト、E2Eテストなど、あらゆる種類のテストコードの実装は求められません
- `make test` コマンドは存在しますが、新規テストファイルの作成は不要です

## Commit & Pull Request Guidelines

### Commit Messages
- Conventional Commits を推奨
- 例: `feat(server): add user handler`、`fix(frontend): layout`、`chore(infra): ecs service`

### Pull Request Requirements
- 目的/背景を明記
- 変更範囲（frontend/server/infra）を記載
- 関連Issueを参照
- UI変更はスクリーンショットを添付
- Infraは `terraform plan` の抜粋を含める
- CI Greenを必須（`pnpm lint`/`make test`）

## Security & Configuration Tips

- 機密情報（API Keys、`.env*`、資格情報）はコミットしない
- Next.jsの公開値は `NEXT_PUBLIC_` プレフィックスのみ使用
- Terraformの状態は原則リモートステート（S3等）を使用
- 不要な状態ファイルのコミットを避ける
- AWS資格情報はGitHub Secretsを使用（Workflows参照）

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

### 技術的特徴
- **音声の動的生成**: VOICEVOX（青山龍星(しっとり)）をEC2上で動作させ、アプリケーションサーバからHTTPで音声を生成
- **AR的演出**: カメラ映像にホラー用オーバーレイ（影・虚な人）を重ねる
- **画像認識**: 館内5カ所の撮影画像を基準画像と類似度比較して場所到達を判定（閾値0.5〜0.6）
- **匿名継承コンセプト**: 過去プレイヤーの回答が匿名化されて次のプレイヤーのクイズのダミーになる
- **セッション管理**: UUID v4ベースのセッションID、60分有効期限、逐次保存でリロード対応
- **クイズ生成**: S4の回答を元に4択クイズを生成。過去プレイヤー回答を匿名でダミー選択肢に使用
- **カメラフィルター**: Canvas 2D APIによるリアルタイム画像処理。ホラー演出用のフィルター（色味+ノイズ+影）
- **タイマー機能**: S6は7分、S8は3分のカウントダウン。クライアント・サーバー両方で管理

### インフラストラクチャ
- **コンピュート**: ECS Fargate クラスター（コンテナ化されたAPIとWebサービス）
- **ロードバランサー**: ALB（パスベースルーティング：`/` → Web、`/api/*` → API）
- **データベース**: RDS MySQL（セッション管理、プレイヤー回答、メッセージ保存）
- **ストレージ**: S3（静的アセット、撮影画像、音声ファイル用）
- **音声生成**: EC2上のVOICEVOXサーバ
- **画像認識**: gRPCマイクロサービス（Python）
- **ネットワーク**: カスタムVPC（パブリック/プライベートサブネット）
- **コンテナレジストリ**: ECR（APIとWebイメージ用）
- **CI/CD**: GitHub Actions（ECR/ECS自動デプロイ）

## 開発段階

現在は本格的なプロダクション対応に移行済み：
- ドメイン駆動設計（DDD）の完全実装
- MySQLデータベース対応（ゲームセッション、プレイヤー回答、メッセージ保存）
- Dockerコンテナ化とCI/CD自動デプロイ
- AWS ECS Fargateでの本番環境運用

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
- **必ず `.kiro/steering/` 内のファイルを参照し、最新の仕様に準拠すること**
- **コード変更後は必ずlintチェックを実行し、エラーがないことを確認してからコミットすること**
- **このプロジェクトはPoCのため、テストコードの作成は一切不要**
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

