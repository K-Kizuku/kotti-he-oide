"""
設定管理モジュール。

環境変数を読み込み、S3 やしきい値等の設定を提供する。
本番環境では ECS のタスク定義から設定、ローカルでは .env などから設定する想定。
"""

from __future__ import annotations

import os
from dataclasses import dataclass


@dataclass(frozen=True)
class AppConfig:
    # S3 関連
    s3_bucket: str
    s3_prefix: str
    aws_region: str

    # gRPC サーバー
    grpc_port: int
    grpc_tls_cert_file: str
    grpc_tls_key_file: str

    # 類似度のデフォルト閾値（最新仕様: 0.5〜0.6 推奨 → 0.6 を既定値とする）
    default_threshold: float = 0.6

    @staticmethod
    def load() -> "AppConfig":
        """環境変数から設定をロードする。"""
        # インフラ側の環境変数と整合を取るため S3_BUCKET_NAME も許容
        bucket = os.getenv("REFERENCE_S3_BUCKET", os.getenv("S3_BUCKET_NAME", ""))
        prefix = os.getenv("REFERENCE_S3_PREFIX", "")
        region = os.getenv("AWS_REGION", os.getenv("AWS_DEFAULT_REGION", "ap-northeast-1"))

        # gRPC 設定
        try:
            grpc_port = int(os.getenv("GRPC_PORT", "50051"))
        except ValueError:
            grpc_port = 50051
        tls_cert = os.getenv("GRPC_TLS_CERT_FILE", "")
        tls_key = os.getenv("GRPC_TLS_KEY_FILE", "")

        # default_threshold は 0.0-1.0 の範囲でパース
        try:
            default_th = float(os.getenv("DEFAULT_SIMILARITY_THRESHOLD", "0.6"))
        except ValueError:
            default_th = 0.6
        default_th = min(max(default_th, 0.0), 1.0)

        return AppConfig(
            s3_bucket=bucket,
            s3_prefix=prefix,
            aws_region=region,
            grpc_port=grpc_port,
            grpc_tls_cert_file=tls_cert,
            grpc_tls_key_file=tls_key,
            default_threshold=default_th,
        )
