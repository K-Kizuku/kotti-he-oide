# image_recognition マイクロサービス（Python / FastAPI）

gRPC（python / grpcio）による画像認識サーバ実装と、`uv`（Python パッケージマネージャー）での依存管理、Docker ビルド環境を提供します。proto はリポジトリ共通の `schema/proto` に配置されます。

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
- rpc:
  - `Hello(HelloRequest) returns (HelloReply)`
  - `RecognizeImage(RecognizeImageRequest) returns (RecognizeImageResponse)`
  - `HealthCheck(HealthCheckRequest) returns (HealthCheckResponse)`

### メッセージ定義（抜粋）
```proto
message RecognizeImageRequest {
  bytes image_data = 1;   // JPEG/PNG/WebP/BMP 等
  float threshold = 2;    // 0.0-1.0、未指定時は既定（0.6）
  string place_id = 3;    // 任意: カテゴリ/場所IDで参照画像を絞り込む
}

message RecognizeImageResponse {
  bool is_match = 1;
  float similarity_score = 2; // 0.0-1.0
  string error_message = 3;   // 異常時の説明
}
```

## プロジェクト構成
```
services/image_recognition/
├─ app/
│  ├─ __init__.py
│  ├─ server.py         # gRPC サーバー エントリポイント
│  └─ gen/              # 生成される Python スタブ
│     └─ image_recognition/v1/*_pb2*.py
│  └─ models/
│     └─ place_ids.py   # place_id 語彙の共有定義
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
make type    # mypy 型チェック
make lint    # ruff リント
make test    # pytest
```

## 環境変数
- `REFERENCE_S3_BUCKET`: 参照画像を格納した S3 バケット名
- `REFERENCE_S3_PREFIX`: 参照画像のキー Prefix（任意）
- `AWS_REGION` or `AWS_DEFAULT_REGION`: S3 用リージョン
- `DEFAULT_SIMILARITY_THRESHOLD`: 類似度の既定しきい値（デフォルト 0.6）
- `GRPC_PORT`: gRPC ポート（デフォルト: 50051）
- `GRPC_TLS_CERT_FILE`, `GRPC_TLS_KEY_FILE`: 指定時に TLS 有効化

## メモ
- インフラ（ECS/ECR/Terraform など）は別担当が実装する想定です。本ディレクトリでは扱いません。
- 機密情報は `.env` 等に保持し、コミットしないでください。
- コード内コメントは日本語を推奨します。

## 関連ドキュメント
- `services/image_recognition/PLACE_IDS.md` — place_idの語彙とS3配置規約
- `services/image_recognition/INTEGRATION.md` — Goサーバ/フロントとの繋ぎ込みガイド

## ヘルスチェック仕様
- `HealthCheck` はサーバープロセスが応答可能かのみを返し、常に `healthy=true, status="ok"` を返します。
- S3疎通や参照画像のロード可否は起動ログに出力します（失敗してもサーバーは起動継続）。

## 識別不可時のフォールバック（簡易実装）
- 記述子が抽出できない、または（place_id指定時に）該当参照が0件の場合は、比較処理を行わず80%の確率で「同じ」と判定します。
- 判定はログに `recognize fallback_random` として記録されます。
