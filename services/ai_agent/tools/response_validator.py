"""
回答バリデーションツール
プレイヤー回答の有効性を検証
"""
import json
import re
from typing import Dict, Any

from config.prompts import RESPONSE_VALIDATION_PROMPT
from models.domain import ValidationResult
from utils.logger import logger


class ResponseValidatorTool:
    """プレイヤー回答のバリデーションを行うツール"""

    # 無効回答パターン
    INVALID_PATTERNS = [
        r"^なし$",
        r"^特にない$",
        r"^わからない$",
        r"^思いつかない$",
        r"^ない$",
        r"^無し$",
        r"^分からない$",
        r"^わかりません$",
        r"^特になし$",
    ]

    def __init__(self, agent):
        """
        Args:
            agent: Strands Agent インスタンス
        """
        self.agent = agent

    async def validate(self, question: str, answer: str) -> ValidationResult:
        """
        回答の有効性を検証

        Args:
            question: 質問文
            answer: プレイヤーの回答

        Returns:
            ValidationResult: 検証結果
        """
        try:
            # まず無効パターンをチェック
            answer_normalized = answer.strip().lower()
            for pattern in self.INVALID_PATTERNS:
                if re.match(pattern, answer_normalized, re.IGNORECASE):
                    return ValidationResult(
                        is_valid=False,
                        reason="無効回答パターン検出",
                        feedback="もう少し具体的に教えてください。",
                        confidence=1.0
                    )

            # 空または短すぎる回答
            if len(answer.strip()) < 2:
                return ValidationResult(
                    is_valid=False,
                    reason="回答が短すぎます",
                    feedback="もう少し詳しく教えてください。",
                    confidence=1.0
                )

            # LLMで文脈的妥当性を判定
            prompt = RESPONSE_VALIDATION_PROMPT.format(
                question=question,
                answer=answer
            )

            # Strands Agentで処理
            result = self.agent(prompt)
            response_text = result.message if hasattr(result, 'message') else str(result)

            # JSONレスポンスをパース
            validation_data = self._parse_json_response(response_text)

            return ValidationResult(
                is_valid=validation_data.get("is_valid", True),
                reason=validation_data.get("reason", "LLM判定"),
                feedback=validation_data.get("feedback"),
                confidence=validation_data.get("confidence", 0.5)
            )

        except Exception as e:
            logger.error(
                "Validation failed",
                extra={"question": question, "error": str(e)},
                exc_info=True
            )
            # エラー時は有効と判定（フェイルセーフ）
            return ValidationResult(
                is_valid=True,
                reason="検証エラー - デフォルト有効",
                confidence=0.0
            )

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
                # JSON blocks without markdown
                json_match = re.search(r'(\{.*\})', response_text, re.DOTALL)
                json_str = json_match.group(1) if json_match else response_text

            return json.loads(json_str)

        except Exception as e:
            logger.warning(f"Failed to parse JSON response: {e}")
            # デフォルト値を返す
            return {
                "is_valid": True,
                "reason": "JSONパースエラー",
                "confidence": 0.0
            }
