# 赤煉瓦文化館 〜こっちにおいで〜

体験型Webホラーゲームのモノレポです。福岡市赤煉瓦文化館（現エンジニアカフェ）を実空間のまま舞台にし、1942年→1972年→2002年→2025年を行き来する没入体験をWebだけで実現します。


**目次**
- プロジェクト概要 / コンセプト
- リポジトリ構成
- クイックスタート（ローカル開発）
- 実行コマンド一覧（Frontend / Server / Infra / Microservice）
- アーキテクチャ概要（アプリ / インフラ）
- ゲーム機能の要点（S0〜S9）
- 品質ルール・開発ガイドライン
- デプロイとCI/CD
- セキュリティと設定
- 既知の注意点（DB種別の差異 ほか）
- 参照ドキュメント


## プロジェクト概要 / コンセプト
- タイトル: 赤煉瓦文化館 〜こっちにおいで〜
- 目的: 恐怖→解放→内省の感情遷移を通じ、建物が受け継がれてきた意味を体感してもらう。
- 想定プレイヤー: エンジニアカフェ来館者 / 歴史紹介イベント参加者 / 20〜40代社会人。
- プレイ時間: 20〜30分（館内移動込み）。
- セッション: UUID v4 で発行、サーバー側60分有効（逐次保存により復帰可）。


## リポジトリ構成
```
frontend/        # Next.js 15.5.2 + React 19.1.0 + TypeScript（App Router, Turbopack）
server/          # Go 1.25.1 REST API（DDD + レイヤード）
services/
  image_recognition/  # 画像認識 gRPC マイクロサービス（Python, OpenCV）
infra/           # Terraform（AWS: ECS Fargate/ALB/ECR/RDS/S3/EC2 ほか）
docs/            # 企画書・仕様書（proposal.md / specification.md）
.github/workflows/    # CI/CD（フロント/サーバー/マイクロサービス/インフラ）
.kiro/steering/  # 仕様の単一情報源（必ず参照）
```


## クイックスタート（ローカル開発）

前提ツール
- Node.js 20系 / pnpm（Corepack推奨）
- Go 1.25.1
- Python 3.12 + `uv`（画像認識サービスを動かす場合）
- Docker（任意） / Terraform 1.6+（インフラ作業時） / AWS CLI
- Buf CLI v2（gRPCスタブ生成）

環境変数（最低限）
- フロントエンド: `NEXT_PUBLIC_API_URL`（例: `http://localhost:8080`）
- サーバー: `PORT`（既定 8080）、`DATABASE_URL`、`SESSION_TTL_MINUTES`（既定60）
- Web Push: `VAPID_PUBLIC_KEY`, `VAPID_PRIVATE_KEY`（`go run cmd/generate-vapid/main.go` で生成可）

1) サーバー（Go）
```bash
cd server
# 依存取得
make deps
# 開発起動（:8080）
make run
```

2) フロントエンド（Next.js）
```bash
cd frontend
corepack enable && pnpm install
export NEXT_PUBLIC_API_URL=http://localhost:8080
pnpm dev   # http://localhost:3000
```

3) 画像認識マイクロサービス（任意）
```bash
cd services/image_recognition
uv sync
# リポジトリルートでスタブ生成（推奨）
buf generate
uv run python -m app.server  # :50051
```

動作確認
- フロント: http://localhost:3000 を開く →「ゲームを開始」で S0 へ
- サーバー: http://localhost:8080/api/healthz で `{"status":"ok"}` を確認


## 実行コマンド一覧（抜粋）

Frontend（pnpm）
```bash
cd frontend
pnpm dev     # 開発サーバー（Turbopack）
pnpm build   # 本番ビルド
pnpm start   # 本番起動
pnpm lint    # ESLint（必須）
```

Server（Go）
```bash
cd server
make run     # 開発起動 (:8080)
make dev     # 自動リロード（air 必要）
make build   # bin/server へビルド
make test    # テスト実行（PoCのため新規作成不要）
make deps    # 依存取得
make fmt     # gofmt（必須）
make lint    # golangci-lint（必須）
```

Infrastructure（Terraform）
```bash
cd infra
terraform init
terraform plan
terraform apply
terraform destroy
```

Docker（本番想定ビルド）
```bash
docker build -t api:latest ./server
docker build -t web:latest ./frontend
```


## アーキテクチャ概要
- フロントエンド: Next.js + React（App Router, Canvas 2D フィルター, Web Push UI, Service Worker）
- バックエンド: Go + net/http（DDD + レイヤード、セッション/回答/クイズ/メッセージ/Web Push API）
- マイクロサービス: 画像認識（Python gRPC, OpenCV, S3参照画像, 類似度スコア）
- インフラ: AWS ECS Fargate（API/Web/Microservice）、ALB（`/`→Web, `/api/*`→API）、RDS、S3、EC2(VOICEVOX)

