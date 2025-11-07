"""
gRPC サーバーのエントリポイント。

service: ImageRecognitionService
rpc:
  - Hello(HelloRequest) returns (HelloReply)
  - RecognizeImage(RecognizeImageRequest) returns (RecognizeImageResponse)
  - HealthCheck(HealthCheckRequest) returns (HealthCheckResponse)

起動例:
    uv run python -m app.server
"""

from __future__ import annotations

import asyncio
import logging
import os

import grpc

from pathlib import Path
import sys

# 生成物（app/gen）を import path に追加
_GEN_PATH = Path(__file__).resolve().parent / "gen"
if str(_GEN_PATH) not in sys.path:
    sys.path.append(str(_GEN_PATH))

from image_recognition.v1 import image_recognition_pb2 as pb2  # type: ignore
from image_recognition.v1 import image_recognition_pb2_grpc as pb2_grpc  # type: ignore
from app.services.image_service import ImageService
from app.services.s3_service import S3Service
from app.utils.config import AppConfig
from app.utils.logger import configure_logging


class ImageRecognitionServiceRPC(pb2_grpc.ImageRecognitionServiceServicer):
    def __init__(self, svc: ImageService) -> None:
        self._svc = svc
        self._logger = logging.getLogger("image_recognition.rpc")

    async def Hello(self, request: pb2.HelloRequest, context: grpc.aio.ServicerContext) -> pb2.HelloReply:  # type: ignore[override]
        name = request.name or "world"
        return pb2.HelloReply(message=f"hello {name}")

    async def RecognizeImage(
        self, request: pb2.RecognizeImageRequest, context: grpc.aio.ServicerContext
    ) -> pb2.RecognizeImageResponse:  # type: ignore[override]
        # optional にしたため、未指定かどうかを presence で判定
        try:
            has_threshold = request.HasField("threshold")
        except Exception:
            has_threshold = False
        threshold = request.threshold if has_threshold else None
        # place_id は将来の互換のため存在確認してから参照
        place_id: str | None = None
        try:
            has_place = hasattr(request, "place_id") and request.HasField("place_id")  # type: ignore[attr-defined]
        except Exception:
            has_place = False
        if has_place:
            place_id = str(request.place_id)  # type: ignore[attr-defined]
        try:
            result = self._svc.recognize_image(bytes(request.image_data), threshold, place_id=place_id)
            return pb2.RecognizeImageResponse(
                is_match=result.is_match,
                similarity_score=result.similarity_score,
                error_message=result.error_message or "",
            )
        except ValueError as e:
            # しきい値などの入力値エラー
            await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(e))

    async def HealthCheck(
        self, request: pb2.HealthCheckRequest, context: grpc.aio.ServicerContext
    ) -> pb2.HealthCheckResponse:  # type: ignore[override]
        # サーバープロセスが応答可能であることのみを返す
        return pb2.HealthCheckResponse(healthy=True, status="ok")


async def serve() -> None:
    configure_logging()
    logger = logging.getLogger("image_recognition")

    # 設定ロードと依存初期化
    cfg = AppConfig.load()
    port = cfg.grpc_port
    s3 = S3Service(cfg) if cfg.s3_bucket else None
    svc = ImageService(cfg)
    if s3 is not None and cfg.s3_bucket:
        # S3 から参照画像読み込み（失敗してもサーバは起動を継続）
        try:
            refs = s3.download_reference_images()
            svc.set_references(refs)
            if len(svc.references) == 0:
                logger.error("no reference images loaded from S3", extra={"extra_fields": {"bucket": cfg.s3_bucket, "prefix": cfg.s3_prefix}})
        except Exception:
            logger.exception(
                "failed to load reference images from S3; server will start without references",
                extra={"extra_fields": {"bucket": cfg.s3_bucket, "prefix": cfg.s3_prefix}},
            )
    else:
        logger.error("reference images not configured (S3 disabled or bucket not set)")

    server = grpc.aio.server()  # asyncio ベースの gRPC サーバー
    pb2_grpc.add_ImageRecognitionServiceServicer_to_server(ImageRecognitionServiceRPC(svc), server)
    # TLS 設定（環境変数指定時のみ有効）
    cert_file = cfg.grpc_tls_cert_file
    key_file = cfg.grpc_tls_key_file
    if cert_file and key_file and os.path.exists(cert_file) and os.path.exists(key_file):
        with open(cert_file, "rb") as f:
            cert = f.read()
        with open(key_file, "rb") as f:
            key = f.read()
        credentials = grpc.ssl_server_credentials(((key, cert),))
        server.add_secure_port(f"0.0.0.0:{port}", credentials)
        logger.info(
            "gRPC server starting (TLS)",
            extra={"extra_fields": {"port": port, "tls": True}},
        )
    else:
        server.add_insecure_port(f"0.0.0.0:{port}")
        logger.info(
            "gRPC server starting (insecure)",
            extra={"extra_fields": {"port": port, "tls": False}},
        )
    await server.start()
    logger.info("gRPC server started")
    await server.wait_for_termination()


if __name__ == "__main__":
    asyncio.run(serve())
