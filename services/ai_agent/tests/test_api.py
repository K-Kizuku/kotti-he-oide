"""
基本的なAPI機能テスト
PoCのため最小限のテストのみ実装
"""
import pytest
from fastapi.testclient import TestClient


def test_import_server():
    """サーバーモジュールのインポートテスト"""
    try:
        from app import server
        assert server.app is not None
    except Exception as e:
        pytest.fail(f"Failed to import server: {e}")


def test_import_agent():
    """エージェントモジュールのインポートテスト"""
    try:
        from app import agent
        assert agent.GameAgent is not None
    except Exception as e:
        pytest.fail(f"Failed to import agent: {e}")


def test_import_models():
    """モデルのインポートテスト"""
    try:
        from models import requests, responses, domain
        assert requests.QuizGenerationRequest is not None
        assert responses.QuizGenerationResponse is not None
        assert domain.ValidationResult is not None
    except Exception as e:
        pytest.fail(f"Failed to import models: {e}")
