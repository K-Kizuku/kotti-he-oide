"""
place_id の語彙定義（サーバ/フロント共有用）。

S3 の参照画像は `REFERENCE_S3_PREFIX/{place_id}/...` のディレクトリに配置し、
image_recognition サービスはキーの最初のセグメントをカテゴリ（=place_id）とみなして判定対象を絞り込みます。
"""

from __future__ import annotations

from typing import Dict, List

# 機械可読な固定ID
PLACE_IDS: List[str] = [
    "spiral_stairs",         # 螺旋階段を見上げる高い天井
    "main_hall_fireplace",   # メインホールの暖炉のレンガ
    "back_entrance_hinge",   # 裏玄関の扉の蝶番
    "entrance_door",         # 入口エントランスの扉
    "upstairs_parlor_piano", # 階上応接室のピアノ
]

# 表示用ラベル（日本語）
PLACE_ID_LABELS_JA: Dict[str, str] = {
    "spiral_stairs": "螺旋階段を見上げる高い天井",
    "main_hall_fireplace": "メインホールの暖炉のレンガ",
    "back_entrance_hinge": "裏玄関の扉の蝶番",
    "entrance_door": "入口エントランスの扉",
    "upstairs_parlor_piano": "階上応接室のピアノ",
}

# 推奨S3配置（例）: s3://<bucket>/<prefix>/<place_id>/xxx.jpg
PLACE_ID_TO_S3_PREFIX: Dict[str, str] = {pid: pid for pid in PLACE_IDS}

