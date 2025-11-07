"""
AgentCore Memory統合サービス
過去プレイヤー回答の保存・検索
"""
from typing import List, Optional, Dict, Any
from datetime import datetime, timezone

from bedrock_agentcore.memory import MemoryClient
from config.settings import settings
from utils.logger import logger
from utils.validator import InputValidator
from models.domain import PastAnswer, MemorySearchResult


class MemoryService:
    """AgentCore Memoryを使用した回答保存・検索サービス"""

    def __init__(self):
        """
        MemoryServiceの初期化
        """
        self.client = MemoryClient(region_name=settings.AWS_REGION)
        self.memory_id = settings.AGENTCORE_MEMORY_ID
        self.validator = InputValidator()

        logger.info("MemoryService initialized", extra={"memory_id": self.memory_id})

    async def store_answer(
        self,
        session_id: str,
        question_id: str,
        answer: str,
        metadata: Optional[Dict[str, Any]] = None
    ) -> bool:
        """
        プレイヤーの回答をMemoryに保存

        Args:
            session_id: セッションID
            question_id: 質問ID
            answer: 回答テキスト
            metadata: メタデータ（タイムスタンプ等）

        Returns:
            保存成功した場合True
        """
        try:
            # 個人情報チェック
            if self.validator.contains_personal_info(answer):
                logger.warning(
                    "Answer contains personal information, skipping storage",
                    extra={"session_id": session_id, "question_id": question_id}
                )
                return False

            # サニタイズ
            sanitized_answer = self.validator.sanitize(answer)

            # メタデータの準備
            event_metadata = metadata or {}
            event_metadata.update({
                "timestamp": datetime.now(timezone.utc).isoformat(),
                "scene": "s4",
                "anonymized": True,
                "question_id": question_id
            })

            # Memory IDが設定されていない場合はスキップ
            if not self.memory_id:
                logger.warning("AGENTCORE_MEMORY_ID not set, skipping storage")
                return False

            # AgentCore Memoryに保存
            self.client.create_event(
                memory_id=self.memory_id,
                actor_id=session_id,  # セッションIDのみ使用（匿名化）
                session_id=session_id,
                messages=[
                    (sanitized_answer, "USER"),
                ]
            )

            logger.info(
                "Answer stored successfully",
                extra={
                    "session_id": session_id,
                    "question_id": question_id,
                    "answer_length": len(sanitized_answer)
                }
            )
            return True

        except Exception as e:
            logger.error(
                "Failed to store answer",
                extra={
                    "session_id": session_id,
                    "question_id": question_id,
                    "error": str(e)
                },
                exc_info=True
            )
            return False

    async def search_past_answers(
        self,
        query: str,
        question_id: Optional[str] = None,
        limit: int = 10
    ) -> MemorySearchResult:
        """
        過去プレイヤー回答を検索

        Args:
            query: 検索クエリ（セマンティック検索）
            question_id: 質問ID（フィルタリング用）
            limit: 取得件数

        Returns:
            検索結果
        """
        try:
            # Memory IDが設定されていない場合は空の結果を返す
            if not self.memory_id:
                logger.warning("AGENTCORE_MEMORY_ID not set, returning empty results")
                return MemorySearchResult(answers=[], total_count=0)

            # 名前空間とフィルターの構築
            namespace = "/answers"
            if question_id:
                namespace = f"/answers/{question_id}"

            # AgentCore Memoryから検索
            results = self.client.retrieve_memories(
                memory_id=self.memory_id,
                namespace=namespace,
                query=query
            )

            # 結果を PastAnswer に変換
            past_answers: List[PastAnswer] = []
            for result in results[:limit]:
                # メモリ結果の構造に応じて調整が必要
                # ここでは基本的な構造を想定
                past_answers.append(
                    PastAnswer(
                        answer_text=result.get("content", ""),
                        question_id=question_id or "unknown",
                        timestamp=datetime.now(timezone.utc),  # 実際のタイムスタンプを使用
                        anonymized=True
                    )
                )

            logger.info(
                "Past answers retrieved",
                extra={
                    "query": query,
                    "question_id": question_id,
                    "results_count": len(past_answers)
                }
            )

            return MemorySearchResult(
                answers=past_answers,
                total_count=len(past_answers)
            )

        except Exception as e:
            logger.error(
                "Failed to search past answers",
                extra={
                    "query": query,
                    "question_id": question_id,
                    "error": str(e)
                },
                exc_info=True
            )
            # エラー時は空の結果を返す
            return MemorySearchResult(answers=[], total_count=0)

    async def create_memory_if_not_exists(self, memory_name: str = "GameAgentMemory") -> Optional[str]:
        """
        Memoryリソースが存在しない場合は作成

        Args:
            memory_name: Memory名

        Returns:
            Memory ID
        """
        try:
            # 既にMemory IDが設定されている場合はそれを返す
            if self.memory_id:
                return self.memory_id

            # 新しいMemoryを作成
            memory = self.client.create_memory(
                name=memory_name,
                description="Memory for storing anonymized player answers"
            )

            memory_id = memory.get("id")
            logger.info(
                "Memory created successfully",
                extra={"memory_id": memory_id, "memory_name": memory_name}
            )

            return memory_id

        except Exception as e:
            logger.error(
                "Failed to create memory",
                extra={"memory_name": memory_name, "error": str(e)},
                exc_info=True
            )
            return None
