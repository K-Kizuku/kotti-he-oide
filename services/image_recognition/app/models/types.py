from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Dict, Optional

import numpy as np


@dataclass
class ReferenceImage:
    """参照画像データモデル。"""

    id: str
    category: str
    s3_key: str
    features: Optional[np.ndarray]
    metadata: Dict[str, Any]


@dataclass
class RecognitionResult:
    """認識結果データモデル。"""

    is_match: bool
    similarity_score: float
    processing_time: float
    error_message: Optional[str] = None

