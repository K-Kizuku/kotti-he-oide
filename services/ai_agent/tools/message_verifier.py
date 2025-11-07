"""
メッセージ検証ツール
S7でのメッセージ再入力を検証
"""
import json
import re
from typing import Dict, Any

from config.prompts import MESSAGE_VERIFICATION_PROMPT
from models.domain import VerificationResult
from utils.logger import logger


class MessageVerifierTool:
    """メッセージ再入力の検証を行うツール"""

    def __init__(self, agent):
        """
        Args:
            agent: Strands Agent インスタンス
        """
        self.agent = agent

    async def verify_message(self, original: str, reinput: str) -> VerificationResult:
        """
        メッセージの一致度を検証

        Args:
            original: 元のメッセージ（S4で入力）
            reinput: 再入力されたメッセージ（S7で入力）

        Returns:
            VerificationResult: 検証結果
        """
        try:
            # 完全一致チェック
            if original.strip() == reinput.strip():
                return VerificationResult(
                    matched=True,
                    similarity_score=1.0,
                    reason="完全一致",
                    hint=None
                )

            # 正規化して比較（空白、改行を統一）
            original_normalized = self._normalize_text(original)
            reinput_normalized = self._normalize_text(reinput)

            if original_normalized == reinput_normalized:
                return VerificationResult(
                    matched=True,
                    similarity_score=0.95,
                    reason="正規化後に一致",
                    hint=None
                )

            # LLMで意味的同等性を判定
            prompt = MESSAGE_VERIFICATION_PROMPT.format(
                original=original,
                reinput=reinput
            )

            # Strands Agentで処理
            result = self.agent(prompt)
            response_text = result.message if hasattr(result, 'message') else str(result)

            # JSONレスポンスをパース
            verification_data = self._parse_json_response(response_text)

            return VerificationResult(
                matched=verification_data.get("matched", False),
                similarity_score=verification_data.get("similarity_score", 0.0),
                reason=verification_data.get("reason", "LLM判定"),
                hint=verification_data.get("hint")
            )

        except Exception as e:
            logger.error(
                "Message verification failed",
                extra={"error": str(e)},
                exc_info=True
            )
            # エラー時は不一致と判定
            return VerificationResult(
                matched=False,
                similarity_score=0.0,
                reason="検証エラー",
                hint="もう一度正確に入力してください。"
            )

    def _normalize_text(self, text: str) -> str:
        """
        テキストを正規化

        Args:
            text: 正規化するテキスト

        Returns:
            正規化されたテキスト
        """
        # 空白と改行を統一
        normalized = re.sub(r'\s+', '', text)
        # 小文字化
        normalized = normalized.lower()
        return normalized

    def _parse_json_response(self, response_text: str) -> Dict[str, Any]:
        """
        LLMレスポンスからJSONを抽出してパース

        Args:
            response_text: LLMのレスポンステキスト

        Returns:
            パースされたJSON
        """
        try:
            # JSONブロックを探す
            json_match = re.search(r'```json\s*(\{.*?\})\s*```', response_text, re.DOTALL)
            if json_match:
                json_str = json_match.group(1)
            else:
                json_match = re.search(r'(\{.*\})', response_text, re.DOTALL)
                json_str = json_match.group(1) if json_match else response_text

            return json.loads(json_str)

        except Exception as e:
            logger.warning(f"Failed to parse JSON response: {e}")
            return {
                "matched": False,
                "similarity_score": 0.0,
                "reason": "JSONパースエラー"
            }
