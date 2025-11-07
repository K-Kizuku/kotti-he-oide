"""
レスポンスモデル定義
FastAPI エンドポイントで使用
"""
from typing import List, Dict, Any, Optional
from pydantic import BaseModel, Field


class QuizOption(BaseModel):
    """クイズの選択肢"""
    index: int = Field(..., description="選択肢のインデックス (0-3)")
    text: str = Field(..., description="選択肢のテキスト")
    is_correct: bool = Field(..., description="正解かどうか")
    source: str = Field(..., description="ソース (player_answer, player_other_answer, past_player, system_generic)")


class QuizGenerationResponse(BaseModel):
    """クイズ生成レスポンス"""
    quiz_id: str = Field(..., description="生成されたクイズのID")
    place_id: str = Field(..., description="場所ID")
    question_text: str = Field(..., description="質問文")
    options: List[QuizOption] = Field(..., description="4つの選択肢")
    answer_index: int = Field(..., description="正解のインデックス")


class ResponseValidationResponse(BaseModel):
    """回答バリデーションレスポンス"""
    is_valid: bool = Field(..., description="有効かどうか")
    reason: str = Field(..., description="判定理由")
    feedback: Optional[str] = Field(None, description="無効の場合のフィードバックメッセージ")
    confidence: float = Field(..., description="判定の信頼度 (0.0-1.0)")


class DialogueGenerationResponse(BaseModel):
    """対話生成レスポンス"""
    dialogue_text: str = Field(..., description="生成された対話テキスト")
    voice_text: str = Field(..., description="VOICEVOX用テキスト")
    estimated_duration_ms: int = Field(..., description="推定音声長（ミリ秒）")


class MessageVerificationResponse(BaseModel):
    """メッセージ検証レスポンス"""
    matched: bool = Field(..., description="一致したかどうか")
    similarity_score: float = Field(..., description="類似度スコア (0.0-1.0)")
    reason: str = Field(..., description="判定理由")
    hint: Optional[str] = Field(None, description="不一致の場合のヒント")


class MemoryStoreResponse(BaseModel):
    """Memory保存レスポンス"""
    success: bool = Field(..., description="保存成功したかどうか")
    memory_id: str = Field(..., description="保存されたMemoryのID")


class HealthCheckResponse(BaseModel):
    """ヘルスチェックレスポンス"""
    status: str = Field(..., description="ステータス (ok, degraded, error)")
    model: str = Field(..., description="使用中のモデルID")
    memory: str = Field(..., description="Memory接続状態 (ok, error)")
    timestamp: str = Field(..., description="チェック時刻")


class InvocationResponse(BaseModel):
    """
    統合レスポンスモデル
    AgentCore Runtime の /invocations エンドポイント用
    """
    output: Dict[str, Any] = Field(..., description="操作結果")
