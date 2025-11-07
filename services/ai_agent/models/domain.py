"""
ドメインモデル定義
内部データ構造
"""
from typing import Optional, List
from dataclasses import dataclass
from datetime import datetime


@dataclass
class PastAnswer:
    """過去プレイヤーの回答（匿名）"""
    answer_text: str
    question_id: str
    timestamp: datetime
    anonymized: bool = True


@dataclass
class MemorySearchResult:
    """Memory検索結果"""
    answers: List[PastAnswer]
    total_count: int


@dataclass
class ValidationResult:
    """バリデーション結果"""
    is_valid: bool
    reason: str
    feedback: Optional[str] = None
    confidence: float = 0.0


@dataclass
class VerificationResult:
    """検証結果"""
    matched: bool
    similarity_score: float
    reason: str
    hint: Optional[str] = None
