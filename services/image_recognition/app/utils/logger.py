"""
CloudWatch Logs 向け JSON フォーマッタとロガー設定。
ローカルでも同じ形式で出力される。
"""

from __future__ import annotations

import json
import logging
from datetime import datetime, timezone
from typing import Any, Dict


class CloudWatchJSONFormatter(logging.Formatter):
    def format(self, record: logging.LogRecord) -> str:  # noqa: D401
        payload: Dict[str, Any] = {
            "timestamp": datetime.now(tz=timezone.utc).isoformat(),
            "level": record.levelname,
            "logger": record.name,
            "message": record.getMessage(),
            "module": record.module,
            "function": record.funcName,
            "line": record.lineno,
        }
        if hasattr(record, "extra_fields") and isinstance(record.extra_fields, dict):  # type: ignore[attr-defined]
            payload.update(record.extra_fields)  # type: ignore[attr-defined]
        if record.exc_info:
            payload["exception"] = self.formatException(record.exc_info)
        return json.dumps(payload, ensure_ascii=False)


def configure_logging(level: int = logging.INFO) -> None:
    handler = logging.StreamHandler()
    handler.setLevel(level)
    handler.setFormatter(CloudWatchJSONFormatter())

    root = logging.getLogger()
    root.setLevel(level)
    root.handlers = [handler]

