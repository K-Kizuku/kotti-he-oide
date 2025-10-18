# 画像認識マイクロサービス 引き継ぎメモ

本ドキュメントは 2025-09-07 時点の進捗と、後任者が実装・運用を引き継ぐための要点をまとめたものです。

## 実装サマリ
- 追加領域: `services/image_recognition/`（Python / gRPC / uv）
- RPC: `Hello`, `RecognizeImage`, `HealthCheck`
- プロトコル: `schema/proto/image_recognition/v1/image_recognition.proto`
  - `RecognizeImageRequest.threshold` は proto3 `optional` として実装（未指定時はサーバ既定を適用可能）
- 主要コンポーネント
  - `app/server.py`: gRPC サーバ（asyncio, grpc.aio）
  - `app/services/image_service.py`: 類似度計算・しきい値判定
  - `app/services/s3_service.py`: 参照画像の列挙/取得
  - `app/utils/image_processor.py`: 前処理/特徴抽出（ORB）/類似度算出（BFMatcher）
  - `app/utils/config.py`: 環境変数管理（`REFERENCE_S3_BUCKET` または `S3_BUCKET_NAME`）
  - `app/utils/logger.py`: CloudWatch JSON ログ
  - 生成スタブ: `app/gen/image_recognition/v1/*_pb2*.py`

## 現状の動作状況
- Python サービスは起動可能（S3 未設定でも起動はするが、参照画像が空のため常に `is_match=false` になる）。
- 単体テストを一部実装（前処理/類似度）。統合・E2E テストはスコープ外のため未実装。
- Dockerfile は整備済み（uv で依存解決→proto 生成→gRPC 起動）。

## ローカル開発手順
```bash
cd services/image_recognition
uv sync                       # 依存解決
uv run python -m app.server   # :50051 で起動（環境変数で変更可）

# プロトコル変更時のスタブ再生成
uv run python -m grpc_tools.protoc \
  -I ../../schema/proto \
  --python_out app/gen \
  --grpc_python_out app/gen \
  ../../schema/proto/image_recognition/v1/image_recognition.proto
```

環境変数（例は `.env.example` 参照）
- `REFERENCE_S3_BUCKET` または `S3_BUCKET_NAME`
- `REFERENCE_S3_PREFIX`（任意）
- `AWS_REGION` または `AWS_DEFAULT_REGION`
- `DEFAULT_SIMILARITY_THRESHOLD`（未指定時 0.8）
- `GRPC_PORT`（未指定時 50051）

## Go バックエンド連携（次作業）
1) Go スタブ再生成（Buf）
```bash
# リポジトリルート
buf generate
```

2) gRPC クライアント呼び出し追加（`server/internal/interfaces/http/handler/ml_handler.go` など）
```go
resp, err := client.RecognizeImage(ctx, &pb.RecognizeImageRequest{
    ImageData: imgBytes,    // リクエストから取得
    Threshold: 0.8,         // 任意（省略時はサーバ既定値）。proto3 optional に対応
})
```

3) HTTP エンドポイント設計案
- `POST /api/ml/recognize`（`multipart/form-data` で画像、`threshold` はクエリまたはフォーム）
- 成功: `{ is_match, similarity_score }`

## S3 取り扱い
- 参照画像は `s3://<bucket>/<category>/<filename>` の構成を想定。
- 起動時に全キーを列挙しダウンロード → 特徴量を計算してメモリ常駐。
- IAM: ECS タスクロールに `s3:GetObject`, `s3:ListBucket` を付与する。

## 既知の制約・改善候補
- 類似度: ORB+BF(Hamming) と簡易正規化。要件は満たすが、精度チューニング余地あり（SIFT/FLANN など）。
- HealthCheck: 現状は常に healthy を返す。S3 ロード完了/例外の反映を推奨。
- TLS: 現状は平文 gRPC。社内ネットワーク限定運用前提。必要に応じてサーバ証明書を設定。
- メモリ: 参照画像が多い場合は特徴量の永続化/キャッシュ戦略（オンデマンド読み込み）を検討。
- ログ: CloudWatch JSON 形式。コンテキスト（処理時間/閾値）をすでに出力。

## デプロイ（次作業）
- Terraform/ECS:
  - microservice のポートは 50051（`microservice_container_port`）に統一、タスク環境変数 `GRPC_PORT` を追加。
  - Cloud Map を有効化（`microservice.<name_prefix>.local`）。API サービスへ `IMAGE_RECOGNITION_GRPC_ADDR` を注入。
  - 参照画像バケットは `S3_BUCKET_NAME`（Python 側は `REFERENCE_S3_BUCKET` も許容）。
- CI: `buf generate` → Lint/Test → コンテナビルド → ECR プッシュ → ECS 更新を GitHub Actions に追加済み。
  - ワークフロー: `.github/workflows/server-deploy.yml`, `.github/workflows/microservice-deploy.yml`
  - インフラ変更の検証: `.github/workflows/infra-plan.yml`（`terraform plan` 実行）

## トラブルシュート
- 「同じエラーで10分以上」迷ったら一旦タイムアウトし、以下を確認：
  - スタブ再生成の漏れ（proto 更新時）
  - `sys.path` に `app/gen` が入っているか（`app/server.py` で追加済み）
  - 参照画像未設定時は `is_match=false` が期待挙動
  - OpenCV が WebP 等を扱えない場合はホストのライブラリ依存を確認（Docker では問題なし）

## 連絡事項（今回の環境で遭遇した点）
- サンドボックス環境では `uv run` が一部キャッシュアクセスで拒否される事象があり、昇格実行で回避しました。通常開発環境では再現しない見込みです。

以上です。必要であれば Go ハンドラーの実装雛形追加や CI 設定の草案作成も対応可能です。
