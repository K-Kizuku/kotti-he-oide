import os
from pathlib import Path

import numpy as np

from app.utils.image_processor import ImageFormatError, extract_features, preprocess_image


def test_preprocess_raises_on_empty():
    try:
        preprocess_image(b"")
    except ImageFormatError:
        return
    assert False, "ImageFormatError expected"


def test_preprocess_and_extract_on_dummy_png(tmp_path: Path):
    # 100x100 の黒画像を PNG で作る（OpenCV 省略、raw で最小限）
    import cv2  # type: ignore

    arr = np.zeros((100, 100, 3), dtype=np.uint8)
    ok, buff = cv2.imencode(".png", arr)
    assert ok
    gray = preprocess_image(bytes(buff))
    _, desc = extract_features(gray)
    assert gray.ndim == 2
    assert desc.dtype == np.uint8
