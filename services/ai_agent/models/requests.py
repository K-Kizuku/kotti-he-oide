"""
リクエストモデル定義
FastAPI エンドポイントで使用
"""
from typing import List, Dict, Any, Optional
from pydantic import BaseModel, Field


class PlayerAnswer(BaseModel):
    """プレイヤーの回答"""
    question_id: str = Field(..., description="質問ID")
    question_text: str = Field(..., description="質問文")
    answer: str = Field(..., description="回答テキスト")


class QuizGenerationRequest(BaseModel):
    """クイズ生成リクエスト"""
    session_id: str = Field(..., description="セッションID")
    place_id: str = Field(..., description="場所ID (spiral_stairs, fireplace, etc.)")
    player_answers: List[PlayerAnswer] = Field(..., description="S4の10問の回答")


class ResponseValidationRequest(BaseModel):
    """回答バリデーションリクエスト"""
    question: str = Field(..., description="質問文")
    answer: str = Field(..., description="プレイヤーの回答")


class PlayerContext(BaseModel):
    """プレイヤー情報"""
    arrival_method: Optional[str] = Field(None, description="来館方法")
    usual_activity: Optional[str] = Field(None, description="普段の活動")
    first_visit: Optional[bool] = Field(None, description="初回来館かどうか")


class DialogueGenerationRequest(BaseModel):
    """対話生成リクエスト"""
    scene: str = Field(..., description="シーン識別子 (s1_greeting, s1_purpose, etc.)")
    player_context: PlayerContext = Field(..., description="プレイヤー情報")


class MessageVerificationRequest(BaseModel):
    """メッセージ検証リクエスト"""
    original: str = Field(..., description="S4で入力された元のメッセージ")
    reinput: str = Field(..., description="S7で再入力されたメッセージ")


class MemoryStoreRequest(BaseModel):
    """Memory保存リクエスト"""
    session_id: str = Field(..., description="セッションID")
    question_id: str = Field(..., description="質問ID")
    answer: str = Field(..., description="回答テキスト")
    metadata: Dict[str, Any] = Field(default_factory=dict, description="メタデータ")


class InvocationRequest(BaseModel):
    """
    統合リクエストモデル
    AgentCore Runtime の /invocations エンドポイント用
    """
    operation: str = Field(..., description="操作タイプ (generate_quiz, validate_response, etc.)")
    payload: Dict[str, Any] = Field(..., description="操作固有のペイロード")