API エンドポイント例（抜粋）
- Health: `GET /api/healthz`
- セッション: `POST /api/session`, `GET /api/session/{session_id}`
- S4回答: `POST /api/session/{session_id}/answers`, `GET /api/session/{session_id}/answers`
- S6: `POST /api/session/{session_id}/s6/initialize`, `POST /api/session/{session_id}/s6/verify-location`, `GET /api/session/{session_id}/s6/quiz/{place_id}`, `POST /api/session/{session_id}/s6/answer`, `GET /api/session/{session_id}/s6/progress`
- メッセージ: `POST /api/session/{session_id}/message`, `GET /api/messages`
- Web Push: `GET /api/push/vapid-public-key`, `POST /api/push/subscribe`, `DELETE /api/push/subscriptions/{subscription_id}`, `POST /api/push/send/{session_id}`


## ゲーム機能の要点（S0〜S9）
- S0: 起動・注意書き・カメラ/音声許可（許可失敗時はカメラなしモードで継続可能）
- S1: 1942年パート（VOICEVOX音声、来館者情報を保存）
- S2: お気に入りの場所5箇所の説明（後続S6の探索の前提）
- S3: 2階診査室への移動指示 + ホラー演出（フィルター/影/虚な人）
- S4: 診査室 - 内省10問（“なし/特にない”は弾く、逐次保存）
- S5: 死亡届受理通知 + 7分タイマー開始（1972年）
- S6: 存在証明書探索（5箇所撮影 + 類似度判定 + 4択クイズ）最重要
- S7: 2002年パート - メッセージ受け取り
- S8: メインホール探索（3分）
- S9: 2025年 - 自分のメッセージを刻む（匿名継承の一部に）

技術的ポイント
- 音声動的生成: VOICEVOX（EC2）をHTTPで呼び出し
- 画像認識: 類似度判定（推奨閾値 0.5〜0.6, 調整可能）
- タイマー: クライアント/サーバー双方で管理（サーバーを正とする）
- フォールバック: 画像判定失敗時はWeb選択で代替進行
- ジャンプスケアは最大2回まで


## 品質ルール・開発ガイドライン（必須）
- 仕様の参照: `.kiro/steering/` を常に正とする
- Lint必須: フロント `pnpm lint` / サーバー `make fmt && make lint`
- コードスタイル: Frontend 2スペース, TS/Reactの命名規約遵守。Goは `gofmt`/命名規約、`panic`禁止。
- テスト方針: 本プロジェクトはPoC。新規テストコードの作成は不要。
- セキュリティ: 機密情報（API Keys, `.env*` 等）をコミットしない。Next.js 公開値は `NEXT_PUBLIC_` のみ。
- セッション: UUID v4、60分TTL、逐次保存でリロード復帰可。

コミット/PR
- Conventional Commits 推奨（例: `feat(server): add user handler`）
- PRには目的/背景、変更範囲（frontend/server/infra）、関連Issue、UI変更のスクショ、Infraは `terraform plan` の抜粋、CI Green を添付


## デプロイとCI/CD
- GitHub Actions で ECR ビルド→ECS 更新（フロント/サーバー/マイクロサービス/インフラ）
- パスベースでワークフローが発火（`frontend/**`, `server/**`, `services/image_recognition/**`, `infra/**`）
- 環境変数/シークレットは GitHub Secrets/Variables を使用（AWS資格情報、`NEXT_PUBLIC_API_URL` 等）


## セキュリティと設定
- Web Push は HTTPS + VAPID キーが必須。`VAPID_PUBLIC_KEY` / `VAPID_PRIVATE_KEY` を環境変数に設定
- Terraform の状態は原則リモートステート（S3等）を使用。不要な状態ファイルのコミットは避ける
- プレイヤーデータは匿名性を厳守（IP・端末ID等は保存しない）


## 既知の注意点（差異について）
- RDBMS 仕様: 本プロジェクトでは RDS MySQL 8.0 を採用しています。`server/` 実装および `infra/rds.tf` は MySQL ベースで統一されています（例: `NewMySQLDB`, `engine = "mysql"`）。
- 開発中機能: 画像認識やWeb Pushの一部はPoC段階。フォールバック（Web選択/通知なし進行）を必ず実装。


## 参照ドキュメント
- 仕様の単一情報源: `.kiro/steering/product.md`, `.kiro/steering/tech.md`, `.kiro/steering/structure.md`
- 企画/仕様: `docs/proposal.md`, `docs/specification.md`
- インフラ手順: `infra/README.md`
- サーバーセットアップ: `server/SETUP.md`, 実装状況: `server/IMPLEMENTATION_STATUS.md`
- Web Push スキーマ: `server/schema.sql`

コミュニケーションは原則日本語。コード内コメントも可能な限り日本語で記述してください。
