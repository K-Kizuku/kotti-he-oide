#!/bin/bash
# ローカルテスト用スクリプト

echo "=== AI Agent Service Local Testing ==="
echo ""

# 1. サーバー起動（バックグラウンド）
echo "1. Starting server..."
uv run python -m app.server &
SERVER_PID=$!
sleep 3

# 2. /ping エンドポイントテスト
echo ""
echo "2. Testing /ping endpoint..."
curl -s http://localhost:8080/ping | jq .

# 3. /health エンドポイントテスト
echo ""
echo "3. Testing /health endpoint..."
curl -s http://localhost:8080/health | jq .

# 4. /invocations エンドポイントテスト（バリデーション）
echo ""
echo "4. Testing /invocations endpoint (validate_response)..."
curl -s -X POST http://localhost:8080/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "operation": "validate_response",
    "payload": {
      "question": "小学生の頃、何に夢中になってた？",
      "answer": "サッカーに夢中でした"
    }
  }' | jq .

# 5. サーバー停止
echo ""
echo "5. Stopping server..."
kill $SERVER_PID

echo ""
echo "=== Testing completed ==="
