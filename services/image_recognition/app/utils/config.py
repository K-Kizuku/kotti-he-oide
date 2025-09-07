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

    # 類似度のデフォルト閾値
    default_threshold: float = 0.8

    @staticmethod
    def load() -> "AppConfig":
        """環境変数から設定をロードする。"""
        # インフラ側の環境変数と整合を取るため S3_BUCKET_NAME も許容
        bucket = os.getenv("REFERENCE_S3_BUCKET", os.getenv("S3_BUCKET_NAME", ""))
        prefix = os.getenv("REFERENCE_S3_PREFIX", "")
        region = os.getenv("AWS_REGION", os.getenv("AWS_DEFAULT_REGION", "ap-northeast-1"))

        # default_threshold は 0.0-1.0 の範囲でパース
        try:
            default_th = float(os.getenv("DEFAULT_SIMILARITY_THRESHOLD", "0.8"))
        except ValueError:
            default_th = 0.8
        default_th = min(max(default_th, 0.0), 1.0)

        return AppConfig(
            s3_bucket=bucket,
            s3_prefix=prefix,
            aws_region=region,
            default_threshold=default_th,
        )
