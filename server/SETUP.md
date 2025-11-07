# バックエンドセットアップガイド

## 実装内容

このバックエンドは、赤煉瓦文化館の体験型Webホラーゲーム用のREST APIサーバーです。

### アーキテクチャ

- **DDD（ドメイン駆動設計）+ Clean Architecture**
- **レイヤー構造**:
  - Domain Layer: ビジネスロジック、ドメインモデル
  - Application Layer: ユースケース
  - Infrastructure Layer: データベース接続、リポジトリ実装
  - Interface Layer: HTTPハンドラー、DTO

### 技術スタック

- **言語**: Go 1.25.1
- **データベース**: MySQL 8.0（go-sql-driver/mysql）
- **マイグレーション**: golang-migrate
- **ルーティング**: 標準ライブラリ（net/http）

## 環境変数

以下の環境変数を設定してください：

```bash
# 必須
DATABASE_URL="user:password@tcp(localhost:3306)/kotti_game?parseTime=true"

# VAPID鍵（Web Push機能用）
VAPID_PUBLIC_KEY="your_public_key_here"
VAPID_PRIVATE_KEY="your_private_key_here"

# AWS RDS MySQL の場合
DATABASE_URL="admin:password@tcp(your-rds-endpoint:3306)/kotti_game?parseTime=true"

# オプション（デフォルト値あり）
PORT=8080
SESSION_TTL_MINUTES=60
```

### VAPID鍵の生成

Web Push機能を使用するには、VAPID鍵ペアを生成する必要があります：

```bash
# VAPID鍵ペアを生成
cd server
go run cmd/generate-vapid/main.go
```

出力された鍵を環境変数に設定：

```bash
export VAPID_PUBLIC_KEY="BNm..."
export VAPID_PRIVATE_KEY="abc..."
```

または `.env` ファイルに追加してください。

## データベースセットアップ

### 1. MySQLデータベースを作成

```bash
# ローカルMySQL
mysql -u root -p
CREATE DATABASE kotti_game CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# AWS RDS MySQL の場合
mysql -h your-rds-endpoint -u admin -p
CREATE DATABASE kotti_game CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. マイグレーションツールのインストール

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

### 3. マイグレーションの実行

```bash
export DATABASE_URL="user:password@tcp(localhost:3306)/kotti_game?parseTime=true"
make migrate-up

