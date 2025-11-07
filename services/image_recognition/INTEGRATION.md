# 画像認識サービス 繋ぎ込みガイド（Goサーバ/フロント向け）

この文書は、Goバックエンドから Python gRPC 画像認識サービスを呼び出す際に必要な情報と、フロントで扱う `place_id` の語彙をまとめたものです。

## 接続情報（Goサーバ → Python gRPC）
- アドレス環境変数: `IMAGE_RECOGNITION_GRPC_ADDR`（例: `127.0.0.1:50051`）
- プロトコル: gRPC（開発はInsecure、本番はVPC内＋任意TLS）
- サービス/メソッド:
  - `image_recognition.v1.ImageRecognitionService/RecognizeImage`
  - `image_recognition.v1.ImageRecognitionService/HealthCheck`

## リクエスト/レスポンス
- Request (RecognizeImageRequest)
  - `image_data: bytes` 必須（JPEG/PNG/WebP/BMP）
  - `threshold: float` 任意（0.0–1.0, 省略時はサーバ既定0.6）
  - `place_id: string` 任意（下記語彙のいずれか）
- Response (RecognizeImageResponse)
  - `is_match: bool`（`similarity_score >= threshold`）
  - `similarity_score: float`（0.0–1.0）
  - `error_message: string`（異常時のみ）

## フォールバック仕様（重要）
- 次の条件では比較処理をスキップし、80%の確率で「同じ」と判定します。
  - 記述子が抽出できなかった（画像が特徴不足）
  - `place_id` 指定時に該当カテゴリの参照画像が0件
- ログには `recognize fallback_random` として出力されます。

## HealthCheck
- サーバプロセスの応答可否のみを返し、常に `healthy=true, status="ok"`。
- S3疎通・参照画像のロード成否は起動ログで確認してください。

## place_id 語彙（固定）
- `spiral_stairs` — 螺旋階段を見上げる高い天井
- `main_hall_fireplace` — メインホールの暖炉のレンガ
- `back_entrance_hinge` — 裏玄関の扉の蝶番
- `entrance_door` — 入口エントランスの扉
- `upstairs_parlor_piano` — 階上応接室のピアノ

> 補足: 参照画像は `s3://<bucket>/<prefix>/<place_id>/*.jpg` に配置。サービスはキー先頭のディレクトリ名でカテゴリ判定します。

## Go側の呼び出し例（既存）
- `server/internal/interfaces/http/handler/ml_handler.go` にHTTP→gRPCのプロキシ実装あり。
  - `POST /api/ml/recognize` にて、画像ボディ（またはmultipart）＋`threshold`任意。
  - place_id を付与する場合は、gRPCクライアントに `RecognizeImageRequest.PlaceId` を追加してください（proto更新後）。

## エラーハンドリング
- 画像デコードに失敗: INVALID_ARGUMENT
- しきい値範囲外: INVALID_ARGUMENT
- その他サーバ内部: gRPCエラー + `error_message` に詳細（ログも出力）

## TLS（任意）
- 環境変数 `GRPC_TLS_CERT_FILE`, `GRPC_TLS_KEY_FILE` を設定するとTLSで待受。
- その場合、クライアント側もTLSで接続してください（社内VPC内での利用を推奨）。

