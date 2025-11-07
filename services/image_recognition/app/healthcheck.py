from __future__ import annotations

import os
import sys
import time
import grpc

from pathlib import Path

# 生成スタブ
from image_recognition.v1 import image_recognition_pb2 as pb2  # type: ignore
from image_recognition.v1 import image_recognition_pb2_grpc as pb2_grpc  # type: ignore


def main() -> int:
    port = int(os.getenv("GRPC_PORT", "50051"))
    addr = f"127.0.0.1:{port}"

    cert_file = os.getenv("GRPC_TLS_CERT_FILE", "")
    key_file = os.getenv("GRPC_TLS_KEY_FILE", "")

    # TLS: サーバ証明書を信頼（自己署名想定）
    if cert_file and key_file and Path(cert_file).exists() and Path(key_file).exists():
        with open(cert_file, "rb") as f:
            root = f.read()
        credentials = grpc.ssl_channel_credentials(root_certificates=root)
        channel: grpc.Channel = grpc.secure_channel(addr, credentials)
    else:
        channel = grpc.insecure_channel(addr)

    stub = pb2_grpc.ImageRecognitionServiceStub(channel)
    try:
        # 短めのタイムアウト
        resp = stub.HealthCheck(pb2.HealthCheckRequest(), timeout=2.0)
        return 0 if resp.healthy else 1
    except Exception:
        return 1
    finally:
        # すこし待ってクリーンアップ
        time.sleep(0.01)


if __name__ == "__main__":
    sys.exit(main())

