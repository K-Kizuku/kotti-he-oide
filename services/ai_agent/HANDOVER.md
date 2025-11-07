# AI Agentサービス 引き継ぎ書

**作成日**: 2025-11-08
**作成者**: Claude Code
**対象**: 後任開発者

## 目次

1. [プロジェクト概要](#プロジェクト概要)
2. [実装完了項目](#実装完了項目)
3. [未実装項目](#未実装項目)
4. [アーキテクチャ詳細](#アーキテクチャ詳細)
5. [セットアップ手順](#セットアップ手順)
6. [既知の問題と制約](#既知の問題と制約)
7. [次のステップ](#次のステップ)
8. [参考資料](#参考資料)

---

## プロジェクト概要

### 目的
「赤煉瓦文化館 〜こっちにおいで〜」ゲーム向けの動的コンテンツ生成AIエージェントサービス。Strands Agentsフレームワークを使用し、AWS Bedrock AgentCore Runtime上で動作する。

### 技術スタック
- **言語**: Python 3.12+
- **フレームワーク**: Strands Agents, FastAPI
- **LLM**: Amazon Bedrock Claude 4 Sonnet
- **メモリ**: AWS AgentCore Memory
- **パッケージ管理**: uv
- **デプロイ**: Docker (ARM64), AWS ECS Fargate

### ディレクトリ構造
```
services/ai_agent/
├── app/                    # アプリケーション層
│   ├── agent.py           # GameAgent（メインエージェント）
│   └── server.py          # FastAPI サーバー
├── config/                 # 設定
│   ├── prompts.py         # プロンプトテンプレート
│   └── settings.py        # 環境変数管理
├── models/                 # データモデル
│   ├── domain.py          # ドメインモデル
│   ├── requests.py        # リクエストモデル
│   └── responses.py       # レスポンスモデル
├── services/               # ビジネスロジック
│   ├── fallback_service.py    # フォールバック処理
│   └── memory_service.py      # AgentCore Memory統合
├── tools/                  # エージェントツール
│   ├── dialogue_generator.py  # 対話生成
│   ├── message_verifier.py    # メッセージ検証
│   ├── quiz_generator.py      # クイズ生成
│   └── response_validator.py  # 回答バリデーション
├── utils/                  # ユーティリティ
│   ├── filter.py          # コンテンツフィルタリング
│   ├── logger.py          # 構造化ロギング
│   └── validator.py       # 入力検証
├── tests/                  # テスト（最小限）
├── scripts/                # スクリプト
├── Dockerfile              # ARM64対応
├── Makefile                # 開発用コマンド
├── pyproject.toml          # 依存関係定義
└── README.md               # セットアップガイド
```

---

## 実装完了項目

### ✅ 1. プロジェクト基盤

#### 1.1 パッケージ管理
- [x] `pyproject.toml` - Python 3.12以上、uv使用
- [x] `uv.lock` - 依存関係ロックファイル（自動生成済み）
- [x] `.env.example` - 環境変数テンプレート
- [x] `.gitignore` - Python/uv用

**依存関係**:
```toml
dependencies = [
  "fastapi>=0.109.0",
  "uvicorn[standard]>=0.27.0",
  "pydantic>=2.5.0",
  "strands-agents>=0.1.0",
  "bedrock-agentcore>=0.1.0",
  "boto3>=1.34.0",
  "python-dotenv>=1.0.0",
  "httpx>=0.26.0",
  "aws-opentelemetry-distro>=0.10.1",
]
```

#### 1.2 開発環境
- [x] `Makefile` - setup, dev, run, lint, type, test コマンド
- [x] `.dockerignore` - Dockerビルド最適化
- [x] `Dockerfile` - ARM64対応、ADOT計装済み
- [x] `README.md` - セットアップ・デプロイ手順

#### 1.3 ドキュメント
- [x] `README.md` - 概要、セットアップ、API、デプロイ
- [x] `HANDOVER.md` - 本ファイル（引き継ぎ書）

**コマンド例**:
```bash
make setup    # 依存関係インストール
make dev      # 開発サーバー起動（ホットリロード）
make lint     # Ruffでコード品質チェック
```

---

### ✅ 2. 設定とユーティリティ

#### 2.1 設定管理 (`config/settings.py`)
- [x] 環境変数からの設定読み込み
- [x] AWS認証情報管理
- [x] Bedrock モデルID設定
- [x] AgentCore Memory ID設定
- [x] タイムアウト・キャッシュ設定
- [x] 設定バリデーション

**環境変数**:
```bash
AWS_REGION=ap-northeast-1
AWS_ACCESS_KEY_ID=<YOUR_KEY>
AWS_SECRET_ACCESS_KEY=<YOUR_SECRET>
BEDROCK_MODEL_ID=us.anthropic.claude-sonnet-4-20250514-v1:0
AGENTCORE_MEMORY_ID=<MEMORY_ID>
LOG_LEVEL=INFO
FALLBACK_MODE=enabled
```

#### 2.2 プロンプトテンプレート (`config/prompts.py`)
- [x] クイズ生成プロンプト
- [x] 回答バリデーションプロンプト
- [x] 対話生成プロンプト
- [x] メッセージ検証プロンプト
- [x] フォールバッククイズテンプレート（5箇所分）
- [x] フォールバック対話テンプレート

#### 2.3 ロギング (`utils/logger.py`)
- [x] JSON形式の構造化ログ
- [x] CloudWatch Logs対応
- [x] session_id、operation、duration_ms等の自動記録
- [x] エラー時のスタックトレース記録

**ログ形式例**:
```json
{
  "timestamp": "2025-11-08T00:00:00.000Z",
  "level": "INFO",
  "service": "ai-agent",
  "message": "Quiz generated successfully",
  "session_id": "uuid-v4",
  "operation": "generate_quiz",
  "duration_ms": 8500
}
```

#### 2.4 入力検証 (`utils/validator.py`)
- [x] テキスト長チェック（最大2000文字）
- [x] 許可文字パターン検証（日本語・英数字・記号）
- [x] HTMLタグ除去
- [x] 制御文字除去
- [x] 個人情報検出（メール、電話番号、郵便番号）

#### 2.5 コンテンツフィルタリング (`utils/filter.py`)
- [x] 不適切コンテンツ検出（基本パターン）
- [x] 安全なフォールバックテキスト提供

---

### ✅ 3. データモデル

#### 3.1 リクエストモデル (`models/requests.py`)
- [x] `PlayerAnswer` - プレイヤーの回答
- [x] `QuizGenerationRequest` - クイズ生成リクエスト
- [x] `ResponseValidationRequest` - 回答バリデーションリクエスト
- [x] `PlayerContext` - プレイヤー情報
- [x] `DialogueGenerationRequest` - 対話生成リクエスト
- [x] `MessageVerificationRequest` - メッセージ検証リクエスト
- [x] `MemoryStoreRequest` - Memory保存リクエスト
- [x] `InvocationRequest` - 統合リクエスト

#### 3.2 レスポンスモデル (`models/responses.py`)
- [x] `QuizOption` - クイズ選択肢
- [x] `QuizGenerationResponse` - クイズ生成レスポンス
- [x] `ResponseValidationResponse` - 回答バリデーションレスポンス
- [x] `DialogueGenerationResponse` - 対話生成レスポンス
- [x] `MessageVerificationResponse` - メッセージ検証レスポンス
- [x] `MemoryStoreResponse` - Memory保存レスポンス
- [x] `HealthCheckResponse` - ヘルスチェックレスポンス
- [x] `InvocationResponse` - 統合レスポンス

#### 3.3 ドメインモデル (`models/domain.py`)
- [x] `PastAnswer` - 過去プレイヤー回答
- [x] `MemorySearchResult` - Memory検索結果
- [x] `ValidationResult` - バリデーション結果
- [x] `VerificationResult` - 検証結果

---

### ✅ 4. AgentCore Memory統合

#### 4.1 MemoryService (`services/memory_service.py`)
- [x] MemoryClientの初期化
- [x] `store_answer()` - 回答保存（匿名化）
- [x] `search_past_answers()` - セマンティック検索
- [x] `create_memory_if_not_exists()` - Memory作成
- [x] 個人情報除外ロジック
- [x] エラーハンドリング（Memory障害時も継続）

**実装詳細**:
```python
# 回答保存時
self.client.create_event(
    memory_id=self.memory_id,
    actor_id=session_id,  # セッションIDのみ（匿名）
    session_id=session_id,
    messages=[(sanitized_answer, "USER")]
)

# 検索時
results = self.client.retrieve_memories(
    memory_id=self.memory_id,
    namespace="/answers/{question_id}",
    query=query
)
```

---

### ✅ 5. エージェントツール

#### 5.1 クイズ生成ツール (`tools/quiz_generator.py`)
- [x] Strands Agent統合
- [x] AgentCore Memory検索（過去回答取得）
- [x] LLMによるクイズ生成
- [x] 4択構成（正解、プレイヤー別回答、過去回答、汎用）
- [x] JSON レスポンスパース
- [x] エラーハンドリング

**生成ロジック**:
1. 過去プレイヤー回答をMemoryから検索
2. プロンプトにプレイヤー回答 + 過去回答を含める
3. LLMで4択クイズ生成
4. JSONパース → QuizGenerationResponse に変換

#### 5.2 回答バリデーションツール (`tools/response_validator.py`)
- [x] 無効回答パターン検出（正規表現）
- [x] LLMによる文脈的妥当性判定
- [x] フィードバックメッセージ生成
- [x] エラー時フェイルセーフ（デフォルト有効）

**無効パターン**:
```python
INVALID_PATTERNS = [
    r"^なし$", r"^特にない$", r"^わからない$",
    r"^思いつかない$", r"^ない$", r"^無し$",
    r"^分からない$", r"^わかりません$", r"^特になし$"
]
```

#### 5.3 対話生成ツール (`tools/dialogue_generator.py`)
- [x] プレイヤー情報の考慮
- [x] 1942年設定維持
- [x] 自然な話し言葉生成（VOICEVOX用）
- [x] 発話長制御（100文字以内）
- [x] 推定音声長計算

#### 5.4 メッセージ検証ツール (`tools/message_verifier.py`)
- [x] 完全一致チェック
- [x] 正規化比較（空白・改行統一）
- [x] LLMによる意味的同等性判定
- [x] 不一致時のヒント生成

**検証ステップ**:
1. 完全一致 → 即座に正解（score=1.0）
2. 正規化一致 → 正解（score=0.95）
3. LLM判定 → 意味的同等性チェック

---

### ✅ 6. メインエージェント

#### 6.1 GameAgent (`app/agent.py`)
- [x] Strands Agent初期化（デフォルトBedrock使用）
- [x] MemoryService統合
- [x] FallbackService統合
- [x] 全ツールの初期化と統合
- [x] シングルトンパターン
- [x] エラーハンドリングとフォールバック

**初期化**:
```python
self.agent = Agent()  # Strands Agent（自動でBedrock使用）
self.memory_service = MemoryService()
self.fallback_service = FallbackService()

# ツール初期化
self.quiz_generator = QuizGeneratorTool(self.agent, self.memory_service)
self.response_validator = ResponseValidatorTool(self.agent)
# ...
```

---

### ✅ 7. FastAPIアプリケーション

#### 7.1 エンドポイント (`app/server.py`)

**必須エンドポイント（AgentCore Runtime要件）**:
- [x] `GET /ping` - ヘルスチェック
- [x] `POST /invocations` - メインエントリーポイント

**追加エンドポイント**:
- [x] `GET /health` - 詳細ヘルスチェック

#### 7.2 操作ハンドラー
- [x] `handle_generate_quiz()` - クイズ生成
- [x] `handle_validate_response()` - 回答バリデーション
- [x] `handle_generate_dialogue()` - 対話生成
- [x] `handle_verify_message()` - メッセージ検証
- [x] `handle_store_answer()` - 回答保存

#### 7.3 ミドルウェア
- [x] 起動時GameAgent初期化
- [x] グローバル例外ハンドラー
- [x] エラーログ記録

**リクエスト形式**:
```json
{
  "operation": "generate_quiz",
  "payload": {
    "session_id": "uuid-v4",
    "place_id": "spiral_stairs",
    "player_answers": [...]
  }
}
```

**レスポンス形式**:
```json
{
  "output": {
    "quiz_id": "quiz_spiral_stairs_abc123",
    "place_id": "spiral_stairs",
    "question_text": "...",
    "options": [...],
    "answer_index": 0
  }
}
```

---

### ✅ 8. フォールバック機能

#### 8.1 FallbackService (`services/fallback_service.py`)
- [x] 固定クイズテンプレート（5箇所分）
- [x] 固定対話テンプレート
- [x] フォールバックモード設定
- [x] エラー時の安全な代替コンテンツ提供

**フォールバック動作**:
- AI Agent障害時でもゲーム継続可能
- 固定コンテンツで基本機能を維持
- 環境変数`FALLBACK_MODE=enabled`で制御

---

### ✅ 9. デプロイ準備

#### 9.1 Docker設定
- [x] ARM64対応Dockerfile
- [x] uv使用のマルチステージビルド
- [x] ADOT計装（OpenTelemetry）
- [x] ポート8080露出
- [x] `.dockerignore`でビルド最適化

**Dockerfile特徴**:
```dockerfile
FROM --platform=linux/arm64 ghcr.io/astral-sh/uv:python3.12-bookworm-slim
# ...
CMD ["opentelemetry-instrument", "uv", "run", "uvicorn", "app.server:app", "--host", "0.0.0.0", "--port", "8080"]
```

#### 9.2 ビルド・実行コマンド
```bash
# ローカルビルド
docker buildx build --platform linux/arm64 -t ai-agent:arm64 --load .

# ローカルテスト
docker run --platform linux/arm64 -p 8080:8080 --env-file .env ai-agent:arm64

# ECRプッシュ（要設定）
docker buildx build --platform linux/arm64 \
  -t <account-id>.dkr.ecr.ap-northeast-1.amazonaws.com/kotti-ai-agent:latest \
  --push .
```

---

### ✅ 10. セキュリティ

#### 10.1 実装済み
- [x] 入力テキスト長制限（2000文字）
- [x] 許可文字パターンチェック
- [x] HTMLタグ除去
- [x] 個人情報検出と除外
- [x] コンテンツフィルタリング
- [x] テキストサニタイゼーション

#### 10.2 AWS統合
- [x] IAMロールベース認証対応
- [x] 環境変数からの認証情報読み込み
- [x] セッションIDのみ使用（個人情報なし）

---

## 未実装項目

### ❌ 1. インフラストラクチャ

#### 1.1 Terraform（`infra/ai_agent.tf`）
**未実装理由**: インフラ担当が別途実装予定

**必要なリソース**:
- [ ] IAMロール（AgentCore Runtime用）
- [ ] IAMポリシー（Bedrock, AgentCore Memory）
- [ ] ECRリポジトリ（kotti-ai-agent）
- [ ] AgentCore Runtime デプロイ設定
- [ ] CloudWatch Logs ロググループ
- [ ] CloudWatch Alarms（エラー率、レイテンシ）

**参考実装**（設計書より）:
```hcl
# infra/ai_agent.tf

resource "aws_iam_role" "ai_agent_role" {
  name = "game-ai-agent-role"
  # ...
}

resource "aws_iam_role_policy" "ai_agent_bedrock_policy" {
  # Bedrock, Memory権限
}

resource "aws_ecr_repository" "ai_agent" {
  name = "kotti-ai-agent"
}
```

#### 1.2 AgentCore Memory初期セットアップ
**未実装理由**: AWS環境未設定

**必要な作業**:
- [ ] AgentCore Memory作成
  ```python
  memory = client.create_memory(
      name="GameAgentMemory",
      description="Memory for storing anonymized player answers"
  )
  ```
- [ ] Memory IDを環境変数に設定
- [ ] CloudWatch Transaction Search有効化（Observability用）

**参考コマンド**:
```python
# memory_service.pyに実装済みのヘルパー使用
memory_service = MemoryService()
memory_id = await memory_service.create_memory_if_not_exists("GameAgentMemory")
```

---

### ❌ 2. Go Server統合

#### 2.1 AIAgentService インターフェース
**未実装理由**: サーバー側の実装範囲外

**必要な実装** (`server/internal/domain/service/ai_agent_service.go`):
```go
package service

type AIAgentService interface {
    GenerateQuiz(ctx context.Context, req *GenerateQuizRequest) (*Quiz, error)
    ValidateResponse(ctx context.Context, question, answer string) (*ValidationResult, error)
    GenerateDialogue(ctx context.Context, scene string, playerCtx *PlayerContext) (string, error)
    VerifyMessage(ctx context.Context, original, reinput string) (*VerificationResult, error)
    StoreAnswer(ctx context.Context, sessionID, questionID, answer string) error
}
```

#### 2.2 AgentClient 実装
**未実装理由**: サーバー側の実装範囲外

**必要な実装** (`server/internal/infrastructure/ai/agent_client.go`):
```go
package ai

type AgentClient struct {
    baseURL    string
    httpClient *http.Client
    timeout    time.Duration
}

func NewAgentClient(baseURL string) *AgentClient {
    return &AgentClient{
        baseURL: baseURL,
        httpClient: &http.Client{Timeout: 15 * time.Second},
    }
}

func (c *AgentClient) GenerateQuiz(ctx context.Context, req *GenerateQuizRequest) (*Quiz, error) {
    // HTTP POST to /invocations
    // operation: "generate_quiz"
}
```

#### 2.3 フィーチャーフラグ
**未実装理由**: サーバー側の実装範囲外

**必要な実装**:
```go
type FeatureFlags struct {
    UseAIQuizGeneration      bool
    UseAIResponseValidation  bool
    UseAIDialogueGeneration  bool
    UseAIMessageVerification bool
}
```

---

### ❌ 3. テストとCI/CD

#### 3.1 ユニットテスト
**未実装理由**: PoCのためテスト最小限

**実装済み**:
- [x] 基本インポートテスト（`tests/test_api.py`）

**未実装**:
- [ ] 各ツールのユニットテスト
- [ ] MemoryServiceのモックテスト
- [ ] FastAPI エンドポイント統合テスト
- [ ] エラーハンドリングテスト

**テストカバレッジ目標**: 不要（PoCのため）

#### 3.2 統合テスト
**未実装理由**: AWS環境未設定

**必要なテスト**:
- [ ] AgentCore Memory接続テスト
- [ ] Bedrock API呼び出しテスト
- [ ] エンドツーエンドテスト（全操作）
- [ ] フォールバック動作テスト

#### 3.3 CI/CDパイプライン
**未実装理由**: インフラ未構築

**必要な設定** (`.github/workflows/ai-agent.yml`):
- [ ] Lintチェック（ruff）
- [ ] 型チェック（mypy）
- [ ] Dockerビルド
- [ ] ECRプッシュ
- [ ] ECS/AgentCore Runtimeデプロイ

---

### ❌ 4. パフォーマンス最適化

#### 4.1 キャッシング
**未実装理由**: 時間制約

**設計書に記載あり**:
- [ ] `ResponseCache` クラス（LLMレスポンスキャッシング）
- [ ] エンベディングキャッシュ
- [ ] Redis統合（オプション）

**実装案**:
```python
class ResponseCache:
    def __init__(self, ttl_seconds=3600):
        self.cache = {}
        self.ttl = ttl_seconds

    def get_cache_key(self, prompt: str, params: dict) -> str:
        content = f"{prompt}:{json.dumps(params, sort_keys=True)}"
        return hashlib.sha256(content.encode()).hexdigest()
```

#### 4.2 バッチ処理
**未実装理由**: 時間制約

**必要な実装**:
- [ ] S6クイズ5問の並列生成
  ```python
  tasks = [
      self.generate_quiz(session_id, player_answers, place_id)
      for place_id in places
  ]
  quizzes = await asyncio.gather(*tasks)
  ```

#### 4.3 プロンプト最適化
**未実装理由**: 初期実装完了、要チューニング

**今後の改善**:
- [ ] トークン数削減
- [ ] Few-shot examples追加
- [ ] プロンプトA/Bテスト

---

### ❌ 5. 監視とObservability

#### 5.1 CloudWatch Metricsダッシュボード
**未実装理由**: インフラ未構築

**必要なメトリクス**:
- [ ] InvocationCount
- [ ] InvocationLatency (P95 < 15s)
- [ ] ErrorRate (< 5%)
- [ ] TokenUsage
- [ ] MemorySearchLatency (P95 < 2s)
- [ ] FallbackRate (< 10%)

#### 5.2 アラート設定
**未実装理由**: インフラ未構築

**必要なアラーム**:
- [ ] エラー率 > 5%
- [ ] P95レイテンシ > 15s
- [ ] Memory検索失敗率 > 10%

**参考Terraform**:
```hcl
resource "aws_cloudwatch_metric_alarm" "ai_agent_error_rate" {
  alarm_name          = "ai-agent-high-error-rate"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 2
  metric_name         = "ErrorRate"
  namespace           = "AWS/Bedrock/AgentCore"
  threshold           = 5.0
}
```

#### 5.3 X-Rayトレーシング
**未実装理由**: ADOT計装のみ実装済み

**実装済み**:
- [x] ADOT（OpenTelemetry）計装
- [x] Dockerfile内で`opentelemetry-instrument`使用

**未実装**:
- [ ] カスタムスパン追加
- [ ] セッションID伝播（baggageへの設定）
  ```python
  from opentelemetry import baggage, context
  ctx = baggage.set_baggage("session.id", session_id)
  context.attach(ctx)
  ```

---

### ❌ 6. 高度な機能（将来拡張）

#### 6.1 マルチモーダル対応
**未実装理由**: MVP範囲外

**設計書に記載あり**:
- [ ] 画像を含めたクイズ生成
- [ ] Claude 4 Sonnetのマルチモーダル機能活用

#### 6.2 リアルタイム音声対話
**未実装理由**: MVP範囲外

**必要な統合**:
- [ ] VOICEVOX Server直接呼び出し
- [ ] 音声ストリーミング

#### 6.3 適応型難易度調整
**未実装理由**: MVP範囲外

**設計書に記載あり**:
- [ ] プレイヤーパフォーマンス分析
- [ ] 動的難易度調整

#### 6.4 感情分析
**未実装理由**: MVP範囲外

**設計書に記載あり**:
- [ ] 回答テキストからの感情分析
- [ ] ホラー演出強度の動的調整

---

### ❌ 7. ドキュメント

#### 7.1 API仕様書
**未実装理由**: 時間制約

**必要な作成**:
- [ ] OpenAPI (Swagger) 仕様書
- [ ] リクエスト/レスポンス例
- [ ] エラーコード一覧

#### 7.2 運用マニュアル
**未実装理由**: 本番運用前

**必要な作成**:
- [ ] デプロイ手順書
- [ ] トラブルシューティングガイド
- [ ] ロールバック手順
- [ ] スケーリング設定

---

## アーキテクチャ詳細

### データフロー

```
┌─────────────────┐
│ Frontend        │
│ (Next.js)       │
└────────┬────────┘
         │ HTTP/REST
         ▼
┌─────────────────┐
│ Go Server       │
│ (Backend)       │
└────────┬────────┘
         │ HTTP POST /invocations
         ▼
┌─────────────────────────────────────────────┐
│ AI Agent Service (FastAPI)                  │
│  ┌──────────────────────────────────────┐   │
│  │ GameAgent (Strands Agents)           │   │
│  │  ├─ QuizGeneratorTool                │   │
│  │  ├─ ResponseValidatorTool            │   │
│  │  ├─ DialogueGeneratorTool            │   │
│  │  └─ MessageVerifierTool              │   │
│  └────────┬─────────────────────────────┘   │
│           │                                  │
│  ┌────────▼──────────┐  ┌────────────────┐  │
│  │ MemoryService     │  │ FallbackService│  │
│  └────────┬──────────┘  └────────────────┘  │
└───────────┼──────────────────────────────────┘
            │
    ┌───────┴────────┐
    ▼                ▼
┌─────────┐    ┌──────────────┐
│ Bedrock │    │ AgentCore    │
│ Claude  │    │ Memory       │
└─────────┘    └──────────────┘
```

### エージェントライフサイクル

1. **起動時**:
   ```python
   # サーバー起動時にGameAgentをシングルトン初期化
   @app.on_event("startup")
   async def startup_event():
       _ = get_game_agent()
   ```

2. **リクエスト処理**:
   ```python
   # /invocations エンドポイント
   operation = body.get("operation")

   if operation == "generate_quiz":
       agent = get_game_agent()
       result = await agent.generate_quiz(...)
   ```

3. **LLM呼び出し**:
   ```python
   # Strands Agentによる処理
   result = self.agent(prompt)
   response_text = result.message
   ```

4. **Memory統合**:
   ```python
   # 保存
   await memory_service.store_answer(session_id, question_id, answer)

   # 検索
   results = await memory_service.search_past_answers(query)
   ```

---

## セットアップ手順

### 1. 前提条件

- Python 3.12以上
- uv (`curl -LsSf https://astral.sh/uv/install.sh | sh`)
- Docker（デプロイ時）
- AWS認証情報（Bedrock + AgentCore Memory権限）

### 2. ローカル開発環境セットアップ

```bash
# 1. リポジトリクローン
cd services/ai_agent

# 2. 依存関係インストール
make setup
# または
uv sync

# 3. 環境変数設定
cp .env.example .env
# .envを編集してAWS認証情報を設定

# 4. 開発サーバー起動
make dev
# または
uv run uvicorn app.server:app --host 0.0.0.0 --port 8080 --reload
```

### 3. 環境変数設定

`.env`ファイルを作成し、以下を設定:

```bash
# 必須
AWS_REGION=ap-northeast-1
AWS_ACCESS_KEY_ID=AKIA...
AWS_SECRET_ACCESS_KEY=...

# 推奨
BEDROCK_MODEL_ID=us.anthropic.claude-sonnet-4-20250514-v1:0
AGENTCORE_MEMORY_ID=mem_xxxxx  # Memory作成後に設定
LOG_LEVEL=INFO
FALLBACK_MODE=enabled

# オプション
AI_AGENT_TIMEOUT=15
QUIZ_GENERATION_TIMEOUT=20
```

### 4. 動作確認

```bash
# ヘルスチェック
curl http://localhost:8080/ping
# → {"status": "healthy"}

curl http://localhost:8080/health
# → {"status": "ok", "model": "...", "memory": "..."}

# 回答バリデーション
curl -X POST http://localhost:8080/invocations \
  -H "Content-Type: application/json" \
  -d '{
    "operation": "validate_response",
    "payload": {
      "question": "小学生の頃、何に夢中になってた？",
      "answer": "サッカーに夢中でした"
    }
  }'
```

### 5. Lintチェック

```bash
make lint
# または
uv run ruff check .
```

---

## 既知の問題と制約

### 1. 技術的制約

#### 1.1 Strands Agents
- **制約**: デフォルトでBedrockを使用、AWS認証が必須
- **影響**: ローカル開発でもAWS認証情報が必要
- **回避策**: なし（フレームワークの仕様）

#### 1.2 AgentCore Memory
- **制約**: Memory ID未設定でも動作するが、過去回答検索は無効
- **影響**: クイズのダミー選択肢が汎用回答のみになる
- **回避策**: Memoryを作成してIDを環境変数に設定

#### 1.3 ARM64専用
- **制約**: AgentCore Runtime要件でARM64必須
- **影響**: Dockerビルドに`--platform linux/arm64`必須
- **回避策**: Docker Buildxを使用

### 2. パフォーマンス

#### 2.1 LLM呼び出しレイテンシ
- **問題**: クイズ生成に8〜12秒かかる可能性
- **影響**: ユーザー体験に影響
- **対策**:
  - タイムアウト15秒に設定済み
  - フォールバック機能で継続可能
  - **今後**: キャッシング実装で改善

#### 2.2 Memory検索
- **問題**: セマンティック検索に2〜3秒
- **影響**: クイズ生成時間に加算
- **対策**: 並列実行で一部軽減
- **今後**: エンベディングキャッシュ

### 3. セキュリティ

#### 3.1 個人情報検出
- **問題**: 正規表現ベースのため完全ではない
- **影響**: 一部の個人情報が漏れる可能性
- **対策**: 基本的なパターン（メール、電話、郵便番号）は検出
- **今後**: より高度な検出ロジック実装

#### 3.2 コンテンツフィルタリング
- **問題**: 基本パターンのみ実装
- **影響**: 不適切コンテンツが生成される可能性
- **対策**: 最小限のフィルター実装済み
- **今後**: より包括的なフィルター追加

### 4. テスト

#### 4.1 テストカバレッジ
- **問題**: ユニットテスト最小限
- **影響**: バグ検出が遅れる可能性
- **対策**: PoCのためテスト作成不要
- **今後**: 本番化する場合はテスト追加

#### 4.2 統合テスト
- **問題**: AWS環境がないと実行不可
- **影響**: CI/CDでのテストが困難
- **対策**: モックを使ったテスト（未実装）
- **今後**: LocalStack等で擬似環境構築

### 5. ドキュメント

#### 5.1 API仕様書
- **問題**: OpenAPI仕様書なし
- **影響**: フロントエンド連携が困難
- **対策**: コード内のdocstringで代用
- **今後**: Swagger UI追加

---

## 次のステップ

### 優先度: 高（すぐに実施）

#### 1. AgentCore Memory作成
```python
# services/ai_agent/ で実行
from services.memory_service import MemoryService

memory_service = MemoryService()
memory_id = await memory_service.create_memory_if_not_exists("GameAgentMemory")
print(f"Memory ID: {memory_id}")

# .envに追加
# AGENTCORE_MEMORY_ID={memory_id}
```

#### 2. Dockerビルドと動作確認
```bash
# ARM64ビルド
docker buildx build --platform linux/arm64 -t ai-agent:arm64 --load .

# ローカルテスト
docker run --platform linux/arm64 -p 8080:8080 --env-file .env ai-agent:arm64

# 別ターミナルで
curl http://localhost:8080/ping
```

#### 3. Go Server統合開始
- `server/internal/domain/service/ai_agent_service.go` 作成
- `server/internal/infrastructure/ai/agent_client.go` 作成
- S6クイズ生成エンドポイントから統合開始

### 優先度: 中（1週間以内）

#### 4. Terraformインフラ構築
```bash
# infra/ai_agent.tf 作成
cd infra
terraform apply
```

必要なリソース:
- IAMロール + ポリシー
- ECRリポジトリ
- CloudWatch Logs

#### 5. ECRデプロイ
```bash
# ECRログイン
aws ecr get-login-password --region ap-northeast-1 | \
  docker login --username AWS --password-stdin <account-id>.dkr.ecr.ap-northeast-1.amazonaws.com

# プッシュ
docker buildx build --platform linux/arm64 \
  -t <account-id>.dkr.ecr.ap-northeast-1.amazonaws.com/kotti-ai-agent:latest \
  --push .
```

#### 6. AgentCore Runtimeデプロイ
```python
import boto3

client = boto3.client('bedrock-agentcore-control', region_name='ap-northeast-1')

response = client.create_agent_runtime(
    agentRuntimeName='kotti-ai-agent',
    agentRuntimeArtifact={
        'containerConfiguration': {
            'containerUri': '<account-id>.dkr.ecr.ap-northeast-1.amazonaws.com/kotti-ai-agent:latest'
        }
    },
    networkConfiguration={"networkMode": "PUBLIC"},
    roleArn='arn:aws:iam::<account-id>:role/AgentRuntimeRole'
)

print(f"Agent Runtime ARN: {response['agentRuntimeArn']}")
```

### 優先度: 低（時間があれば）

#### 7. パフォーマンス最適化
- ResponseCacheクラス実装
- バッチ処理（クイズ5問並列生成）
- プロンプトチューニング

#### 8. 監視設定
- CloudWatch Metricsダッシュボード作成
- アラーム設定
- X-Rayトレーシング強化

#### 9. ドキュメント整備
- OpenAPI仕様書
- 運用マニュアル
- トラブルシューティングガイド

---

## 参考資料

### 公式ドキュメント

#### Strands Agents
- 公式サイト: https://strandsagents.com/latest/documentation/docs/
- GitHub: https://github.com/strands-agents/sdk-python
- デプロイガイド: https://strandsagents.com/latest/documentation/docs/user-guide/deploy/deploy_to_bedrock_agentcore/

#### AWS Bedrock AgentCore
- Developer Guide: https://docs.aws.amazon.com/bedrock-agentcore/latest/devguide/
- Python SDK: https://github.com/aws/bedrock-agentcore-sdk-python
- Memory Guide: https://docs.aws.amazon.com/bedrock-agentcore/latest/devguide/agentcore-sdk-memory.html

#### Amazon Bedrock
- Claude 4 Sonnet: https://docs.anthropic.com/claude/docs/models-overview
- Bedrock API: https://docs.aws.amazon.com/bedrock/

### プロジェクトドキュメント

- 要件定義: `.kiro/specs/ai-agent-integration/requirements.md`
- 設計書: `.kiro/specs/ai-agent-integration/design.md`
- タスクリスト: `.kiro/specs/ai-agent-integration/tasks.md`
- CLAUDE.md: `/CLAUDE.md`（プロジェクト全体）

### コード参考

#### image_recognitionサービス
- `services/image_recognition/` - 同様のマイクロサービス構成
- pyproject.toml、Dockerfileの参考に

---

## トラブルシューティング

### よくある問題

#### 1. `ModuleNotFoundError: No module named 'strands'`
**原因**: 依存関係未インストール
**解決**:
```bash
uv sync
```

#### 2. AWS認証エラー
**原因**: AWS認証情報未設定
**解決**:
```bash
# .envファイル確認
cat .env

# または環境変数設定
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
export AWS_REGION=ap-northeast-1
```

#### 3. Memory検索が空
**原因**: AGENTCORE_MEMORY_ID未設定
**解決**:
1. Memoryを作成
2. `.env`に`AGENTCORE_MEMORY_ID`を追加
3. サーバー再起動

#### 4. LLMレスポンスのJSONパースエラー
**原因**: LLMが指定形式で返さない
**解決**: プロンプトを調整（`config/prompts.py`）

#### 5. Dockerビルドエラー
**原因**: ARM64指定忘れ
**解決**:
```bash
docker buildx build --platform linux/arm64 -t ai-agent:arm64 --load .
```

---

## 連絡先

### 質問・問題報告
- GitHub Issues: https://github.com/K-Kizuku/kotti-he-oide/issues
- プロジェクトチャンネル: （Slackなど、適宜追加）

### コードレビュー
- Pull Request作成時: 設計書準拠を確認
- Lintチェック必須: `make lint`

---

## 更新履歴

| 日付 | バージョン | 変更内容 | 作成者 |
|------|-----------|---------|--------|
| 2025-11-08 | 1.0.0 | 初版作成 | Claude Code |

---

**以上が引き継ぎ書です。不明点があれば遠慮なく質問してください。**
