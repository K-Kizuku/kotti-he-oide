"""
FastAPIサーバー
AgentCore Runtime要件に準拠したHTTP APIを提供
"""
from fastapi import FastAPI, HTTPException, Request
from fastapi.responses import JSONResponse
from datetime import datetime, timezone
from typing import Dict, Any
import traceback

from app.agent import get_game_agent
from models.requests import (
    PlayerAnswer
)
from models.responses import (
    HealthCheckResponse
)
from config.settings import settings
from utils.logger import logger
from utils.validator import InputValidator
from utils.filter import ContentFilter

# FastAPIアプリケーションの初期化
app = FastAPI(
    title="AI Agent Service",
    version="0.1.0",
    description="Game AI Agent powered by Strands Agents and AWS Bedrock AgentCore"
)

# バリデーターとフィルターの初期化
validator = InputValidator()
content_filter = ContentFilter()


@app.on_event("startup")
async def startup_event():
    """起動時の初期化処理"""
    try:
        # GameAgentの初期化（シングルトンパターン）
        _ = get_game_agent()
        logger.info("Server started successfully")
    except Exception as e:
        logger.error(f"Failed to start server: {e}", exc_info=True)
        raise


@app.get("/ping")
async def ping() -> Dict[str, str]:
    """
    ヘルスチェックエンドポイント（AgentCore Runtime要件）
    """
    return {"status": "healthy"}


@app.get("/health")
async def health_check() -> HealthCheckResponse:
    """
    詳細ヘルスチェック
    """
    try:
        agent = get_game_agent()

        # Memory接続チェック
        memory_status = "ok" if agent.memory_service.memory_id else "not_configured"

        return HealthCheckResponse(
            status="ok",
            model=settings.BEDROCK_MODEL_ID,
            memory=memory_status,
            timestamp=datetime.now(timezone.utc).isoformat()
        )
    except Exception as e:
        logger.error(f"Health check failed: {e}", exc_info=True)
        return HealthCheckResponse(
            status="error",
            model=settings.BEDROCK_MODEL_ID,
            memory="error",
            timestamp=datetime.now(timezone.utc).isoformat()
        )


@app.post("/invocations")
async def invocations(request: Request) -> Dict[str, Any]:
    """
    メインエントリーポイント（AgentCore Runtime要件）
    すべてのAI Agent操作を処理
    """
    try:
        # リクエストボディの取得
        body = await request.json()
        operation = body.get("operation")

        if not operation:
            raise HTTPException(status_code=400, detail="Missing 'operation' field")

        # 操作ごとにルーティング
        if operation == "generate_quiz":
            return await handle_generate_quiz(body)
        elif operation == "validate_response":
            return await handle_validate_response(body)
        elif operation == "generate_dialogue":
            return await handle_generate_dialogue(body)
        elif operation == "verify_message":
            return await handle_verify_message(body)
        elif operation == "store_answer":
            return await handle_store_answer(body)
        else:
            raise HTTPException(status_code=400, detail=f"Unknown operation: {operation}")

    except HTTPException:
        raise
    except Exception as e:
        logger.error(
            "Invocation failed",
            extra={"error": str(e), "traceback": traceback.format_exc()},
            exc_info=True
        )
        raise HTTPException(status_code=500, detail=f"Internal server error: {str(e)}")


async def handle_generate_quiz(body: Dict[str, Any]) -> Dict[str, Any]:
    """クイズ生成の処理"""
    try:
        req_data = body.get("payload", {})
        session_id = req_data.get("session_id")
        place_id = req_data.get("place_id")
        player_answers_data = req_data.get("player_answers", [])

        # バリデーション
        if not session_id or not place_id:
            raise HTTPException(status_code=400, detail="Missing required fields")

        # PlayerAnswerオブジェクトに変換
        player_answers = [
            PlayerAnswer(**ans) for ans in player_answers_data
        ]

        # クイズ生成
        agent = get_game_agent()
        result = await agent.generate_quiz(session_id, player_answers, place_id)

        return {"output": result.model_dump()}

    except Exception as e:
        logger.error(f"Quiz generation failed: {e}", exc_info=True)
        raise


