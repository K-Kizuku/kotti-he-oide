# protobuf スキーマ

このディレクトリにはマイクロサービス間通信用の Protocol Buffers（gRPC）スキーマを配置します。Buf v2 を用いて Go/Python のスタブを自動生成します。

- ルート: `schema/proto`
- サービス別/バージョン別にサブディレクトリを作成（例: `image_recognition/v1`）

## 前提
- Buf CLI（v2）: https://buf.build/docs/installation

## 生成（Buf v2）
リポジトリルートに `buf.yaml`（modules）と `buf.gen.yaml`（plugins/out）があります。

```bash
# リポジトリルートで実行
buf generate

# server ディレクトリからでも実行可（Makefile 経由）
cd server && make proto
```

出力先:
- Go: `server/internal/gen/...`
- Python: `services/image_recognition/app/gen/...`

補足:
- `services/image_recognition/Dockerfile` はコンテナ内で `grpcio-tools` により再度生成します（Buf 不要）。
