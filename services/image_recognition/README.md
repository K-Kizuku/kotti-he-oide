# image_recognition マイクロサービス（Python / FastAPI）

gRPC（python / grpcio）による最小サーバ実装と、`uv`（Python パッケージマネージャー）での依存管理、Docker ビルド環境を提供します。proto はリポジトリ共通の `schema/proto` に配置されます。

## 前提
- Python 3.12+（ローカル開発時）
- `uv` インストール済み（未インストールの場合）:
  - macOS/Linux: `curl -LsSf https://astral.sh/uv/install.sh | sh`
  - Homebrew: `brew install uv`

## クイックスタート（ローカル）
```bash
cd services/image_recognition
uv sync                          # 依存の解決と仮想環境作成
## 1) Buf (推奨) でスタブ生成
# リポジトリルートで
buf generate

## 2) 直接 grpcio-tools で生成（代替）
uv run python -m grpc_tools.protoc \
  -I ../../schema/proto \
  --python_out app/gen \
  --grpc_python_out app/gen \
  ../../schema/proto/image_recognition/v1/image_recognition.proto

# gRPC サーバ起動（デフォルト :50051）
uv run python -m app.server
```

- 動作確認（GoサーバのHTTP経由）: 後述の `/api/ml/hello` を参照

## Docker で起動
```bash
cd <リポジトリルート>
docker build -f services/image_recognition/Dockerfile -t image-recognition:dev .
docker run --rm -p 50051:50051 image-recognition:dev
```

## gRPC サービス
- service: `image_recognition.v1.ImageRecognitionService`
- rpc: `Hello(HelloRequest) returns (HelloReply)`

## プロジェクト構成
```
services/image_recognition/
├─ app/
│  ├─ __init__.py
│  ├─ server.py         # gRPC サーバー エントリポイント
│  └─ gen/              # 生成される Python スタブ
│     └─ image_recognition/v1/*_pb2*.py
├─ Dockerfile           # uv と grpcio-tools による生成 + 実行
├─ Makefile             # 開発用コマンド
├─ pyproject.toml       # 依存定義（uv 管理）
├─ .dockerignore
└─ .gitignore
```

## よく使うコマンド
```bash
make setup   # uv sync（初回/更新時）
make dev     # （任意）ローカル開発補助
make run     # gRPC サーバー起動
```

## メモ
- インフラ（ECS/ECR/Terraform など）は別担当が実装する想定です。本ディレクトリでは扱いません。
- 機密情報は `.env` 等に保持し、コミットしないでください。
- コード内コメントは日本語を推奨します。
