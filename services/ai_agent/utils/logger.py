"""
ロギングモジュール
構造化ログをJSON形式で出力
"""
import logging
import json
import sys
from datetime import datetime, timezone
from typing import Any, Dict
from config.settings import settings


class StructuredFormatter(logging.Formatter):
    """JSON形式の構造化ログフォーマッター"""

    def format(self, record: logging.LogRecord) -> str:
        log_data: Dict[str, Any] = {
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "level": record.levelname,
            "service": "ai-agent",
            "message": record.getMessage(),
        }

        # 追加の属性を含める
        if hasattr(record, "session_id"):
            log_data["session_id"] = record.session_id
        if hasattr(record, "operation"):
            log_data["operation"] = record.operation
        if hasattr(record, "duration_ms"):
            log_data["duration_ms"] = record.duration_ms
        if hasattr(record, "token_usage"):
            log_data["token_usage"] = record.token_usage

        # エラー情報を含める
        if record.exc_info:
            log_data["exception"] = self.formatException(record.exc_info)
        if hasattr(record, "stack_trace"):
            log_data["stack_trace"] = record.stack_trace

        return json.dumps(log_data, ensure_ascii=False)


def setup_logger(name: str = "ai-agent") -> logging.Logger:
    """ロガーをセットアップ"""
    logger = logging.getLogger(name)
    logger.setLevel(getattr(logging, settings.LOG_LEVEL.upper()))

    # ハンドラーがすでに存在する場合は追加しない
    if logger.handlers:
        return logger

    # コンソールハンドラー
    handler = logging.StreamHandler(sys.stdout)
    handler.setFormatter(StructuredFormatter())
    logger.addHandler(handler)

    return logger


# グローバルロガーインスタンス
logger = setup_logger()
