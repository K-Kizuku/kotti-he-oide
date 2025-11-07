# AI Agent Service

Strands AgentsフレームワークとAWS Bedrock AgentCore Runtimeを使用した、ゲーム用AIエージェントマイクロサービス。

## 概要

このサービスは、「赤煉瓦文化館 〜こっちにおいで〜」ゲームのための動的コンテンツ生成を担当します：

- **クイズ生成**: プレイヤーの回答を元に個別化された4択クイズを生成
- **回答バリデーション**: 無効回答の検出とフィードバック生成
- **対話生成**: 1942年設定を維持した担当者との対話テキスト生成
- **メッセージ検証**: S7での回答再入力の正確性判定
- **過去回答活用**: AgentCore Memoryを使用した匿名回答の保存・検索

## 技術スタック

- **Framework**: Strands Agents
- **Runtime**: AWS Bedrock AgentCore Runtime
- **LLM**: Amazon Bedrock Claude 4 Sonnet
- **Memory**: AgentCore Memory
- **API**: FastAPI
- **Language**: Python 3.12+
- **Package Manager**: uv

## セットアップ

### 前提条件

- Python 3.12以上
- uv (推奨) または pip
- AWS認証情報
- Docker (デプロイ時)

### インストール

```bash
# uvを使用 (推奨)
make setup

# または直接
uv sync
```

### 環境変数

`.env.example`をコピーして`.env`を作成し、必要な値を設定：

```bash
cp .env.example .env
# .envを編集してAWS認証情報等を設定
```

## 開発

### ローカル起動

```bash
# 開発モード (ホットリロード有効)
make dev

# または本番モード
make run
```

### テスト

```bash
# エンドポイントテスト
curl http://localhost:8080/ping

# クイズ生成テスト
curl -X POST http://localhost:8080/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "operation": "generate_quiz",
    "session_id": "test-session-id",
    "player_answers": [...]
  }'
```

### Linting

```bash
make lint
```

## API エンドポイント

### `/ping` (GET)
ヘルスチェック

### `/invocations` (POST)
メインエントリーポイント

#### リクエスト形式
```json
{
  "operation": "generate_quiz|validate_response|generate_dialogue|verify_message|store_answer",
  "session_id": "uuid-v4",
  ...operation-specific parameters
}
```

## デプロイ

### Docker ビルド

```bash
# ARM64イメージのビルド
docker buildx build --platform linux/arm64 -t ai-agent:arm64 --load .

# ローカルテスト
docker run --platform linux/arm64 -p 8080:8080 \
  --env-file .env \
  ai-agent:arm64
```

### AWS ECR へのプッシュ

```bash
# ECRリポジトリ作成
aws ecr create-repository --repository-name kotti-ai-agent --region ap-northeast-1

# ログイン
aws ecr get-login-password --region ap-northeast-1 | \
  docker login --username AWS --password-stdin <account-id>.dkr.ecr.ap-northeast-1.amazonaws.com

# ビルド & プッシュ
docker buildx build --platform linux/arm64 \
  -t <account-id>.dkr.ecr.ap-northeast-1.amazonaws.com/kotti-ai-agent:latest \
  --push .
```

### AgentCore Runtime デプロイ

```bash
# Terraformでインフラ構築
cd ../../infra
terraform apply

# または手動でboto3を使用
python deploy_agent.py
```

## アーキテクチャ

```
app/
├── server.py           # FastAPI application
├── agent.py            # Strands Agent initialization
├── models/             # Pydantic models
├── tools/              # Agent tools
├── services/           # Business logic services
├── utils/              # Utilities (logger, validator)
└── config/             # Configuration and prompts
```

## 監視

CloudWatchでメトリクスとログを確認：

```bash
# CloudWatch Logs
/aws/bedrock/agentcore/kotti-ai-agent

# Metrics
- InvocationCount
- InvocationLatency
- ErrorRate
- TokenUsage
```

## トラブルシューティング

### AWS認証エラー
```bash
aws configure
# または環境変数を設定
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
export AWS_REGION=ap-northeast-1
```

### タイムアウトエラー
- AgentCore Runtimeのタイムアウトは15秒
- 複雑なクイズ生成は最大20秒まで許容

### Memory検索エラー
- AGENTCORE_MEMORY_IDが正しく設定されているか確認
- IAMロールにMemory操作権限があるか確認

## 参考資料

- [Strands Agents Documentation](https://strandsagents.com/latest/documentation/docs/)
- [AWS Bedrock AgentCore](https://docs.aws.amazon.com/bedrock-agentcore/)
- [設計書](.kiro/specs/ai-agent-integration/design.md)
- [要件定義](.kiro/specs/ai-agent-integration/requirements.md)
