from __future__ import annotations

import io
import logging
from typing import List

import boto3
from botocore.exceptions import BotoCoreError, ClientError

from app.models.types import ReferenceImage
from app.utils.config import AppConfig


class S3Service:
    """S3 から参照画像を読み込み、メモリ上に保持するサービス。"""

    def __init__(self, config: AppConfig) -> None:
        self._config = config
        self._s3 = boto3.client("s3", region_name=config.aws_region)
        self._logger = logging.getLogger("image_recognition.s3")

    def list_reference_keys(self) -> List[str]:
        """参照画像のキー一覧を返す。prefix が空ならバケット直下。"""
        try:
            paginator = self._s3.get_paginator("list_objects_v2")
            keys: List[str] = []
            for page in paginator.paginate(Bucket=self._config.s3_bucket, Prefix=self._config.s3_prefix or None):
                for obj in page.get("Contents", []) or []:
                    k = obj.get("Key")
                    if isinstance(k, str) and not k.endswith("/"):
                        keys.append(k)
            return keys
        except (BotoCoreError, ClientError) as e:
            self._logger.exception("S3 list failed")
            raise RuntimeError(f"failed to list reference images: {e}") from e

    def download_reference_images(self) -> List[ReferenceImage]:
        """S3 から全参照画像をダウンロードし、`ReferenceImage` として返す。"""
        keys = self.list_reference_keys()
        results: List[ReferenceImage] = []
        for key in keys:
            try:
                buf = io.BytesIO()
                self._s3.download_fileobj(self._config.s3_bucket, key, buf)
                data = buf.getvalue()
                # カテゴリはキーのディレクトリ名（"category/xxx.jpg" → category）
                category = key.split("/")[0] if "/" in key else "default"
                results.append(
                    ReferenceImage(
                        id=key,
                        category=category,
                        s3_key=key,
                        features=None,  # 後段で特徴量を算出
                        metadata={"size": len(data), "key": key, "_raw": data},
                    )
                )
            except (BotoCoreError, ClientError) as e:
                self._logger.exception("S3 download failed: %s", key)
                raise RuntimeError(f"failed to download {key}: {e}") from e
        return results

