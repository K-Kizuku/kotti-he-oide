from __future__ import annotations

import logging
import time
from typing import List

import numpy as np

from app.models.types import RecognitionResult, ReferenceImage
from app.utils.config import AppConfig
from app.utils.image_processor import (
    ImageFormatError,
    extract_features,
    match_similarity,
    preprocess_image,
)


class ImageService:
    """画像類似度計算のビジネスロジック。"""

    def __init__(self, config: AppConfig, reference_images: List[ReferenceImage] | None = None) -> None:
        self._config = config
        self._logger = logging.getLogger("image_recognition.image_service")
        self._refs: List[ReferenceImage] = reference_images or []

    @property
    def references(self) -> List[ReferenceImage]:
        return self._refs

    def set_references(self, refs: List[ReferenceImage]) -> None:
        """参照画像を設定し、必要なら特徴量を事前計算する。"""
        processed: List[ReferenceImage] = []
        for r in refs:
            data = r.metadata.get("_raw")
            if isinstance(data, (bytes, bytearray)):
                gray = preprocess_image(bytes(data))
                _, desc = extract_features(gray)
                processed.append(
                    ReferenceImage(
                        id=r.id,
                        category=r.category,
                        s3_key=r.s3_key,
                        features=desc,
                        metadata={k: v for k, v in r.metadata.items() if k != "_raw"},
                    )
                )
            else:
                # 特徴量が既に計算済みならそのまま
                processed.append(r)
        self._refs = processed
        self._logger.info("reference images prepared", extra={"extra_fields": {"count": len(self._refs)}})

    def recognize_image(self, image_data: bytes, threshold: float | None = None) -> RecognitionResult:
        start = time.perf_counter()
        th = self._normalize_threshold(threshold)
        try:
            gray = preprocess_image(image_data)
            _, desc = extract_features(gray)

            best = 0.0
            for ref in self._refs:
                if ref.features is None:
                    # 念のため遅延計算
                    self._logger.debug("lazy compute reference features", extra={"extra_fields": {"id": ref.id}})
                    continue
                sim = match_similarity(desc, ref.features)
                best = max(best, sim)

            result = RecognitionResult(is_match=best >= th, similarity_score=float(best), processing_time=time.perf_counter() - start)
            self._logger.info(
                "recognize done",
                extra={
                    "extra_fields": {
                        "score": result.similarity_score,
                        "is_match": result.is_match,
                        "threshold": th,
                        "refs": len(self._refs),
                        "ms": int(result.processing_time * 1000),
                    }
                },
            )
            return result
        except ImageFormatError as e:
            return RecognitionResult(
                is_match=False,
                similarity_score=0.0,
                processing_time=time.perf_counter() - start,
                error_message=str(e),
            )
        except Exception as e:  # 予期せぬエラー
            self._logger.exception("recognize failed")
            return RecognitionResult(
                is_match=False,
                similarity_score=0.0,
                processing_time=time.perf_counter() - start,
                error_message=f"internal error: {e}",
            )

    def _normalize_threshold(self, threshold: float | None) -> float:
        th = self._config.default_threshold if threshold is None else float(threshold)
        if th < 0.0 or th > 1.0:
            raise ValueError("threshold must be in 0.0-1.0")
        return th

