# 赤煉瓦文化館 REST APIサーバー

福岡市赤煉瓦文化館を舞台とした体験型Webホラーゲームのバックエンドサーバーです。

## 目次

- [環境構築](#環境構築)
- [開発環境でのサーバー起動](#開発環境でのサーバー起動)
- [VOICEVOX統合](#voicevox統合)
- [APIエンドポイント](#apiエンドポイント)
- [データベース](#データベース)
- [開発コマンド](#開発コマンド)

---

## 環境構築

### 必要な環境

- Go 1.25.1 以上
- MySQL 8.0
- VOICEVOX Engine（音声合成用）
- golang-migrate（マイグレーションツール）

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd kotti-backend/server
```

### 2. 依存関係のインストール

```bash
go mod download
go mod tidy
```

### 3. MySQLデータベースの作成

```bash
# ローカルMySQL
mysql -u root -p
CREATE DATABASE kotti_game CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;

# AWS RDS MySQL の場合
mysql -h your-rds-endpoint -u admin -p
CREATE DATABASE kotti_game CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

### 4. 環境変数の設定

```bash
# 必須環境変数
export DATABASE_URL="user:password@tcp(localhost:3306)/kotti_game?parseTime=true"

# オプション環境変数
export PORT=8080
export SESSION_TTL_MINUTES=60

# VOICEVOX設定
export VOICEVOX_API_URL="http://127.0.0.1:50021"
export VOICEVOX_SPEAKER_ID=84  # 青山龍星(しっとり)
export AUDIO_OUTPUT_DIR="./audio_files"
export AUDIO_URL_PREFIX="/audio"
export USE_S3_FOR_AUDIO=false  # 本番環境ではtrue
# export S3_BUCKET_NAME="your-bucket-name"  # S3使用時のみ
# export S3_REGION="ap-northeast-1"         # S3使用時のみ
```

### 5. マイグレーションの実行

```bash
# golang-migrateのインストール（macOS）
brew install golang-migrate

# マイグレーション適用
make migrate-up
```

### 6. サーバーの起動

```bash
# 開発モード
make run

# または本番ビルド
make build
./bin/server
```

---

## Docker Composeでの起動（推奨）

Docker Composeを使用することで、MySQL・マイグレーション・APIサーバーを一括で起動できます。

### 手順1: 環境変数ファイルの作成

```bash
# .env.exampleをコピー
cp .env.example .env

# VAPID鍵ペアを生成（プッシュ通知機能に必要）
go run cmd/generate-vapid/main.go

# 生成されたVAPID_PUBLIC_KEYとVAPID_PRIVATE_KEYを.envファイルにコピー
# または、環境変数として直接exportする
```

**重要**: VAPID鍵は必須です。生成しないとサーバーが起動しません。

### 手順2: VOICEVOXサーバーの起動（別ターミナル）

VOICEVOXは別途ローカルで起動してください。

```bash
# VOICEVOXアプリケーションを起動（デフォルトポート: 50021）
# または Docker版を使用
```

#### VOICEVOXサーバーの確認

```bash
# 疎通確認
curl http://127.0.0.1:50021/version

# スピーカー一覧確認
curl http://127.0.0.1:50021/speakers | jq
```

### 手順3: Docker Composeでサービス起動

```bash
# サービスをバックグラウンドで起動
docker compose up -d

# ログを確認
docker compose logs -f app
```

以下のサービスが自動的に起動します：

1. **MySQL 8.0** (ポート 3306)
2. **マイグレーション実行** (自動実行後に終了)
3. **APIサーバー** (ポート 8080)

### 手順4: 動作確認

```bash
# ヘルスチェック
curl http://localhost:8080/api/healthz

# 音声生成テスト
curl -X POST http://localhost:8080/api/voice/generate \
  -H "Content-Type: application/json" \
  -d '{"text": "こんにちは"}' | jq
```

### サービスの停止

```bash
# サービスを停止
docker compose down

# データベースのボリュームも削除する場合
docker compose down -v
```

### トラブルシューティング

#### サービスの状態確認

```bash
# 全サービスの状態を確認
docker compose ps

# ログを確認
docker compose logs app
docker compose logs db
```

#### コンテナに入る

```bash
# APIサーバーコンテナに入る
docker compose exec app sh

# MySQLコンテナに入る
docker compose exec db mysql -u kotti -p kotti_game
```

#### リビルド

```bash
# コードを変更した場合はリビルド
docker compose up -d --build
```

---

## 開発環境でのサーバー起動（ローカル実行）

### 手順1: VOICEVOXサーバーの起動

VOICEVOX Engineをローカルで起動します（別ターミナル）。

```bash
# VOICEVOXアプリケーションを起動するか、Docker版を使用
# デフォルトポート: 50021
```

#### VOICEVOXサーバーの確認

```bash
# 疎通確認
curl http://127.0.0.1:50021/version

# スピーカー一覧確認
curl http://127.0.0.1:50021/speakers | jq
```

### 手順2: MySQLサーバーの起動

```bash
# macOSの場合（Homebrewでインストールしている場合）
brew services start mysql

# または直接起動
mysql.server start
```

### 手順3: 環境変数の設定

```bash
export DATABASE_URL="root:password@tcp(localhost:3306)/kotti_game?parseTime=true"
export VOICEVOX_API_URL="http://127.0.0.1:50021"
```

### 手順4: APIサーバーの起動

```bash
cd server
make run
```

以下のようなログが表示されれば成功です：

```
Server starting on port 8080
VOICEVOX API URL: http://127.0.0.1:50021
Audio output directory: ./audio_files
```

---

## VOICEVOX統合

### 概要

VOICEVOX Engineを使用して、テキストから音声を生成します。S1シーン（1942年パート）で担当者の音声として使用します。

### 使用するスピーカー

- **青山龍星(しっとり)**: Speaker ID = 84

### APIエンドポイント

#### POST /api/voice/generate

テキストから音声を生成します。

**リクエスト**

```bash
curl -X POST http://localhost:8080/api/voice/generate \
  -H "Content-Type: application/json" \
  -d '{
    "text": "こんにちは、赤煉瓦文化館へようこそ。私は担当者の青山と申します。"
  }'
```

**レスポンス**

```json
{
  "audio_url": "/audio/abc123def456_1699999999.wav"
}
```

**オプション: Speaker IDを指定**

```bash
curl -X POST http://localhost:8080/api/voice/generate \
  -H "Content-Type: application/json" \
  -d '{
    "text": "こんにちは",
    "speaker_id": 84
  }'
```

### 音声ファイルの取得

生成された音声ファイルは、レスポンスの`audio_url`から取得できます。

```bash
# 音声ファイルのダウンロード
curl http://localhost:8080/audio/abc123def456_1699999999.wav -o voice.wav

# ブラウザで直接再生
# http://localhost:8080/audio/abc123def456_1699999999.wav
```

### テスト手順

#### 1. VOICEVOXの疎通確認

```bash
# VOICEVOXバージョン確認
curl http://127.0.0.1:50021/version

# 期待されるレスポンス例
# "0.25.0"
```

#### 2. スピーカー一覧の確認

```bash
curl http://127.0.0.1:50021/speakers | jq '.[] | select(.name == "青山龍星")'
```

#### 3. 音声生成テスト

```bash
# テキストから音声を生成
curl -X POST http://localhost:8080/api/voice/generate \
  -H "Content-Type: application/json" \
  -d '{"text": "テストメッセージです"}' | jq

# 期待されるレスポンス
# {
#   "audio_url": "/audio/xxxxxxxxxxxx_1234567890.wav"
# }
```

#### 4. 音声ファイルの取得

```bash
# レスポンスのaudio_urlを使用
curl http://localhost:8080/audio/xxxxxxxxxxxx_1234567890.wav -o test.wav

# ファイルサイズ確認
ls -lh test.wav

# macOSで再生
afplay test.wav
```

### トラブルシューティング

#### エラー: "Failed to generate audio from VOICEVOX"

**原因**: VOICEVOXサーバーが起動していない、または接続できない。

**解決方法**:
```bash
# VOICEVOXサーバーの状態確認
curl http://127.0.0.1:50021/version

# 接続できない場合はVOICEVOXを再起動
```

#### エラー: "failed to create audio output directory"

**原因**: 音声ファイルの出力ディレクトリに書き込み権限がない。

**解決方法**:
```bash
# ディレクトリを作成し、権限を付与
mkdir -p ./audio_files
chmod 755 ./audio_files
```

#### エラー: "Text is required"

**原因**: リクエストボディに`text`フィールドが空、またはnull。

**解決方法**:
```bash
# 正しいリクエスト形式
curl -X POST http://localhost:8080/api/voice/generate \
  -H "Content-Type: application/json" \
  -d '{"text": "音声化するテキスト"}'
```

---

## APIエンドポイント

### ヘルスチェック

```bash
GET /api/healthz
```

### セッション管理

```bash
POST   /api/session                      # 新規セッション作成
GET    /api/session/{session_id}         # セッション情報取得
POST   /api/session/{session_id}/s6/start  # S6開始
```

### S4回答管理

```bash
POST   /api/session/{session_id}/answers   # 回答保存
GET    /api/session/{session_id}/answers   # 回答一覧取得
```

### S6進捗・クイズ管理

```bash
POST   /api/session/{session_id}/s6/initialize       # S6進捗初期化
POST   /api/session/{session_id}/s6/verify-location  # 場所到達記録
GET    /api/session/{session_id}/s6/quiz/{place_id}  # クイズ取得
POST   /api/session/{session_id}/s6/answer           # クイズ回答
GET    /api/session/{session_id}/s6/progress         # 進捗状況取得
```

### メッセージ管理

```bash
POST   /api/session/{session_id}/message   # メッセージ保存
GET    /api/messages                        # メッセージ一覧取得
```

### 音声生成（VOICEVOX）

```bash
POST   /api/voice/generate   # 音声生成
```

### プッシュ通知

```bash
GET    /api/push/vapid-public-key                     # VAPID公開鍵取得
POST   /api/push/subscribe                            # サブスクリプション登録
DELETE /api/push/subscriptions/{subscription_id}      # サブスクリプション削除
POST   /api/push/send/{session_id}                    # プッシュ通知送信
```

---

## データベース

### マイグレーション管理

```bash
# マイグレーション適用
make migrate-up

# マイグレーションロールバック
make migrate-down

# 新規マイグレーション作成
make migrate-create NAME=add_new_feature
```

### テーブル一覧

| テーブル名 | 説明 |
|-----------|------|
| sessions | ゲームセッション |
| session_answers | S4の内省質問への回答 |
| session_s6_progress | S6探索進捗 |
| quiz_questions | 生成されたクイズ |
| player_messages | プレイヤーが刻んだメッセージ |
| location_images | 場所の基準画像（未使用） |
| push_subscriptions | プッシュ通知サブスクリプション |
| push_logs | プッシュ通知送信ログ |

---

## 開発コマンド

### ビルド・実行

```bash
make run          # サーバー起動（開発モード）
make run-cmd      # cmd/server/main.go から起動
make build        # バイナリビルド（bin/server）
```

### コード品質

```bash
make fmt          # コードフォーマット（必須）
make lint         # Linterチェック（必須）
```

### データベース

```bash
make migrate-up           # マイグレーション適用
make migrate-down         # マイグレーションロールバック
make migrate-create NAME=migration_name  # 新規マイグレーション作成
```

### テスト

```bash
make test         # テスト実行（PoCのため新規テスト作成は不要）
```

---

## アーキテクチャ

本プロジェクトは **DDD (Domain-Driven Design)** と **Clean Architecture** の原則に基づいています。

### レイヤー構成

```
┌─────────────────────────────────────────┐
│   Interface Layer (外部インターフェース層)   │
│   - HTTP Handlers                        │
│   - DTOs (Data Transfer Objects)        │
└─────────────────────────────────────────┘
              ↓ 依存
┌─────────────────────────────────────────┐
│   Application Layer (アプリケーション層)    │
│   - Use Cases                            │
└─────────────────────────────────────────┘
              ↓ 依存
┌─────────────────────────────────────────┐
│   Domain Layer (ドメイン層)               │
│   - Entities (Models)                    │
│   - Value Objects                        │
│   - Repository Interfaces                │
│   - Domain Services                      │
└─────────────────────────────────────────┘
              ↑ 実装
┌─────────────────────────────────────────┐
│   Infrastructure Layer (インフラ層)       │
│   - Repository Implementations (MySQL)   │
│   - VOICEVOX Client                      │
│   - Database Connection                  │
└─────────────────────────────────────────┘
```

### ディレクトリ構造

```
server/
├── cmd/server/              # アプリケーションエントリーポイント
├── internal/
│   ├── interfaces/          # Interface Layer
│   │   └── http/
│   │       ├── handler/     # HTTPハンドラー
│   │       └── dto/         # データ転送オブジェクト
│   ├── application/         # Application Layer
│   │   └── usecase/         # ユースケース
│   ├── domain/              # Domain Layer
│   │   ├── model/           # エンティティ
│   │   ├── valueobject/     # 値オブジェクト
│   │   ├── repository/      # リポジトリインターフェース
│   │   └── service/         # ドメインサービス
│   └── infrastructure/      # Infrastructure Layer
│       ├── persistence/     # リポジトリ実装（MySQL）
│       ├── database/        # DB接続管理
│       └── voicevox/        # VOICEVOXクライアント
├── pkg/
│   ├── config/              # 環境設定
│   └── errors/              # カスタムエラー
├── migrations/              # DBマイグレーションファイル
└── audio_files/             # 生成された音声ファイル（自動作成）
```

---

## 本番環境デプロイ

### AWS ECS Fargate

本番環境はAWS ECS Fargateにデプロイされます。

#### 必要なリソース

- ECS Cluster
- Application Load Balancer (ALB)
- RDS MySQL 8.0
- ECR (Container Registry)
- S3 Bucket (音声ファイル保存用)
- EC2 (VOICEVOX Engine専用)

#### 環境変数（本番）

```bash
DATABASE_URL="admin:password@tcp(rds-endpoint:3306)/kotti_game?parseTime=true"
PORT=8080
SESSION_TTL_MINUTES=60
VOICEVOX_API_URL="http://voicevox-ec2-private-ip:50021"
VOICEVOX_SPEAKER_ID=84
USE_S3_FOR_AUDIO=true
S3_BUCKET_NAME="kotti-game-audio-files"
S3_REGION="ap-northeast-1"
```

---

## ライセンス

This project is part of a hackathon project for 赤煉瓦文化館.

---

## 参考資料

- [VOICEVOX Engine API仕様](./voicevox.openapi.json)
- [実装状況ドキュメント](./IMPLEMENTATION_STATUS.md)
- [セットアップガイド](./SETUP.md)
- [プロジェクト全体仕様](../CLAUDE.md)