async def handle_validate_response(body: Dict[str, Any]) -> Dict[str, Any]:
    """回答バリデーションの処理"""
    try:
        req_data = body.get("payload", {})
        question = req_data.get("question")
        answer = req_data.get("answer")

        # バリデーション
        if not question or not answer:
            raise HTTPException(status_code=400, detail="Missing required fields")

        # 入力検証
        is_valid, error_msg = validator.validate_text(answer)
        if not is_valid:
            raise HTTPException(status_code=400, detail=error_msg)

        # 回答バリデーション
        agent = get_game_agent()
        result = await agent.validate_response(question, answer)

        return {
            "output": {
                "is_valid": result.is_valid,
                "reason": result.reason,
                "feedback": result.feedback,
                "confidence": result.confidence
            }
        }

    except Exception as e:
        logger.error(f"Response validation failed: {e}", exc_info=True)
        raise


async def handle_generate_dialogue(body: Dict[str, Any]) -> Dict[str, Any]:
    """対話生成の処理"""
    try:
        req_data = body.get("payload", {})
        scene = req_data.get("scene")
        player_context_data = req_data.get("player_context", {})

        # バリデーション
        if not scene:
            raise HTTPException(status_code=400, detail="Missing required fields")

        # PlayerContextオブジェクトに変換
        from models.requests import PlayerContext
        player_context = PlayerContext(**player_context_data)

        # 対話生成
        agent = get_game_agent()
        result = await agent.generate_dialogue(scene, player_context)

        # コンテンツフィルタリング
        if not content_filter.is_appropriate(result.dialogue_text):
            logger.warning("Inappropriate content detected in dialogue")
            result.dialogue_text = content_filter.get_safe_fallback()
            result.voice_text = content_filter.get_safe_fallback()

        return {"output": result.model_dump()}

    except Exception as e:
        logger.error(f"Dialogue generation failed: {e}", exc_info=True)
        raise


async def handle_verify_message(body: Dict[str, Any]) -> Dict[str, Any]:
    """メッセージ検証の処理"""
    try:
        req_data = body.get("payload", {})
        original = req_data.get("original")
        reinput = req_data.get("reinput")

        # バリデーション
        if not original or not reinput:
            raise HTTPException(status_code=400, detail="Missing required fields")

        # メッセージ検証
        agent = get_game_agent()
        result = await agent.verify_message(original, reinput)

        return {
            "output": {
                "matched": result.matched,
                "similarity_score": result.similarity_score,
                "reason": result.reason,
                "hint": result.hint
            }
        }

    except Exception as e:
        logger.error(f"Message verification failed: {e}", exc_info=True)
        raise


async def handle_store_answer(body: Dict[str, Any]) -> Dict[str, Any]:
    """回答保存の処理"""
    try:
        req_data = body.get("payload", {})
        session_id = req_data.get("session_id")
        question_id = req_data.get("question_id")
        answer = req_data.get("answer")
        metadata = req_data.get("metadata", {})

        # バリデーション
        if not session_id or not question_id or not answer:
            raise HTTPException(status_code=400, detail="Missing required fields")

        # 回答保存
        agent = get_game_agent()
        success = await agent.store_answer(session_id, question_id, answer, metadata)

        return {
            "output": {
                "success": success,
                "memory_id": agent.memory_service.memory_id or "not_configured"
            }
        }

    except Exception as e:
        logger.error(f"Answer storage failed: {e}", exc_info=True)
        raise


@app.exception_handler(Exception)
async def global_exception_handler(request: Request, exc: Exception):
    """グローバル例外ハンドラー"""
    logger.error(
        "Unhandled exception",
        extra={
            "path": request.url.path,
            "method": request.method,
            "error": str(exc),
            "traceback": traceback.format_exc()
        },
        exc_info=True
    )
    return JSONResponse(
        status_code=500,
        content={"detail": "Internal server error"}
    )


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