# AWS RDS MySQL の場合
export DATABASE_URL="admin:password@tcp(your-rds-endpoint:3306)/kotti_game?parseTime=true"
make migrate-up
```

## サーバーの起動

### 開発モード

```bash
make run
```

### 本番ビルド

```bash
make build
./bin/server
```

## API エンドポイント

### セッション管理

- `POST /api/session` - 新規セッション作成
- `GET /api/session/{session_id}` - セッション情報取得
- `POST /api/session/{session_id}/s6/start` - S6開始（7分タイマー）

### S4（診査室パート）回答管理

- `POST /api/session/{session_id}/answers` - 回答保存（逐次保存可能）
- `GET /api/session/{session_id}/answers` - 全回答取得

### S6（存在証明書探索パート）進捗管理

- `POST /api/session/{session_id}/s6/initialize` - 進捗初期化（5箇所分）
- `POST /api/session/{session_id}/s6/verify-location` - 場所到達記録（Web選択）
- `GET /api/session/{session_id}/s6/quiz/{place_id}` - クイズ取得
- `POST /api/session/{session_id}/s6/answer` - クイズ回答送信
- `GET /api/session/{session_id}/s6/progress` - 進捗取得

### メッセージ管理

- `POST /api/session/{session_id}/message` - メッセージ保存（S9）
- `GET /api/messages?place_id={place_id}&limit={limit}` - メッセージ一覧取得

### Web Push通知

- `GET /api/push/vapid-public-key` - VAPID公開鍵取得
- `POST /api/push/subscribe` - プッシュ通知サブスクリプション登録
- `DELETE /api/push/subscriptions/{subscription_id}` - サブスクリプション削除
- `POST /api/push/send/{session_id}` - 即時プッシュ通知送信

### ヘルスチェック

- `GET /api/healthz` - ヘルスチェック

## 実装されている機能

### ✅ 実装済み

1. **セッション管理**
   - UUID形式のセッションID生成
   - 60分の有効期限管理
   - 有効期限切れチェック

2. **S4回答管理**
   - 10問の質問への回答保存（逐次保存対応）
   - リロード後も回答を復元可能

3. **S6進捗管理**
   - 5箇所の場所管理（spiral_stairs, fireplace, back_door_hinge, entrance_door, piano_room）
   - 7分タイマー管理
   - 場所到達判定（Web選択のみ）
   - クイズ生成（プレイヤーの回答から4択生成）
   - クイズ回答検証

4. **メッセージ管理**
   - 最大120文字のメッセージ保存
   - 匿名性保証（session_idのみ保存）
   - 場所別・全体でのメッセージ取得

5. **クイズ生成ロジック**
   - 正解: プレイヤー自身の回答
   - ダミー1: プレイヤーの別回答
   - ダミー2: 過去プレイヤーの匿名回答
   - ダミー3: システム汎用回答

6. **エラーハンドリング**
   - ドメインエラーの適切なHTTPステータスコード変換
   - セッション有効期限切れ（410 Gone）
   - S6時間切れ（408 Request Timeout）
   - 不正な入力（400 Bad Request）

7. **Web Push通知機能** ⭐ NEW
   - VAPID認証によるセキュアなプッシュ通知
   - セッションベースのサブスクリプション管理
   - 即時プッシュ通知送信（ホラー演出用）
   - 自動サブスクリプション無効化（404/410エラー時）
   - プッシュ送信ログ記録
   - MySQL永続化（再起動後も維持）

### ❌ 未実装（後回し）

1. **VOICEVOX統合** - `POST /api/voice/generate`
2. **画像類似度判定** - `POST /api/session/{session_id}/s6/verify-location` の画像アップロード版

※ 画像類似度判定は、Web選択のみで代替可能です。

## 開発コマンド

```bash
# コードフォーマット
make fmt

# ビルド
make build

# テスト実行
make test

# マイグレーション作成
make migrate-create NAME=add_new_feature

# マイグレーションロールバック
make migrate-down
```

## 実装の特徴

### DDD原則

- ビジネスロジックはドメイン層に集約
- Value Objectsによる値の検証
- リポジトリパターンによるデータアクセスの抽象化
- ドメインサービスによる複雑なビジネスロジックの実装

### Clean Architecture

- 依存関係は内側（Domain）に向かう
- インフラストラクチャ層はドメイン層のインターフェースを実装
- ユースケース層でビジネスフローを管理

### セキュリティ

- プレイヤーの匿名性保証（IPアドレス等は保存しない）
- セッションIDによる認証
- CORS対応（開発用）

## トラブルシューティング

### データベース接続エラー

```bash
# MySQLが起動しているか確認
mysqladmin ping

# 接続情報を確認
mysql -h localhost -u user -p -e "SELECT 1;"

# AWS RDS MySQL の場合
mysql -h your-rds-endpoint -u admin -p -e "SELECT 1;"
```

### マイグレーションエラー

```bash
# マイグレーション状態を確認
# MySQL用のURLフォーマットに注意
export DATABASE_URL="user:password@tcp(localhost:3306)/kotti_game"
migrate -path migrations -database "mysql://$DATABASE_URL" version

# 強制的にバージョンを設定（最終手段）
migrate -path migrations -database "mysql://$DATABASE_URL" force 1
```

## 今後の拡張

1. **VOICEVOX統合** - 音声生成機能の実装
2. **画像認識統合** - 既存のgRPCサービスを使った画像類似度判定
3. **セッションクリーンアップバッチ** - 期限切れセッションの自動削除
4. **ログ管理** - 構造化ログの導入（slog, zap）
5. **メトリクス** - Prometheusメトリクスの追加
6. **テストコード** - ユニットテスト、統合テストの追加
