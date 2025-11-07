"""
コンテンツフィルタリングモジュール
生成されたコンテンツの不適切性チェック
"""
import re
from typing import List


class ContentFilter:
    """生成されたコンテンツの不適切性チェック"""

    # 不適切なパターン（基本的なもののみ）
    INAPPROPRIATE_PATTERNS: List[str] = [
        r'暴力的',
        r'差別的',
        r'性的',
        r'攻撃的',
        # 必要に応じて追加
    ]

    @staticmethod
    def is_appropriate(text: str) -> bool:
        """
        コンテンツの適切性をチェック

        Args:
            text: チェックするテキスト

        Returns:
            適切な場合True
        """
        text_lower = text.lower()

        for pattern in ContentFilter.INAPPROPRIATE_PATTERNS:
            if re.search(pattern, text_lower):
                return False

        return True

    @staticmethod
    def get_safe_fallback() -> str:
        """
        安全なフォールバックコンテンツを返す

        Returns:
            安全なフォールバックテキスト
        """
        return "申し訳ございません。適切な応答を生成できませんでした。"
