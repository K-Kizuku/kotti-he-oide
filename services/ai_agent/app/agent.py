"""
メインエージェント
Strands Agentsを使用したゲームAIエージェント
"""
from strands import Agent

from config.settings import settings
from services.memory_service import MemoryService
from services.fallback_service import FallbackService
from tools.quiz_generator import QuizGeneratorTool
from tools.response_validator import ResponseValidatorTool
from tools.dialogue_generator import DialogueGeneratorTool
from tools.message_verifier import MessageVerifierTool
from utils.logger import logger


class GameAgent:
    """ゲーム用AIエージェント"""

    def __init__(self):
        """GameAgentの初期化"""
        try:
            # Strands Agentを初期化（デフォルトでBedrock使用）
            self.agent = Agent()

            # サービスの初期化
            self.memory_service = MemoryService()
            self.fallback_service = FallbackService()

            # ツールの初期化
            self.quiz_generator = QuizGeneratorTool(self.agent, self.memory_service)
            self.response_validator = ResponseValidatorTool(self.agent)
            self.dialogue_generator = DialogueGeneratorTool(self.agent)
            self.message_verifier = MessageVerifierTool(self.agent)

            logger.info(
                "GameAgent initialized successfully",
                extra={
                    "model": settings.BEDROCK_MODEL_ID,
                    "region": settings.AWS_REGION
                }
            )

        except Exception as e:
            logger.error(
                "Failed to initialize GameAgent",
                extra={"error": str(e)},
                exc_info=True
            )
            raise

    async def generate_quiz(self, session_id: str, player_answers: list, place_id: str):
        """クイズ生成"""
        try:
            return await self.quiz_generator.generate_quiz(
                session_id=session_id,
                player_answers=player_answers,
                place_id=place_id
            )
        except Exception as e:
            logger.error(f"Quiz generation failed: {e}")
            # フォールバック
            if self.fallback_service.is_fallback_enabled():
                return self.fallback_service.get_fallback_quiz(place_id)
            raise

    async def validate_response(self, question: str, answer: str):
        """回答バリデーション"""
        return await self.response_validator.validate(question, answer)

    async def generate_dialogue(self, scene: str, player_context):
        """対話生成"""
        try:
            return await self.dialogue_generator.generate_dialogue(scene, player_context)
        except Exception as e:
            logger.error(f"Dialogue generation failed: {e}")
            # フォールバック
            if self.fallback_service.is_fallback_enabled():
                return self.fallback_service.get_fallback_dialogue(scene)
            raise

    async def verify_message(self, original: str, reinput: str):
        """メッセージ検証"""
        return await self.message_verifier.verify_message(original, reinput)

    async def store_answer(self, session_id: str, question_id: str, answer: str, metadata: dict = None):
        """回答保存"""
        return await self.memory_service.store_answer(
            session_id=session_id,
            question_id=question_id,
            answer=answer,
            metadata=metadata
        )


# グローバルインスタンス（サーバー起動時に初期化）
_agent_instance = None


def get_game_agent() -> GameAgent:
    """GameAgentのシングルトンインスタンスを取得"""
    global _agent_instance
    if _agent_instance is None:
        _agent_instance = GameAgent()
    return _agent_instance
