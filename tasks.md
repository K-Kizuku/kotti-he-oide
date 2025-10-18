# 実装計画（進捗反映）

- [x] 1. プロジェクト基盤とProtocol Buffers設定
  - services/image_recognition ディレクトリと Python 環境（uv/mypy/pytest/ruff）を整備
  - schema/proto/image_recognition.proto にサービス定義を拡張（Hello/RecognizeImage/HealthCheck、threshold は optional）
  - Python スタブを生成し `services/image_recognition/app/gen` に配置（設計書の `src/proto_generated` は実装上 `app/gen` に読み替え）
  - _要件: 1.1, 4.1, 4.2_

- [x] 2. 基本設定とユーティリティの実装
  - `pyproject.toml` に依存とmypy/ruff設定を定義
  - CloudWatch JSON ログフォーマッタを `app/utils/logger.py` に実装
  - 環境変数管理を `app/utils/config.py` に実装（`REFERENCE_S3_BUCKET`/`S3_BUCKET_NAME` 両対応、既定しきい値）
  - _要件: 5.1, 5.2_

- [x] 3. データモデルとS3サービスの実装
  - `ReferenceImage`/`RecognitionResult` を型注釈付きで追加
  - `S3Service` でキー列挙・一括ダウンロード・例外時の中断/ログ出力を実装
  - _要件: 2.1, 2.2, 2.3_

- [x] 4. 画像処理エンジンの実装
  - `preprocess_image`（形式統一/縮小/グレースケール化）
  - ORB による特徴抽出と BFMatcher+比率テストでのマッチングを実装
  - _要件: 3.1, 3.2, 7.1, 7.4_

- [x] 5. 類似度計算サービスの実装
  - 類似度 = 良質マッチ数 / max(desc 長) を 0.0-1.0 に正規化
  - 閾値バリデーションとデフォルト適用を `ImageService` に実装
  - _要件: 4.1, 4.3, 7.2, 7.3_

- [x] 6. gRPCサーバーの実装
  - `app/server.py` に RPC 実装（Hello/RecognizeImage/HealthCheck）
  - S3 参照画像の起動時ロード（未設定時は警告ログで空参照）
  - _要件: 1.2, 1.3, 1.4, 5.3, 5.4_

- [x] 7. エラーハンドリングと例外処理の実装
  - 画像形式エラー `ImageFormatError` を実装
  - INVALID_ARGUMENT 等の gRPC ステータス返却（しきい値不正時）
  - JSON 様式の詳細ログ出力
  - _要件: 3.3, 3.4, 5.1_

- [x] 8. 単体テストの実装
  - 前処理/特徴抽出/類似度の最小テストを追加（`services/image_recognition/tests`）
  - _要件: 7.1, 7.2, 7.3_

- [ ] 9. 統合テストとE2Eテストの実装（今回スコープ外）
  - 依頼により今回は未実施。将来対応時は Docker Compose で Go→gRPC の結合を検証
  - _要件: 1.1, 1.2, 1.3, 1.4_

- [ ] 10. Docker化とデプロイメント設定（部分的に対応済み）
  - [x] Dockerfile 整備（uv 依存解決・proto 生成・gRPC 起動）
  - [x] ECS タスク定義/TF 変数整合（`GRPC_PORT` を追加、`microservice_container_port` を 50051 に統一）
  - [x] API→Microservice の名前解決に Cloud Map を導入（`microservice.<name_prefix>.local`）
  - [x] API タスクに `IMAGE_RECOGNITION_GRPC_ADDR` を注入
  - [x] CI での Buf 生成 + コンテナビルド/ECR プッシュ（サーバ/マイクロサービス）
  - _要件: 6.1, 6.2, 6.3_

 - [x] 11. Go 側 HTTP プロキシの実装（画像認識）
   - `POST /api/ml/recognize` を追加。`multipart/form-data`（`image` or `file`）および生バイトを受付。
   - クエリ/フォームの `threshold`（任意）に対応。proto3 optional に合わせて `*float32` を設定。
   - 実装ファイル: `server/internal/interfaces/http/handler/ml_handler.go`、ルーティング: `server/main.go`。
   - Go/Python の gRPC スタブを `buf generate` で生成済み（Go: `server/internal/gen/...`、Python: `services/image_recognition/app/gen/...`）。

## 備考
- Go 側のスタブ更新は未実施（現状 Hello のみ利用）。`buf generate` を実行してから `ml_handler.go` に RecognizeImage 呼び出しを追加する。
- ローカル検証ではサンドボックス制約により一部コマンドで権限昇格が必要でしたが、通常開発環境では `uv sync && make proto && make run` で起動可能です。
