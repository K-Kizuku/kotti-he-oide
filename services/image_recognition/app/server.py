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
        try:
            result = self._svc.recognize_image(bytes(request.image_data), threshold)
            return pb2.RecognizeImageResponse(
                is_match=result.is_match,
                similarity_score=result.similarity_score,
                error_message=result.error_message or "",
            )
        except ValueError as e:
            await context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(e))

    async def HealthCheck(
        self, request: pb2.HealthCheckRequest, context: grpc.aio.ServicerContext
    ) -> pb2.HealthCheckResponse:  # type: ignore[override]
        # ここでは簡易に healthy 固定（S3 などの疎通結果で将来拡張）
        return pb2.HealthCheckResponse(healthy=True, status="ok")


async def serve() -> None:
    port = int(os.environ.get("GRPC_PORT", "50051"))
    configure_logging()
    logger = logging.getLogger("image_recognition")

    # 設定ロードと依存初期化
    cfg = AppConfig.load()
    s3 = S3Service(cfg) if cfg.s3_bucket else None
    svc = ImageService(cfg)
    if s3 is not None and cfg.s3_bucket:
        # S3 から参照画像読み込み（失敗時は例外で起動中止）
        refs = s3.download_reference_images()
        svc.set_references(refs)
    else:
        logger.warning("S3 disabled or bucket not set; running with no references")

    server = grpc.aio.server()  # asyncio ベースの gRPC サーバー
    pb2_grpc.add_ImageRecognitionServiceServicer_to_server(ImageRecognitionServiceRPC(svc), server)
    server.add_insecure_port(f"0.0.0.0:{port}")
    logger.info("gRPC server starting", extra={"extra_fields": {"port": port}})
    await server.start()
    logger.info("gRPC server started")
    await server.wait_for_termination()


if __name__ == "__main__":
    asyncio.run(serve())
