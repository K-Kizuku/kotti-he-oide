from __future__ import annotations

import logging
from typing import Tuple

import cv2
import numpy as np

_logger = logging.getLogger("image_recognition.image_processor")


class ImageFormatError(ValueError):
    """未対応の形式・破損画像等の例外。"""


def preprocess_image(image_data: bytes, max_size: int = 512) -> np.ndarray:
    """画像バイト列を読み込み、RGB 相当の ndarray (H, W, 3) を返す。

    - 入力: JPEG/PNG/WebP/BMP 等（OpenCV が対応する形式）
    - リサイズ: 長辺が `max_size` を超える場合に縮小
    """
    if not image_data:
        raise ImageFormatError("empty image data")

    arr = np.frombuffer(image_data, dtype=np.uint8)
    bgr = cv2.imdecode(arr, cv2.IMREAD_COLOR)
    if bgr is None:
        raise ImageFormatError("decode failed or unsupported format")

    h, w = bgr.shape[:2]
    scale = float(max(h, w)) / float(max_size)
    if scale > 1.0:
        nh, nw = int(h / scale), int(w / scale)
        bgr = cv2.resize(bgr, (nw, nh), interpolation=cv2.INTER_AREA)

    # OpenCV は BGR なのでグレースケールへ（特徴抽出は大抵グレースケール使用）
    gray = cv2.cvtColor(bgr, cv2.COLOR_BGR2GRAY)
    return gray


def extract_features(gray_image: np.ndarray) -> Tuple[np.ndarray, np.ndarray]:
    """ORB による特徴点と記述子を抽出して返す。(keypoints, descriptors)

    SIFT は opencv-contrib が必要なため、実装は ORB をデフォルトとする。
    """
    if gray_image.ndim != 2:
        raise ValueError("gray image expected")

    orb = cv2.ORB_create(nfeatures=1000)
    kps, desc = orb.detectAndCompute(gray_image, None)
    if desc is None:
        # 特徴が検出できない場合は空配列
        desc = np.zeros((0, 32), dtype=np.uint8)
    return np.array(kps, dtype=object), desc


def match_similarity(desc1: np.ndarray, desc2: np.ndarray) -> float:
    """記述子同士のマッチングに基づく簡易類似度を 0.0-1.0 で返す。

    - Hamming 距離の BFMatcher を使用。
    - 比率テスト（Lowe's ratio test）で良質マッチのみ採用。
    - 類似度 = 良質マッチ数 / max(len(desc1), len(desc2))。
    """
    if desc1.size == 0 or desc2.size == 0:
        return 0.0

    bf = cv2.BFMatcher(cv2.NORM_HAMMING, crossCheck=False)
    matches = bf.knnMatch(desc1, desc2, k=2)

    good = []
    for m_n in matches:
        if len(m_n) < 2:
            continue
        m, n = m_n
        if m.distance < 0.75 * n.distance:
            good.append(m)

    denom = float(max(len(desc1), len(desc2)))
    sim = float(len(good)) / denom if denom > 0 else 0.0
    # 数値安定性のためクリップ
    return float(max(0.0, min(1.0, sim)))

