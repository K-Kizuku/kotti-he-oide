"""
フォールバックサービス
AI Agent System障害時の代替コンテンツ提供
"""
from typing import Optional
import uuid

from config.prompts import FALLBACK_QUIZZES, FALLBACK_DIALOGUES
from models.responses import QuizGenerationResponse, QuizOption, DialogueGenerationResponse
from utils.logger import logger


class FallbackService:
    """フォールバック処理サービス"""

    @staticmethod
    def get_fallback_quiz(place_id: str) -> Optional[QuizGenerationResponse]:
        """
        固定クイズを返す

        Args:
            place_id: 場所ID

        Returns:
            固定クイズレスポンス
        """
        try:
            quiz_template = FALLBACK_QUIZZES.get(place_id)
            if not quiz_template:
                logger.warning(f"No fallback quiz found for place_id: {place_id}")
                return None

            # QuizOptionに変換
            options = [
                QuizOption(
                    index=i,
                    text=opt["text"],
                    is_correct=opt["is_correct"],
                    source="fallback"
                )
                for i, opt in enumerate(quiz_template["options"])
            ]

            # 正解のインデックスを見つける
            answer_index = next(
                (i for i, opt in enumerate(options) if opt.is_correct),
                0
            )

            response = QuizGenerationResponse(
                quiz_id=f"fallback_{place_id}_{uuid.uuid4().hex[:8]}",
                place_id=place_id,
                question_text=quiz_template["question"],
                options=options,
                answer_index=answer_index
            )

            logger.info(
                "Fallback quiz generated",
                extra={"place_id": place_id, "quiz_id": response.quiz_id}
            )

            return response

        except Exception as e:
            logger.error(
                "Failed to generate fallback quiz",
                extra={"place_id": place_id, "error": str(e)},
                exc_info=True
            )
            return None

    @staticmethod
    def get_fallback_dialogue(scene: str) -> Optional[DialogueGenerationResponse]:
        """
        固定対話を返す

        Args:
            scene: シーン識別子

        Returns:
            固定対話レスポンス
        """
        try:
            dialogue_text = FALLBACK_DIALOGUES.get(scene)
            if not dialogue_text:
                logger.warning(f"No fallback dialogue found for scene: {scene}")
                # デフォルトの対話を返す
                dialogue_text = "承知いたしました。"

            # 推定音声長を計算（1文字あたり約150ms）
            estimated_duration_ms = len(dialogue_text) * 150

            response = DialogueGenerationResponse(
                dialogue_text=dialogue_text,
                voice_text=dialogue_text,
                estimated_duration_ms=estimated_duration_ms
            )

            logger.info(
                "Fallback dialogue generated",
                extra={"scene": scene}
            )

            return response

        except Exception as e:
            logger.error(
                "Failed to generate fallback dialogue",
                extra={"scene": scene, "error": str(e)},
                exc_info=True
            )
            return None

    @staticmethod
    def is_fallback_enabled() -> bool:
        """
        フォールバックモードが有効かチェック

        Returns:
            有効な場合True
        """
        from config.settings import settings
        return settings.FALLBACK_MODE
