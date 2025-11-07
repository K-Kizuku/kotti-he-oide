"""
入力検証モジュール
プレイヤー入力のバリデーションとサニタイゼーション
"""
import re
from typing import Tuple
from config.settings import settings


class InputValidator:
    """入力データの検証とサニタイゼーション"""

    # 許可される文字パターン（日本語、英数字、基本的な記号）
    ALLOWED_CHARACTERS = r'^[\u3000-\u303f\u3040-\u309f\u30a0-\u30ff\uff00-\uffef\u4e00-\u9faf\u3400-\u4dbfa-zA-Z0-9\s.,!?、。！？\n\r]+$'

    @staticmethod
    def validate_text(text: str) -> Tuple[bool, str]:
        """
        テキストの基本的な検証

        Args:
            text: 検証するテキスト

        Returns:
            (is_valid, error_message): 検証結果とエラーメッセージ
        """
        if not text or text.strip() == "":
            return False, "テキストが空です"

        if len(text) > settings.MAX_TEXT_LENGTH:
            return False, f"テキストが長すぎます（最大{settings.MAX_TEXT_LENGTH}文字）"

        # 許可された文字のみかチェック
        if not re.match(InputValidator.ALLOWED_CHARACTERS, text):
            return False, "許可されていない文字が含まれています"

        return True, ""

    @staticmethod
    def sanitize(text: str) -> str:
        """
        テキストのサニタイゼーション

        Args:
            text: サニタイズするテキスト

        Returns:
            サニタイズされたテキスト
        """
        # HTMLタグの除去
        text = re.sub(r'<[^>]+>', '', text)

        # 制御文字の除去（改行、タブは保持）
        text = ''.join(
            char for char in text
            if ord(char) >= 32 or char in '\n\r\t'
        )

        # 連続する空白を1つに
        text = re.sub(r'\s+', ' ', text)

        return text.strip()

    @staticmethod
    def contains_personal_info(text: str) -> bool:
        """
        個人情報が含まれているかチェック

        Args:
            text: チェックするテキスト

        Returns:
            個人情報が含まれている場合True
        """
        # メールアドレスパターン
        if re.search(r'\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b', text):
            return True

        # 電話番号パターン（日本）
        if re.search(r'\d{2,4}-?\d{2,4}-?\d{4}', text):
            return True

        # 郵便番号パターン
        if re.search(r'\d{3}-?\d{4}', text):
            return True

        return False
