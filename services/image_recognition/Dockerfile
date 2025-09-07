# ベースイメージ（できるだけ軽量な slim 系を使用）
FROM python:3.12-slim AS base

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    UV_LINK_MODE=copy \
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    PATH="/root/.local/bin:${PATH}"

WORKDIR /app

# uv をインストール（公式インストーラ）
RUN apt-get update \
 && apt-get install -y --no-install-recommends curl ca-certificates \
 && rm -rf /var/lib/apt/lists/* \
 && curl -LsSf https://astral.sh/uv/install.sh | sh

# 依存解決（レイヤキャッシュ最適化のためソースより先に pyproject をコピー）
# 依存解決（レイヤキャッシュ最適化のため先に pyproject のみ）
COPY services/image_recognition/pyproject.toml ./pyproject.toml
RUN uv sync --no-dev

# スキーマとアプリをコピー（ビルドコンテキストはリポジトリルートを想定）
COPY schema/proto ./schema/proto
COPY services/image_recognition/app ./app

# Python の gRPC スタブを生成
RUN mkdir -p app/gen && \
    python -m grpc_tools.protoc \
      -I schema/proto \
      --python_out app/gen \
      --grpc_python_out app/gen \
      schema/proto/image_recognition/v1/image_recognition.proto

EXPOSE 50051

# 本番実行コマンド（gRPC サーバ）
CMD ["uv", "run", "python", "-m", "app.server"]
