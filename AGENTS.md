# Repository Guidelines

## Project Structure & Module Organization
- `frontend/` Next.js 15（TypeScript）。`src/app` にルーティング、`public/` にアセット。
- `server/` Go（クリーンアーキ風）。`internal/` ドメイン/アプリ層、`cmd/server` エントリ、`Makefile` で開発タスク。
- `infra/` Terraform。VPC/ALB/ECS/ECR/RDS 等のコードと `terraform.tfvars`。
- `.github/workflows/` CI/CD（ECR ビルド→ECS 更新）。

## Build, Test, and Development Commands
Frontend（要 Corepack/Pnpm）
- `cd frontend && corepack enable && pnpm i`
- `pnpm dev` 開発サーバ（http://localhost:3000）
- `pnpm build && pnpm start` 本番ビルド/起動
- `pnpm lint` ESLint（next/core-web-vitals）

Server（Go 1.25 系）
- `cd server && make deps` 依存取得 / 整理
- `make run` もしくは `make run-cmd` 起動
- `make test` ユニットテスト実行（`go test ./...`）
- `make fmt && make lint` フォーマット / 静的解析（要 `golangci-lint`）

Infra（要 AWS 資格情報）
- `cd infra && terraform init`
- `terraform plan -var-file=terraform.tfvars`
- `terraform apply` 変更適用（要レビュー）

## Coding Style & Naming Conventions
- Frontend: 2 スペース、TypeScript。コンポーネントは `PascalCase.tsx`、ルートは小文字ケバブ。CSS は `*.module.css`。ESLint を通すこと。
- Server: `gofmt` 準拠。公開識別子は `CamelCase`、パッケージ名は小文字。`panic` 禁止、エラーは戻り値で連鎖させる。

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

## Agent-Specific Notes
- 自動化エージェントは本ガイドのコマンドのみ実行し、設定ファイルの不要変更を避ける。ネストした `AGENTS.md` がある場合は最も深い階層を優先。

