# Requirements Document

## Introduction

本ドキュメントは、「赤煉瓦文化館 〜こっちにおいで〜」体験型Webホラーゲームに、Strands Agents フレームワークを用いたAI Agent機能を統合するための要件定義書です。

AI Agentは、ゲーム内の対話生成、クイズ生成、プレイヤー回答の分析、音声テキスト生成など、動的なコンテンツ生成を担当します。現在の固定テキストベースのシステムから、プレイヤーの回答に応じて適応的に反応するインテリジェントなシステムへと進化させることを目的とします。

### システムアーキテクチャ

データフロー：
```
フロントエンド → バックエンド(Go Server) → AI Agent(Strands Agents) → バックエンド
                                                      ↓
                                              (必要に応じて)
                                                      ↓
                                            VOICEVOX Server
```

AI Agentは **AWS Bedrock AgentCore Runtime** 上にデプロイされ、Strands Agentsフレームワークで実装されます。過去プレイヤーの回答管理には **AgentCore Memory機能** を使用します。

## Glossary

- **AI Agent System**: Strands Agentsフレームワークを使用して構築されたAIエージェントシステム。AWS Bedrock AgentCore Runtime上で動作し、ゲーム内の動的コンテンツ生成を担当
- **Game Server**: Go言語で実装されたバックエンドAPIサーバー。DDD + クリーンアーキテクチャで構築。AI Agent Systemとの通信を仲介
- **VOICEVOX Server**: EC2上で動作する音声合成サーバー（青山龍星(しっとり)）。Game ServerからHTTPリクエストで音声を生成
- **Session Context**: プレイヤーのゲームセッション情報（回答履歴、進行状況、タイムスタンプ等）
- **Quiz Generator**: S4の内省質問回答を元に、S6で使用する4択クイズを生成するAI機能
- **Response Validator**: プレイヤーの回答が「なし」「特にない」等の無効回答でないかを検証するAI機能
- **Dialogue Generator**: 担当者との対話テキストを動的に生成するAI機能
- **AWS Bedrock AgentCore Runtime**: AI Agent Systemのデプロイ先となるAWSマネージドサービス
- **AgentCore Memory**: AgentCore Runtimeの機能。過去プレイヤーの回答を保存・検索するために使用
- **Amazon Bedrock**: LLMモデルプロバイダー。Claude 4 Sonnetを使用
- **Community Tools**: strands-agents-toolsパッケージで提供されるコミュニティ駆動のツール群

## Requirements

### Requirement 1: AI Agent基盤の構築

**User Story:** システム管理者として、Strands Agentsフレームワークを使用したAI Agent基盤を構築し、ゲームサーバーから利用できるようにしたい。これにより、動的なコンテンツ生成機能を提供できる。

#### Acceptance Criteria

1. WHEN システムが起動するとき、THE AI Agent System SHALL AWS Bedrock AgentCore Runtime上で初期化され、Amazon Bedrock（Claude 4 Sonnet）をモデルプロバイダーとして使用する
2. THE AI Agent System SHALL Strands Agentsフレームワークを使用してPythonで実装される
3. THE AI Agent System SHALL Game ServerからHTTP/REST API経由でリクエストを受け取る
4. THE AI Agent System SHALL セッションコンテキストを受け取り、適切なプロンプトを生成してLLMに送信する
5. THE AI Agent System SHALL 環境変数（AWS_REGION、AWS_ACCESS_KEY_ID等）から認証情報を読み込む
6. WHERE エラーが発生した場合、THE AI Agent System SHALL 詳細なエラーログを出力し、適切なエラーレスポンスを返す

### Requirement 2: クイズ生成機能

**User Story:** ゲーム開発者として、プレイヤーのS4内省質問回答を元に、S6で使用する4択クイズを動的に生成したい。これにより、各プレイヤーに個別化された体験を提供できる。

#### Acceptance Criteria

