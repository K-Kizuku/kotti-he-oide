# Implementation Plan

このタスクリストは、Strands Agentsフレームワークを使用したAI Agent Systemの実装計画です。各タスクは段階的に実装され、要件定義書と設計書に基づいています。

## タスク一覧

- [ ] 1. AI Agent基盤のセットアップ
  - `services/ai_agent/`ディレクトリの作成
  - pyproject.tomlの作成（image_recognitionを参考に、Python 3.12以上、uv使用）
  - .env.exampleの作成（AWS認証情報、Bedrock設定）
  - .gitignoreの作成
  - Makefileの作成（開発用コマンド）
  - README.mdの作成（セットアップ手順）
  - app/ディレクトリ構造の作成（__init__.py、main.py）
  - Strands Agentsのインストールと設定
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

- [ ] 2. HTTP APIエンドポイントの実装
  - FastAPIアプリケーションのセットアップ
  - ヘルスチェックエンドポイント（`/health`）の実装
  - リクエスト/レスポンスモデルの定義（Pydantic）
  - エラーハンドリングミドルウェアの実装
  - _Requirements: 1.1, 1.2, 1.3_

- [ ] 3. メインエージェントの実装
  - GameAgentクラスの実装（Strands Agentsベース）
  - Amazon Bedrock（Claude 4 Sonnet）の設定
  - ツールの登録とエージェント初期化
  - CloudWatchトレーシングの設定
  - _Requirements: 1.1, 1.2, 8.1, 8.2, 8.3, 8.4, 8.5_

- [ ] 4. AgentCore Memory統合
  - MemoryServiceクラスの実装
  - AgentCore Memory接続の設定
  - 回答保存機能の実装（store_answer）
  - セマンティック検索機能の実装（search_past_answers）
  - 個人情報除外とデータ暗号化の実装
  - _Requirements: 5.1, 5.2, 5.3, 5.7, 9.1, 9.2, 9.6_


- [ ] 5. クイズ生成機能の実装
  - QuizGeneratorToolクラスの実装
  - クイズ生成プロンプトテンプレートの作成
  - 4択クイズ生成ロジックの実装（正解、ダミー3つ）
  - 過去プレイヤー回答の活用（AgentCore Memory検索）
  - クイズ生成APIエンドポイント（`/api/v1/quiz/generate`）の実装
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 5.2, 5.3, 5.4_

- [ ] 6. 回答バリデーション機能の実装
  - ResponseValidatorToolクラスの実装
  - 無効回答パターン検出ロジックの実装
  - LLMによる文脈的妥当性判定の実装
  - フィードバックメッセージ生成の実装
  - 回答バリデーションAPIエンドポイント（`/api/v1/response/validate`）の実装
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

- [ ] 7. 対話生成機能の実装
  - DialogueGeneratorToolクラスの実装
  - 対話生成プロンプトテンプレートの作成（1942年設定維持）
  - VOICEVOX用の自然な話し言葉生成
  - 発話長制御（100文字以内）の実装
  - 対話生成APIエンドポイント（`/api/v1/dialogue/generate`）の実装
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_

- [ ] 8. メッセージ検証機能の実装
  - MessageVerifierToolクラスの実装
  - 完全一致チェックの実装
  - LLMによる意味的同等性判定の実装
  - 不一致時のヒント生成の実装
  - メッセージ検証APIエンドポイント（`/api/v1/message/verify`）の実装
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

- [ ] 9. 入力検証とセキュリティの実装
  - InputValidatorクラスの実装（最大2000文字制限）
  - テキストサニタイゼーション機能の実装
  - ContentFilterクラスの実装（不適切コンテンツ検出）
  - レート制限ミドルウェアの実装
  - 個人情報除外ロジックの実装
  - _Requirements: 9.1, 9.2, 9.4, 9.5_


- [ ] 10. エラーハンドリングとフォールバックの実装
  - FallbackStrategyクラスの実装
  - 固定クイズテンプレートの作成
  - 固定対話テンプレートの作成
  - タイムアウト処理の実装（15秒）
  - リトライロジックの実装（最大2回、指数バックオフ）
  - AgentCore Memory障害時のフォールバック処理
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6_

- [ ] 11. パフォーマンス最適化の実装
  - ResponseCacheクラスの実装（LLMレスポンスキャッシング）
  - バッチ処理の実装（5問クイズの並列生成）
  - プロンプト最適化（トークン数削減）
  - エンベディングキャッシュの実装
  - 非同期処理の最適化
  - _Requirements: 10.1, 10.2, 10.4, 10.6_

- [ ] 12. ロギングの実装
  - Pythonロギング設定の実装（app/utils/logger.py）
  - 構造化ログフォーマットの実装（JSON）
  - ログレベル設定（環境変数から読み込み）
  - エラーログとスタックトレースの記録
  - 各ツール実行時のログ出力
  - _Requirements: 8.1, 8.2, 8.3_

