"""
対話生成ツール
担当者との対話テキストを動的生成
"""
import json
import re
from typing import Dict, Any

from config.prompts import DIALOGUE_GENERATION_PROMPT
from models.requests import PlayerContext
from models.responses import DialogueGenerationResponse
from utils.logger import logger


class DialogueGeneratorTool:
    """対話テキスト生成を行うツール"""

    def __init__(self, agent):
        """
        Args:
            agent: Strands Agent インスタンス
        """
        self.agent = agent

    async def generate_dialogue(
        self,
        scene: str,
        player_context: PlayerContext
    ) -> DialogueGenerationResponse:
        """
        対話テキストを生成

        Args:
            scene: シーン識別子 (s1_greeting, s1_purpose, etc.)
            player_context: プレイヤー情報

        Returns:
            DialogueGenerationResponse: 生成された対話
        """
        try:
            # プレイヤー情報をフォーマット
            context_text = self._format_player_context(player_context)

            # プロンプト準備
            prompt = DIALOGUE_GENERATION_PROMPT.format(
                scene=scene,
                player_context=context_text
            )

            # LLMで対話生成
            result = self.agent(prompt)
            response_text = result.message if hasattr(result, 'message') else str(result)

            # JSONレスポンスをパース
            dialogue_data = self._parse_json_response(response_text)

            dialogue_response = DialogueGenerationResponse(
                dialogue_text=dialogue_data.get("dialogue_text", "承知いたしました。"),
                voice_text=dialogue_data.get("voice_text", dialogue_data.get("dialogue_text", "承知いたしました。")),
                estimated_duration_ms=dialogue_data.get("estimated_duration_ms", 3000)
            )

            logger.info(
                "Dialogue generated successfully",
                extra={"scene": scene}
            )

            return dialogue_response

        except Exception as e:
            logger.error(
                "Dialogue generation failed",
                extra={"scene": scene, "error": str(e)},
                exc_info=True
            )
            raise

    def _format_player_context(self, context: PlayerContext) -> str:
        """プレイヤー情報をフォーマット"""
        parts = []

        if context.arrival_method:
            parts.append(f"来館方法: {context.arrival_method}")

        if context.usual_activity:
            parts.append(f"普段の活動: {context.usual_activity}")

        if context.first_visit is not None:
            visit_status = "初回来館" if context.first_visit else "再来館"
            parts.append(f"来館状況: {visit_status}")

        return "\n".join(parts) if parts else "情報なし"

    def _parse_json_response(self, response_text: str) -> Dict[str, Any]:
        """LLMレスポンスからJSONをパース"""
        try:
            json_match = re.search(r'```json\s*(\{.*?\})\s*```', response_text, re.DOTALL)
            if json_match:
                json_str = json_match.group(1)
            else:
                json_match = re.search(r'(\{.*\})', response_text, re.DOTALL)
                json_str = json_match.group(1) if json_match else response_text

            return json.loads(json_str)

        except Exception as e:
            logger.warning(f"Failed to parse JSON response: {e}")
            # デフォルト値を返す
            return {
                "dialogue_text": "承知いたしました。",
                "voice_text": "承知いたしました。",
                "estimated_duration_ms": 3000
            }