1. WHEN S4の10問の回答がすべて完了したとき、THE AI Agent System SHALL 5つの場所それぞれに対応する4択クイズを生成する
2. THE AI Agent System SHALL 各クイズに以下を含める：質問文、4つの選択肢（正解1つ、ダミー3つ）、正解インデックス
3. THE AI Agent System SHALL 正解選択肢としてプレイヤー自身の回答を使用する
4. THE AI Agent System SHALL ダミー選択肢として、プレイヤーの別の回答、過去プレイヤーの匿名回答、システム汎用回答を使用する
5. THE AI Agent System SHALL 生成されたクイズをJSON形式でGame Serverに返す
6. THE AI Agent System SHALL クイズ生成時に、担当者との記憶のシンクロという文脈を維持する
7. THE AI Agent System SHALL 各クイズが明確で曖昧さがなく、プレイヤーが理解しやすい形式であることを保証する

### Requirement 3: 回答バリデーション機能

**User Story:** ゲーム開発者として、プレイヤーの回答が「なし」「特にない」「わからない」等の無効回答でないかをAIで判定したい。これにより、意味のある回答のみを受け付けることができる。

#### Acceptance Criteria

1. WHEN プレイヤーがS4の質問に回答したとき、THE AI Agent System SHALL 回答テキストが有効かどうかを判定する
2. THE AI Agent System SHALL 無効回答パターン（「なし」「特にない」「わからない」「思いつかない」等）を検出する
3. THE AI Agent System SHALL 判定結果として、有効/無効のブール値と、無効の場合の理由を返す
4. THE AI Agent System SHALL 文脈を考慮した判定を行い、質問内容に対して適切な回答かどうかを評価する
5. IF 回答が無効と判定された場合、THEN THE AI Agent System SHALL 「もう少し具体的に教えてください」等の適切なフィードバックメッセージを生成する

### Requirement 4: 対話テキスト生成機能

**User Story:** ゲーム開発者として、担当者との対話テキストをプレイヤーの回答に応じて動的に生成したい。これにより、より没入感のある体験を提供できる。

#### Acceptance Criteria

1. WHEN S1で担当者との対話が必要なとき、THE AI Agent System SHALL プレイヤーの入力（来館方法、普段の活動等）を考慮した対話テキストを生成する
2. THE AI Agent System SHALL 1942年の生命保険診査という設定を維持した対話を生成する
3. THE AI Agent System SHALL 生成されたテキストがVOICEVOX（青山龍星(しっとり)）で音声化されることを考慮し、自然な話し言葉で生成する
4. THE AI Agent System SHALL 対話の長さを適切に制御し、1つの発話が長すぎないようにする（目安：100文字以内）
5. THE AI Agent System SHALL ホラー要素を適度に含めつつ、プレイヤーを不快にさせない範囲で生成する

### Requirement 5: 過去プレイヤー回答の活用（AgentCore Memory使用）

**User Story:** ゲーム開発者として、過去プレイヤーの回答を匿名化してクイズのダミー選択肢として活用したい。これにより、「歴史を紡ぐ」というゲームコンセプトを実現できる。

#### Acceptance Criteria

1. WHEN プレイヤーがS4の質問に回答したとき、THE AI Agent System SHALL AgentCore Memory機能を使用して回答を匿名で保存する
2. WHEN クイズを生成するとき、THE AI Agent System SHALL AgentCore Memoryから過去プレイヤーの匿名回答を検索する
3. THE AI Agent System SHALL 検索時に、現在のクイズ文脈に関連する回答を優先的に取得する（セマンティック検索）
4. THE AI Agent System SHALL 選択された過去回答が正解と明確に区別できることを確認する
5. THE AI Agent System SHALL 過去回答が不足している場合、システム汎用のダミー回答を生成する
6. THE AI Agent System SHALL 個人を特定できる情報（名前、場所等）が含まれていないことを確認する
7. THE AI Agent System SHALL AgentCore Memoryに保存する際、セッションIDのみを識別子として使用し、個人情報を含めない

### Requirement 6: S7メッセージ検証機能

**User Story:** ゲーム開発者として、S7でプレイヤーがS4の回答を一字一句正確に再入力できたかをAIで判定したい。完全一致だけでなく、意味的に同等かどうかも判定できるようにしたい。

#### Acceptance Criteria

