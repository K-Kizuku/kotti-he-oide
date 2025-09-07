"""
gRPC サーバーのエントリポイント。

service: ImageRecognitionService
rpc: Hello(HelloRequest) returns (HelloReply)

起動例:
    uv run python -m app.server
"""

import asyncio
import logging
import os
from concurrent import futures

import grpc

from app.gen.image_recognition.v1 import image_recognition_pb2 as pb2
from app.gen.image_recognition.v1 import image_recognition_pb2_grpc as pb2_grpc


class ImageRecognitionService(pb2_grpc.ImageRecognitionServiceServicer):
    async def Hello(self, request: pb2.HelloRequest, context: grpc.aio.ServicerContext) -> pb2.HelloReply:  # type: ignore[override]
        name = request.name or "world"
        return pb2.HelloReply(message=f"hello {name}")


async def serve() -> None:
    port = int(os.environ.get("GRPC_PORT", "50051"))
    server = grpc.aio.server()  # asyncio ベースの gRPC サーバー
    pb2_grpc.add_ImageRecognitionServiceServicer_to_server(ImageRecognitionService(), server)
    server.add_insecure_port(f"0.0.0.0:{port}")
    logging.info("gRPC server starting on :%d", port)
    await server.start()
    logging.info("gRPC server started")
    await server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    asyncio.run(serve())

