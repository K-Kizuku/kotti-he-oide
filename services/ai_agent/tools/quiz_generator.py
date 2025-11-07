"""
クイズ生成ツール
S4の回答を元にS6用の4択クイズを生成
"""
import json
import re
import uuid
from typing import List, Dict, Any

from config.prompts import QUIZ_GENERATION_PROMPT
from models.requests import PlayerAnswer
from models.responses import QuizGenerationResponse, QuizOption
from services.memory_service import MemoryService
from utils.logger import logger


class QuizGeneratorTool:
    """クイズ生成を行うツール"""

    PLACE_NAMES = {
        "spiral_stairs": "螺旋階段",
        "fireplace": "メインホールの暖炉",
        "hinge": "裏玄関の扉の蝶番",
        "entrance": "入口エントランスの扉",
        "piano": "階上応接室のピアノ"
    }

    def __init__(self, agent, memory_service: MemoryService):
        """
        Args:
            agent: Strands Agent インスタンス
            memory_service: MemoryService インスタンス
        """
        self.agent = agent
        self.memory_service = memory_service

    async def generate_quiz(
        self,
        session_id: str,
        player_answers: List[PlayerAnswer],
        place_id: str
    ) -> QuizGenerationResponse:
        """
        4択クイズを生成

        Args:
            session_id: セッションID
            player_answers: S4の10問の回答
            place_id: 場所ID

        Returns:
            QuizGenerationResponse: 生成されたクイズ
        """
        try:
            # 1. 過去プレイヤー回答をMemoryから検索
            place_name = self.PLACE_NAMES.get(place_id, place_id)
            memory_result = await self.memory_service.search_past_answers(
                query=f"場所:{place_name}",
                limit=10
            )

            # 2. プロンプト準備
            player_answers_text = self._format_player_answers(player_answers)
            past_answers_text = self._format_past_answers(memory_result.answers)

            prompt = QUIZ_GENERATION_PROMPT.format(
                player_answers=player_answers_text,
                place_name=place_name,
                past_answers=past_answers_text or "過去回答なし"
            )

            # 3. LLMでクイズ生成
            result = self.agent(prompt)
            response_text = result.message if hasattr(result, 'message') else str(result)

            # 4. JSONレスポンスをパース
            quiz_data = self._parse_json_response(response_text)

            # 5. QuizGenerationResponseに変換
            options = []
            for i, opt in enumerate(quiz_data.get("options", [])):
                options.append(
                    QuizOption(
                        index=i,
                        text=opt.get("text", ""),
                        is_correct=opt.get("is_correct", False),
                        source=opt.get("source", "unknown")
                    )
                )

            # 正解が設定されていることを確認
            if not any(opt.is_correct for opt in options):
                options[0].is_correct = True

            quiz_response = QuizGenerationResponse(
                quiz_id=f"quiz_{place_id}_{uuid.uuid4().hex[:8]}",
                place_id=place_id,
                question_text=quiz_data.get("question_text", ""),
                options=options,
                answer_index=next((i for i, opt in enumerate(options) if opt.is_correct), 0)
            )

            logger.info(
                "Quiz generated successfully",
                extra={
                    "session_id": session_id,
                    "place_id": place_id,
                    "quiz_id": quiz_response.quiz_id
                }
            )

            return quiz_response

        except Exception as e:
            logger.error(
                "Quiz generation failed",
                extra={"session_id": session_id, "place_id": place_id, "error": str(e)},
                exc_info=True
            )
            raise

    def _format_player_answers(self, answers: List[PlayerAnswer]) -> str:
        """プレイヤー回答をフォーマット"""
        formatted = []
        for ans in answers:
            formatted.append(f"Q{ans.question_id}: {ans.question_text}\nA: {ans.answer}")
        return "\n\n".join(formatted)

    def _format_past_answers(self, past_answers: List[Any]) -> str:
        """過去回答をフォーマット"""
        if not past_answers:
            return ""

        formatted = []
        for i, ans in enumerate(past_answers[:5], 1):  # 最大5件
            formatted.append(f"{i}. {ans.answer_text}")
        return "\n".join(formatted)

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
            raise ValueError("クイズ生成のJSONパースに失敗しました")