1. WHEN S7でプレイヤーが「人生の最期に達成したいこと」を再入力したとき、THE AI Agent System SHALL 元の回答と比較する
2. THE AI Agent System SHALL 完全一致の場合、即座に正解と判定する
3. THE AI Agent System SHALL 文字が若干異なる場合でも、意味が同等であれば正解と判定する
4. THE AI Agent System SHALL 判定結果として、一致度スコア（0.0〜1.0）と、一致/不一致のブール値を返す
5. IF 不一致と判定された場合、THEN THE AI Agent System SHALL どの部分が異なるかのヒントを生成する

### Requirement 7: エラーハンドリングとフォールバック

**User Story:** システム管理者として、AI Agent Systemが利用できない場合でも、ゲームが継続できるフォールバック機能を実装したい。これにより、システムの可用性を向上させる。

#### Acceptance Criteria

1. WHEN AI Agent Systemへの接続が失敗したとき、THE Game Server SHALL 事前定義された固定コンテンツを使用する
2. THE Game Server SHALL AI Agent Systemへのリクエストタイムアウトを15秒に設定する（AgentCore Runtimeの応答時間を考慮）
3. IF タイムアウトが発生した場合、THEN THE Game Server SHALL エラーログを記録し、フォールバックコンテンツを返す
4. THE Game Server SHALL AI Agent Systemの健全性を定期的にチェックする（ヘルスチェックエンドポイント）
5. THE AI Agent System SHALL ヘルスチェックリクエストに対して、ステータス、モデル情報、AgentCore Runtime接続状態を返す
6. WHERE AgentCore Memoryへのアクセスが失敗した場合、THE AI Agent System SHALL システム汎用のダミー回答のみを使用してクイズを生成する

### Requirement 8: 観測可能性とロギング

**User Story:** システム管理者として、AI Agent Systemの動作を監視し、問題を迅速に特定できるようにしたい。これにより、運用品質を向上させる。

#### Acceptance Criteria

1. THE AI Agent System SHALL すべてのリクエストとレスポンスをログに記録する
2. THE AI Agent System SHALL LLMへの各呼び出しに対して、トークン使用量、レイテンシ、コストを記録する
3. THE AI Agent System SHALL エラー発生時に、スタックトレースと関連するセッション情報を記録する
4. THE AI Agent System SHALL CloudWatch Logsにログを送信する
5. THE AI Agent System SHALL Strands Agentsの組み込みトレーシング機能を有効化する

### Requirement 9: セキュリティとデータ保護

**User Story:** システム管理者として、プレイヤーの個人情報を保護し、AIモデルへの不正なアクセスを防ぎたい。これにより、セキュアなシステムを構築できる。

#### Acceptance Criteria

1. THE AI Agent System SHALL プレイヤーの個人を特定できる情報（IP、端末ID等）をLLMに送信しない
2. THE AI Agent System SHALL セッションIDのみを識別子として使用する
3. THE AI Agent System SHALL AWS IAMロールベースの認証を使用してBedrock APIおよびAgentCore Runtimeにアクセスする
4. THE AI Agent System SHALL 入力テキストの最大長を制限し、過度に長い入力を拒否する（最大2000文字）
5. THE AI Agent System SHALL 生成されたコンテンツに不適切な内容が含まれていないかをフィルタリングする
6. THE AI Agent System SHALL AgentCore Memoryに保存されるデータを暗号化する
7. THE Game Server SHALL AI Agent Systemへのリクエストに認証トークンを含める

### Requirement 10: パフォーマンスとスケーラビリティ

**User Story:** システム管理者として、AI Agent Systemが複数の同時リクエストを効率的に処理できるようにしたい。これにより、最大30人の同時プレイヤーをサポートできる。

#### Acceptance Criteria

1. THE AI Agent System SHALL 単一のクイズ生成リクエストを12秒以内に処理する（AgentCore Runtimeのオーバーヘッドを考慮）
2. THE AI Agent System SHALL 単一の回答バリデーションリクエストを5秒以内に処理する
3. THE AI Agent System SHALL 最大30の同時リクエストを処理できる（AgentCore Runtimeのスケーリング機能を活用）
4. THE AI Agent System SHALL AgentCore Memoryへのクエリをキャッシュし、同一セッション内での重複検索を避ける
5. WHERE 負荷が高い場合、THE AI Agent System SHALL AgentCore Runtimeの自動スケーリングによって追加容量が確保される
6. THE AI Agent System SHALL 非同期処理を活用し、複数のLLM呼び出しを並列実行する
