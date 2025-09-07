import numpy as np

from app.utils.image_processor import match_similarity


def test_similarity_bounds():
    # 空特徴同士は 0
    assert match_similarity(np.zeros((0, 32), dtype=np.uint8), np.zeros((0, 32), dtype=np.uint8)) == 0.0

    # 適当なバイナリ記述子（ランダム）
    a = np.random.randint(0, 256, size=(128, 32), dtype=np.uint8)
    b = a.copy()
    sim = match_similarity(a, b)
    assert 0.0 <= sim <= 1.0
