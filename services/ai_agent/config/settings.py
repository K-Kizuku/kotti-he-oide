"""
設定管理モジュール
環境変数から設定を読み込み、アプリケーション全体で使用する
"""
import os
from typing import Optional
from dotenv import load_dotenv

# 環境変数の読み込み
load_dotenv()


class Settings:
    """アプリケーション設定"""

    # AWS Configuration
    AWS_REGION: str = os.getenv("AWS_REGION", "ap-northeast-1")
    AWS_ACCESS_KEY_ID: Optional[str] = os.getenv("AWS_ACCESS_KEY_ID")
    AWS_SECRET_ACCESS_KEY: Optional[str] = os.getenv("AWS_SECRET_ACCESS_KEY")
    AWS_SESSION_TOKEN: Optional[str] = os.getenv("AWS_SESSION_TOKEN")

    # Bedrock Configuration
    BEDROCK_MODEL_ID: str = os.getenv(
        "BEDROCK_MODEL_ID",
        "us.anthropic.claude-sonnet-4-20250514-v1:0"
    )

    # AgentCore Memory Configuration
    AGENTCORE_MEMORY_ID: Optional[str] = os.getenv("AGENTCORE_MEMORY_ID")

    # Logging Configuration
    LOG_LEVEL: str = os.getenv("LOG_LEVEL", "INFO")

    # Application Configuration
    FALLBACK_MODE: bool = os.getenv("FALLBACK_MODE", "enabled").lower() == "enabled"

    # Timeout Configuration (seconds)
    AI_AGENT_TIMEOUT: int = int(os.getenv("AI_AGENT_TIMEOUT", "15"))
    QUIZ_GENERATION_TIMEOUT: int = int(os.getenv("QUIZ_GENERATION_TIMEOUT", "20"))

    # Input Validation
    MAX_TEXT_LENGTH: int = 2000
    MAX_QUIZ_OPTIONS: int = 4

    # Performance
    CACHE_TTL_SECONDS: int = 3600
    MAX_CONCURRENT_REQUESTS: int = 30

    @classmethod
    def validate(cls) -> None:
        """設定の妥当性をチェック"""
        if not cls.AWS_ACCESS_KEY_ID or not cls.AWS_SECRET_ACCESS_KEY:
            # CI/CD環境ではIAMロールを使用する可能性があるため、警告のみ
            print("Warning: AWS credentials not found in environment variables")

        if not cls.AGENTCORE_MEMORY_ID:
            print("Warning: AGENTCORE_MEMORY_ID not set - Memory features may not work")


# シングルトンインスタンス
settings = Settings()

# 起動時に設定を検証
settings.validate()
