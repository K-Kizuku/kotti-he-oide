# サーバー実装状況ドキュメント

**最終更新**: 2025-11-07
**対象**: 赤煉瓦文化館ホラーゲーム REST APIバックエンド

このドキュメントは、後任者が現在の実装状況を正確に把握し、スムーズに開発を継続できるようにまとめたものです。

---

## 目次

1. [プロジェクト概要](#プロジェクト概要)
2. [実装状況サマリー](#実装状況サマリー)
3. [アーキテクチャ詳細](#アーキテクチャ詳細)
4. [ディレクトリ構造](#ディレクトリ構造)
5. [実装済み機能](#実装済み機能)
6. [未実装機能](#未実装機能)
7. [データベース詳細](#データベース詳細)
8. [APIエンドポイント一覧](#apiエンドポイント一覧)
9. [開発環境セットアップ](#開発環境セットアップ)
10. [次のステップ](#次のステップ)
11. [既知の問題・制限事項](#既知の問題制限事項)
12. [開発のヒント](#開発のヒント)

---

## プロジェクト概要

### ゲームコンセプト
福岡市赤煉瓦文化館を舞台とした体験型Webホラーゲームのバックエンドサーバー。プレイヤーは複数の時代（1942年、1972年、2002年、2025年）を行き来し、内省的な質問への回答、7分間の探索、クイズ、メッセージ刻みを通じて、建物の歴史と自己の生きる意義を体験します。

### 技術スタック
- **言語**: Go 1.25.1
- **データベース**: AWS RDS MySQL 8.0
- **DBドライバー**: go-sql-driver/mysql v1.8.1
- **マイグレーション**: golang-migrate
- **アーキテクチャ**: DDD (Domain-Driven Design) + Clean Architecture
- **HTTP**: 標準ライブラリ (net/http)
- **gRPC**: google.golang.org/grpc v1.65.0 (画像認識・VOICEVOX統合用)

### ゲームフロー
```
S0 (注意書き・許可)
  → S1 (1942年パート: 担当者との会話)
  → S2 (お気に入りの場所)
  → S3 (移動指示+ホラー演出)
  → S4 (診査室: 内省10問)
  → S5 (死亡届受理)
  → S6 (存在証明書探索: 7分タイマー)
  → S7 (2002年パート: メッセージ受け取り)
  → S8 (メインホール探索: 3分)
  → S9 (2025年: メッセージ刻み)
```

---

## 実装状況サマリー

### ✅ 実装完了（本番使用可能）

| カテゴリ | 機能 | 実装率 |
|---------|------|--------|
| セッション管理 | UUID発行、有効期限管理、S6タイマー管理 | 100% |
| S4回答管理 | 10問の質問への回答保存・取得 | 100% |
| S6進捗管理 | 5箇所の場所管理、到達判定、クイズ生成・検証 | 100% |
| メッセージ管理 | 最大120文字メッセージ保存・取得（匿名） | 100% |
| ドメイン層 | Value Objects, Entities, Domain Services | 100% |
| アプリケーション層 | Use Cases | 100% |
| インフラ層 | MySQL Repository実装 | 100% |
| インターフェース層 | HTTP Handlers, DTOs | 90% |
| マイグレーション | MySQL 8.0スキーマ定義 | 100% |

### ⚠️ 部分実装（統合待ち）

| 機能 | 状態 | 備考 |
|------|------|------|
| 画像認識統合 | ハンドラー実装済み、ルーティング未設定 | `ml_handler.go`に実装あり |

### ❌ 未実装（仕様定義済み）

| 機能 | 優先度 | 備考 |
|------|--------|------|
| VOICEVOX統合 | 中 | S1で音声生成が必要 |
| S6画像類似度判定の実際の使用 | 低 | Web選択で代替可能 |

---

## アーキテクチャ詳細

### レイヤー構成

本プロジェクトはClean Architecture + DDDに基づき、以下の4層構造を採用しています：

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
│   - Database Connection                  │
└─────────────────────────────────────────┘
```

### 依存関係のルール
- **依存の方向**: 外側から内側（Domain層）へ
- **Infrastructure層**: Domain層のインターフェースを実装（依存性逆転の原則）
- **Domain層**: 他の層に依存しない（純粋なビジネスロジック）

---

## ディレクトリ構造

```
server/
├── cmd/
│   └── server/
│       └── main.go                          # エントリーポイント、ルーティング設定
│
├── internal/
│   ├── interfaces/                          # Interface Layer
│   │   └── http/
│   │       ├── handler/
│   │       │   ├── health_handler.go        # ✅ ヘルスチェック
│   │       │   ├── session_handler.go       # ✅ セッション管理
│   │       │   ├── answer_handler.go        # ✅ S4回答管理
│   │       │   ├── s6_handler.go            # ✅ S6進捗・クイズ管理
│   │       │   ├── message_handler.go       # ✅ メッセージ管理
│   │       │   └── ml_handler.go            # ⚠️ 画像認識プロキシ（ルーティング未設定）
│   │       └── dto/
│   │           ├── session_dto.go           # ✅ セッション用DTO
│   │           ├── answer_dto.go            # ✅ 回答用DTO
│   │           ├── quiz_dto.go              # ✅ クイズ用DTO
│   │           └── message_dto.go           # ✅ メッセージ用DTO
│   │
│   ├── application/                         # Application Layer
│   │   └── usecase/
│   │       ├── session_usecase.go           # ✅ セッション管理ユースケース
│   │       ├── s4_answer_usecase.go         # ✅ S4回答ユースケース
│   │       ├── s6_usecase.go                # ✅ S6進捗・クイズユースケース
│   │       └── message_usecase.go           # ✅ メッセージユースケース
│   │
│   ├── domain/                              # Domain Layer
│   │   ├── model/
│   │   │   ├── session.go                   # ✅ セッションエンティティ
│   │   │   ├── session_answer.go            # ✅ S4回答エンティティ
│   │   │   ├── s6_progress.go               # ✅ S6進捗エンティティ
│   │   │   ├── quiz_question.go             # ✅ クイズエンティティ
│   │   │   └── player_message.go            # ✅ プレイヤーメッセージエンティティ
│   │   │
│   │   ├── valueobject/
│   │   │   ├── session_id.go                # ✅ セッションID（UUID）
│   │   │   ├── place_id.go                  # ✅ 場所ID（5箇所の定義）
│   │   │   ├── question_id.go               # ✅ 質問ID（Q1〜Q10）
│   │   │   └── quiz_id.go                   # ✅ クイズID（UUID）
│   │   │
│   │   ├── repository/
│   │   │   ├── session_repository.go        # ✅ セッションリポジトリインターフェース
│   │   │   ├── session_answer_repository.go # ✅ S4回答リポジトリインターフェース
│   │   │   ├── s6_progress_repository.go    # ✅ S6進捗リポジトリインターフェース
│   │   │   ├── quiz_question_repository.go  # ✅ クイズリポジトリインターフェース
│   │   │   └── player_message_repository.go # ✅ メッセージリポジトリインターフェース
│   │   │
│   │   └── service/
│   │       ├── session_service.go           # ✅ セッション検証サービス
│   │       ├── quiz_service.go              # ✅ クイズ生成サービス（4択）
│   │       └── s6_service.go                # ✅ S6タイマー検証サービス
│   │
│   ├── infrastructure/                      # Infrastructure Layer
│   │   ├── database/
│   │   │   └── mysql.go                     # ✅ MySQL接続管理
│   │   └── persistence/
│   │       ├── mysql_session_repository.go          # ✅ セッションリポジトリ実装
│   │       ├── mysql_session_answer_repository.go   # ✅ S4回答リポジトリ実装
│   │       ├── mysql_s6_progress_repository.go      # ✅ S6進捗リポジトリ実装
│   │       ├── mysql_quiz_question_repository.go    # ✅ クイズリポジトリ実装
│   │       └── mysql_player_message_repository.go   # ✅ メッセージリポジトリ実装
│   │
│   └── gen/                                 # 生成コード
│       └── image_recognition/v1/
│           ├── image_recognition.pb.go      # ✅ gRPC Protobuf定義
│           └── image_recognition_grpc.pb.go # ✅ gRPC クライアント
│
├── pkg/
│   ├── config/
│   │   └── config.go                        # ✅ 環境変数設定
│   └── errors/
│       └── domain_error.go                  # ✅ カスタムエラー定義
│
├── migrations/
│   ├── 000001_create_game_tables.up.sql     # ✅ MySQL 8.0スキーマ定義
│   └── 000001_create_game_tables.down.sql   # ✅ ロールバック定義
│
├── go.mod                                   # ✅ 依存関係管理（MySQL対応済み）
├── go.sum                                   # ✅ 依存関係ロック
├── Makefile                                 # ✅ 開発コマンド（MySQL対応済み）
├── SETUP.md                                 # ✅ セットアップガイド
└── IMPLEMENTATION_STATUS.md                 # 📄 本ドキュメント
```

---

## 実装済み機能

### 1. セッション管理 (Session Management)

#### ファイル
- Domain: `internal/domain/model/session.go`
- Value Object: `internal/domain/valueobject/session_id.go`
- Repository Interface: `internal/domain/repository/session_repository.go`
- Repository Impl: `internal/infrastructure/persistence/mysql_session_repository.go`
- Domain Service: `internal/domain/service/session_service.go`
- Use Case: `internal/application/usecase/session_usecase.go`
- Handler: `internal/interfaces/http/handler/session_handler.go`
- DTO: `internal/interfaces/http/dto/session_dto.go`

#### 主要機能
- **セッション作成**: UUID v4によるセッションID発行
- **有効期限管理**: 作成時刻 + 60分（環境変数で設定可能）
- **S6タイマー管理**: 7分カウントダウンの開始時刻記録
- **セッション検証**: 有効期限切れチェック（ドメインサービス）

#### ビジネスルール
- セッションIDは必ずUUID v4形式
- 有効期限切れのセッションへのアクセスは `SESSION_EXPIRED` エラー
- S6開始時刻から7分経過後のクイズ回答は `S6_TIME_EXPIRED` エラー（最後のピース回答中は例外）

#### データベーステーブル
```sql
CREATE TABLE sessions (
    session_id CHAR(36) PRIMARY KEY,
    current_scene VARCHAR(10) NOT NULL DEFAULT 'S0',
    s6_started_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    INDEX idx_sessions_expires_at (expires_at),
    INDEX idx_sessions_created_at (created_at)
);
```

---

### 2. S4回答管理 (Answer Management)

#### ファイル
- Domain: `internal/domain/model/session_answer.go`
- Value Object: `internal/domain/valueobject/question_id.go`
- Repository Interface: `internal/domain/repository/session_answer_repository.go`
- Repository Impl: `internal/infrastructure/persistence/mysql_session_answer_repository.go`
- Use Case: `internal/application/usecase/s4_answer_usecase.go`
- Handler: `internal/interfaces/http/handler/answer_handler.go`
- DTO: `internal/interfaces/http/dto/answer_dto.go`

#### 主要機能
- **逐次保存**: 質問ごとに回答を即座に保存（Wi-Fi切断対策）
- **重複更新**: 同じ質問への再回答時はUPDATE（ON DUPLICATE KEY UPDATE）
- **全回答取得**: セッション全体の回答一覧を取得
- **回答検証**: "なし"や空文字はバリデーションで弾く

#### 質問ID定義
- `Q1`〜`Q10`: 固定10問（小学生時代の夢中、尊敬した人、人生の願望など）
- `Q8`: "人生の最期に達成したいこと"（S7で再入力が必要）

#### データベーステーブル
```sql
CREATE TABLE session_answers (
    session_id CHAR(36) NOT NULL,
    question_id VARCHAR(10) NOT NULL,
    answer_text TEXT NOT NULL,
    answered_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (session_id, question_id),
    FOREIGN KEY (session_id) REFERENCES sessions(session_id) ON DELETE CASCADE
);
```

---

### 3. S6進捗・クイズ管理 (S6 Progress & Quiz)

#### ファイル
- Domain Models:
  - `internal/domain/model/s6_progress.go`
  - `internal/domain/model/quiz_question.go`
- Value Objects:
  - `internal/domain/valueobject/place_id.go`
  - `internal/domain/valueobject/quiz_id.go`
- Repository Interfaces:
  - `internal/domain/repository/s6_progress_repository.go`
  - `internal/domain/repository/quiz_question_repository.go`
- Repository Impls:
  - `internal/infrastructure/persistence/mysql_s6_progress_repository.go`
  - `internal/infrastructure/persistence/mysql_quiz_question_repository.go`
- Domain Services:
  - `internal/domain/service/quiz_service.go`
  - `internal/domain/service/s6_service.go`
- Use Case: `internal/application/usecase/s6_usecase.go`
- Handler: `internal/interfaces/http/handler/s6_handler.go`
- DTO: `internal/interfaces/http/dto/quiz_dto.go`

#### 主要機能

##### S6進捗管理
- **5箇所の場所管理**:
  - `spiral_stairs`: 螺旋階段を見上げる高い天井
  - `fireplace`: メインホールの暖炉のレンガ
  - `back_door_hinge`: 裏玄関の扉の蝶番
  - `entrance_door`: 入口エントランスの扉
  - `piano_room`: 階上応接室のピアノ
- **到達判定**:
  - `verified_by: "photo"`: 画像類似度判定（未統合）
  - `verified_by: "manual"`: Web選択（現在の代替手段）
- **クイズ生成**: 各場所ごとに1問ずつ、計5問を事前生成
- **正誤判定**: 正解時のみ `correct: true` を記録

##### クイズ生成ロジック（重要）
クイズは **4択形式** で、以下のルールで選択肢を生成します：

1. **正解（AnswerIndex 0〜3のいずれか）**: プレイヤー自身のS4回答
2. **ダミー1**: プレイヤーの別のS4回答（ランダム選択）
3. **ダミー2**: 過去プレイヤーの匿名回答（ランダム選択）
4. **ダミー3**: システム汎用回答（例: "特になし"、"覚えていない"）

選択肢の順序はランダムにシャッフルされ、`AnswerIndex`に正解の位置（0-3）を記録します。

#### ビジネスルール
- S6開始から7分経過後のクイズ回答は拒否（最後のピース回答中を除く）
- 不正解時は即座に再挑戦可能（ホラー演出はフロントエンド側）
- 5箇所すべてで正解しないと先に進めない

#### データベーステーブル

```sql
CREATE TABLE session_s6_progress (
    session_id CHAR(36) NOT NULL,
    place_id VARCHAR(50) NOT NULL,
    verified_by VARCHAR(20),
    quiz_id CHAR(36),
    answered BOOLEAN NOT NULL DEFAULT FALSE,
    correct BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (session_id, place_id),
    FOREIGN KEY (session_id) REFERENCES sessions(session_id) ON DELETE CASCADE
);

CREATE TABLE quiz_questions (
    quiz_id CHAR(36) PRIMARY KEY,
    session_id CHAR(36) NOT NULL,
    place_id VARCHAR(50) NOT NULL,
    question_text TEXT NOT NULL,
    option_1 TEXT NOT NULL,
    option_2 TEXT NOT NULL,
    option_3 TEXT NOT NULL,
    option_4 TEXT NOT NULL,
    answer_index INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(session_id) ON DELETE CASCADE
);
```

---

### 4. メッセージ管理 (Player Message)

#### ファイル
- Domain: `internal/domain/model/player_message.go`
- Repository Interface: `internal/domain/repository/player_message_repository.go`
- Repository Impl: `internal/infrastructure/persistence/mysql_player_message_repository.go`
- Use Case: `internal/application/usecase/message_usecase.go`
- Handler: `internal/interfaces/http/handler/message_handler.go`
- DTO: `internal/interfaces/http/dto/message_dto.go`

#### 主要機能
- **メッセージ保存**: 最大120文字のメッセージを匿名で保存
- **場所指定**: S2で選んだ5箇所のいずれかに紐付け
- **匿名性保証**: IPアドレス・端末IDは保存せず、session_idのみ
- **一覧取得**:
  - 場所別メッセージ取得（limit付き）
  - 全体メッセージ取得（limit付き）

#### ビジネスルール
- メッセージは最大120文字
- 匿名性を保証するため、session_id以外の個人情報は保存しない
- 過去プレイヤーのメッセージは完全匿名で提供（S8で表示）

#### データベーステーブル
```sql
CREATE TABLE player_messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id CHAR(36) NOT NULL,
    place_id VARCHAR(50) NOT NULL,
    message_text VARCHAR(120) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_player_messages_place_id (place_id),
    INDEX idx_player_messages_created_at (created_at)
);
```

---

### 5. エラーハンドリング (Domain Errors)

#### ファイル
- `pkg/errors/domain_error.go`

#### 定義済みエラーコード
```go
const (
    SESSION_NOT_FOUND       = "SESSION_NOT_FOUND"
    SESSION_EXPIRED         = "SESSION_EXPIRED"
    INVALID_SESSION_ID      = "INVALID_SESSION_ID"
    INVALID_PLACE_ID        = "INVALID_PLACE_ID"
    INVALID_QUESTION_ID     = "INVALID_QUESTION_ID"
    QUIZ_NOT_FOUND          = "QUIZ_NOT_FOUND"
    S6_NOT_STARTED          = "S6_NOT_STARTED"
    S6_TIME_EXPIRED         = "S6_TIME_EXPIRED"
    ANSWER_NOT_FOUND        = "ANSWER_NOT_FOUND"
    INSUFFICIENT_ANSWERS    = "INSUFFICIENT_ANSWERS"
    INVALID_QUIZ_ANSWER     = "INVALID_QUIZ_ANSWER"
)
```

#### HTTPステータスコード変換
- `SESSION_EXPIRED` → 410 Gone
- `S6_TIME_EXPIRED` → 408 Request Timeout
- `*_NOT_FOUND` → 404 Not Found
- `INVALID_*` → 400 Bad Request

---

## 未実装機能

### 1. VOICEVOX統合 ❌

#### 必要な実装
- **ハンドラー**: `internal/interfaces/http/handler/voice_handler.go` を新規作成
- **エンドポイント**: `POST /api/voice/generate`
- **リクエスト**: `{ "text": "音声化するテキスト", "speaker_id": 青山龍星のID }`
- **レスポンス**: `{ "audio_url": "S3上の音声ファイルURL" }`

#### 統合手順
1. EC2上でVOICEVOXサーバーを起動（環境変数 `VOICEVOX_API_URL`）
2. ハンドラーでHTTPリクエストをVOICEVOXに送信
3. 生成された音声ファイルをS3にアップロード
4. S3のURLをフロントエンドに返却

#### 使用シーン
- **S1**: 担当者との会話音声（青山龍星(しっとり)）

---

### 2. 画像類似度判定の実際の統合 ⚠️

#### 現状
- **ハンドラー実装**: `internal/interfaces/http/handler/ml_handler.go` に実装済み
- **gRPCクライアント**: `internal/gen/image_recognition/v1/` に定義済み
- **問題**: `cmd/server/main.go` でルーティングが未設定

#### 統合手順
1. `main.go` にMLHandlerを追加：
   ```go
   mlHandler := handler.NewMLHandler()
   mux.HandleFunc("POST /api/session/{session_id}/s6/verify-location", mlHandler.RecognizeImageProxy)
   ```
2. 環境変数 `IMAGE_RECOGNITION_GRPC_ADDR` を設定（デフォルト: `127.0.0.1:50051`）
3. 画像認識gRPCサーバーを起動
4. フロントエンドから画像をmultipart/form-dataで送信

#### 代替手段
現在は **Web選択（"この場所にいることにする"ボタン）** で代替しているため、優先度は低い。

---

### 3. セッションクリーンアップバッチ ❌

#### 必要な実装
期限切れセッションを定期的に削除するバッチ処理。

#### 実装案
- **方法1**: Cronジョブで `DELETE FROM sessions WHERE expires_at < NOW()` を実行
- **方法2**: Goのバックグラウンドゴルーチンで定期実行
- **方法3**: RDSのイベントスケジューラ（MySQL 8.0対応）

---

## データベース詳細

### 接続設定

#### 環境変数
```bash
# 必須
DATABASE_URL="user:password@tcp(localhost:3306)/kotti_game?parseTime=true"

# AWS RDS MySQL の場合
DATABASE_URL="admin:password@tcp(your-rds-endpoint:3306)/kotti_game?parseTime=true"

# オプション
PORT=8080
SESSION_TTL_MINUTES=60
```

#### 接続プール設定
- **MaxOpenConns**: 25
- **MaxIdleConns**: 5
- **ConnMaxLifetime**: 5分

### マイグレーション

#### 実行方法
```bash
# 適用
export DATABASE_URL="user:password@tcp(localhost:3306)/kotti_game?parseTime=true"
make migrate-up

# ロールバック
make migrate-down

# 新規マイグレーション作成
make migrate-create NAME=add_new_feature
```

#### マイグレーションファイル
- `migrations/000001_create_game_tables.up.sql`: 全テーブル作成
- `migrations/000001_create_game_tables.down.sql`: 全テーブル削除

### テーブル一覧

| テーブル名 | 説明 | 主キー | 外部キー |
|-----------|------|-------|---------|
| sessions | ゲームセッション | session_id | - |
| session_answers | S4回答 | (session_id, question_id) | session_id |
| session_s6_progress | S6進捗 | (session_id, place_id) | session_id |
| quiz_questions | クイズ | quiz_id | session_id |
| player_messages | プレイヤーメッセージ | id (AUTO_INCREMENT) | - |
| location_images | 場所基準画像（未使用） | id (AUTO_INCREMENT) | - |

### インデックス戦略

#### sessions
- `idx_sessions_expires_at`: 期限切れセッション検索用
- `idx_sessions_created_at`: セッション作成時刻順ソート用

#### player_messages
- `idx_player_messages_place_id`: 場所別メッセージ取得用
- `idx_player_messages_created_at`: 最新順ソート用

---

## APIエンドポイント一覧

### 実装状況凡例
- ✅: 実装済み、動作確認済み
- ⚠️: 部分実装、統合待ち
- ❌: 未実装

### セッション管理

| メソッド | エンドポイント | 実装状況 | 説明 |
|---------|---------------|---------|------|
| POST | `/api/session` | ✅ | 新規セッション作成 |
| GET | `/api/session/{session_id}` | ✅ | セッション情報取得 |
| POST | `/api/session/{session_id}/s6/start` | ✅ | S6開始（7分タイマー開始） |

#### POST /api/session

**リクエスト**: なし

**レスポンス**:
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "current_scene": "S0",
  "created_at": "2025-11-07T10:00:00Z",
  "expires_at": "2025-11-07T11:00:00Z"
}
```

#### GET /api/session/{session_id}

**レスポンス**:
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "current_scene": "S4",
  "s6_started_at": "2025-11-07T10:30:00Z",
  "created_at": "2025-11-07T10:00:00Z",
  "expires_at": "2025-11-07T11:00:00Z"
}
```

**エラー**:
- 404: セッション未存在
- 410: セッション期限切れ

---

### S4回答管理

| メソッド | エンドポイント | 実装状況 | 説明 |
|---------|---------------|---------|------|
| POST | `/api/session/{session_id}/answers` | ✅ | S4回答保存（逐次） |
| GET | `/api/session/{session_id}/answers` | ✅ | 保存された回答一覧取得 |

#### POST /api/session/{session_id}/answers

**リクエスト**:
```json
{
  "question_id": "Q1",
  "answer_text": "野球に夢中でした"
}
```

**レスポンス**:
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "question_id": "Q1",
  "answer_text": "野球に夢中でした",
  "answered_at": "2025-11-07T10:15:00Z"
}
```

#### GET /api/session/{session_id}/answers

**レスポンス**:
```json
{
  "answers": [
    {
      "question_id": "Q1",
      "answer_text": "野球に夢中でした",
      "answered_at": "2025-11-07T10:15:00Z"
    },
    {
      "question_id": "Q2",
      "answer_text": "母親です",
      "answered_at": "2025-11-07T10:16:00Z"
    }
  ]
}
```

---

### S6進捗・クイズ管理

| メソッド | エンドポイント | 実装状況 | 説明 |
|---------|---------------|---------|------|
| POST | `/api/session/{session_id}/s6/initialize` | ✅ | S6進捗初期化（5箇所） |
| POST | `/api/session/{session_id}/s6/verify-location` | ✅ | 場所到達記録（Web選択） |
| GET | `/api/session/{session_id}/s6/quiz/{place_id}` | ✅ | クイズ取得 |
| POST | `/api/session/{session_id}/s6/answer` | ✅ | クイズ回答送信 |
| GET | `/api/session/{session_id}/s6/progress` | ✅ | 進捗状況取得 |

#### POST /api/session/{session_id}/s6/initialize

**リクエスト**: なし

**レスポンス**:
```json
{
  "message": "S6 progress initialized for 5 places"
}
```

#### POST /api/session/{session_id}/s6/verify-location

**リクエスト**:
```json
{
  "place_id": "spiral_stairs"
}
```

**レスポンス**:
```json
{
  "verified": true,
  "verified_by": "manual"
}
```

#### GET /api/session/{session_id}/s6/quiz/{place_id}

**レスポンス**:
```json
{
  "quiz_id": "660e8400-e29b-41d4-a716-446655440001",
  "place_id": "spiral_stairs",
  "question_text": "小学生の頃、あなたが一番夢中だったことは？",
  "options": [
    "野球に夢中でした",
    "特になし",
    "読書が好きでした",
    "覚えていない"
  ]
}
```

**注意**: `answer_index`（正解の位置）はレスポンスに含まれません。

#### POST /api/session/{session_id}/s6/answer

**リクエスト**:
```json
{
  "quiz_id": "660e8400-e29b-41d4-a716-446655440001",
  "answer_index": 0
}
```

**レスポンス**:
```json
{
  "correct": true,
  "message": "Correct answer!"
}
```

**エラー**:
- 408: S6時間切れ（7分経過）
- 400: 不正な回答インデックス

#### GET /api/session/{session_id}/s6/progress

**レスポンス**:
```json
{
  "progress": [
    {
      "place_id": "spiral_stairs",
      "verified_by": "manual",
      "answered": true,
      "correct": true
    },
    {
      "place_id": "fireplace",
      "verified_by": null,
      "answered": false,
      "correct": false
    }
  ]
}
```

---

### メッセージ管理

| メソッド | エンドポイント | 実装状況 | 説明 |
|---------|---------------|---------|------|
| POST | `/api/session/{session_id}/message` | ✅ | メッセージ保存（S9） |
| GET | `/api/messages?place_id={place_id}&limit={limit}` | ✅ | メッセージ一覧取得 |

#### POST /api/session/{session_id}/message

**リクエスト**:
```json
{
  "place_id": "fireplace",
  "message_text": "この建物を大切に受け継いでください"
}
```

**レスポンス**:
```json
{
  "message": "Message saved successfully"
}
```

#### GET /api/messages

**クエリパラメータ**:
- `place_id` (オプション): 場所指定（未指定時は全体）
- `limit` (オプション): 最大取得件数（デフォルト: 50）

**レスポンス**:
```json
{
  "messages": [
    {
      "place_id": "fireplace",
      "message_text": "この建物を大切に受け継いでください",
      "created_at": "2025-11-07T10:50:00Z"
    },
    {
      "place_id": "fireplace",
      "message_text": "歴史を忘れないで",
      "created_at": "2025-11-07T10:45:00Z"
    }
  ]
}
```

---

### ヘルスチェック

| メソッド | エンドポイント | 実装状況 | 説明 |
|---------|---------------|---------|------|
| GET | `/api/healthz` | ✅ | ヘルスチェック |

**レスポンス**:
```json
{
  "status": "healthy"
}
```

---

### 画像認識（未統合）

| メソッド | エンドポイント | 実装状況 | 説明 |
|---------|---------------|---------|------|
| POST | `/api/ml/recognize` | ⚠️ | 画像類似度判定（ルーティング未設定） |
| GET | `/api/ml/hello` | ⚠️ | gRPC疎通確認（ルーティング未設定） |

---

### VOICEVOX（未実装）

| メソッド | エンドポイント | 実装状況 | 説明 |
|---------|---------------|---------|------|
| POST | `/api/voice/generate` | ❌ | 音声生成 |

---

## 開発環境セットアップ

### 前提条件
- Go 1.25.1 以上
- MySQL 8.0（ローカルまたはAWS RDS）
- golang-migrate（マイグレーションツール）

### セットアップ手順

#### 1. リポジトリクローン
```bash
git clone <repository-url>
cd kotti-backend/server
```

#### 2. 依存関係インストール
```bash
go mod download
go mod tidy
```

#### 3. MySQLデータベース作成
```bash
# ローカルMySQL
mysql -u root -p
CREATE DATABASE kotti_game CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# AWS RDS MySQL
mysql -h your-rds-endpoint -u admin -p
CREATE DATABASE kotti_game CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 4. 環境変数設定
```bash
# .envファイル作成（任意）
export DATABASE_URL="user:password@tcp(localhost:3306)/kotti_game?parseTime=true"
export PORT=8080
export SESSION_TTL_MINUTES=60
```

#### 5. マイグレーション実行
```bash
# golang-migrateインストール（macOS）
brew install golang-migrate

# マイグレーション適用
make migrate-up
```

#### 6. サーバー起動
```bash
# 開発モード
make run

# 本番ビルド
make build
./bin/server
```

#### 7. 動作確認
```bash
# ヘルスチェック
curl http://localhost:8080/api/healthz

# セッション作成
curl -X POST http://localhost:8080/api/session
```

---

## 次のステップ

### 優先度：高

1. **VOICEVOX統合** (S1で音声生成が必要)
   - [ ] `internal/interfaces/http/handler/voice_handler.go` を作成
   - [ ] EC2上でVOICEVOXサーバー起動
   - [ ] S3への音声ファイルアップロード実装
   - [ ] `cmd/server/main.go` にルーティング追加

2. **画像認識統合** (S6で場所到達判定を強化)
   - [ ] `cmd/server/main.go` にMLHandlerのルーティング追加
   - [ ] 画像認識gRPCサーバー起動確認
   - [ ] フロントエンドとの統合テスト

### 優先度：中

3. **セッションクリーンアップバッチ**
   - [ ] 期限切れセッション削除処理の実装
   - [ ] Cronジョブまたはバックグラウンドゴルーチン設定

4. **ログ管理**
   - [ ] 構造化ログ導入（slog, zapなど）
   - [ ] エラーログ、アクセスログの分離
   - [ ] ログレベル設定（環境変数）

### 優先度：低

5. **テストコード**
   - [ ] ユニットテスト（Domain層、Application層）
   - [ ] 統合テスト（Repository層）
   - [ ] E2Eテスト（Handler層）

6. **メトリクス**
   - [ ] Prometheusメトリクス追加
   - [ ] レスポンスタイム計測
   - [ ] エラー率監視

---

## 既知の問題・制限事項

### 1. セッション有効期限の自動延長なし
- **問題**: セッション作成から60分後に強制終了
- **影響**: プレイ中にセッション切れの可能性
- **対策案**: フロントエンドで定期的にセッション情報を取得し、ユーザーに警告

### 2. クイズのダミー選択肢不足時の処理
- **問題**: 過去プレイヤーが少ない場合、ダミー選択肢が重複する可能性
- **現状**: システム汎用回答で補完
- **改善案**: より多様な汎用回答を事前定義

### 3. 画像類似度判定の未統合
- **問題**: ml_handlerは実装済みだがルーティング未設定
- **現状**: Web選択（"この場所にいることにする"）で代替
- **対策**: 統合完了まで代替手段で運用可能

### 4. CORS設定が開発用
- **問題**: `Access-Control-Allow-Origin: *` で全オリジン許可
- **影響**: 本番環境でセキュリティリスク
- **対策**: 本番デプロイ前に許可オリジンを制限（環境変数化推奨）

### 5. トランザクション未実装
- **問題**: 複数テーブル更新時のロールバック処理なし
- **影響**: データ不整合のリスク（S6進捗+クイズ生成など）
- **対策**: 重要な処理にトランザクションを追加

---

## 開発のヒント

### コーディング規約
- **言語**: Go 1.25.1
- **フォーマット**: `go fmt` 必須（`make fmt`）
- **コメント**: 日本語推奨
- **命名**: 英語（関数名、変数名）、説明は日本語

### アーキテクチャガイドライン
- **依存方向**: 必ず内側（Domain層）に向ける
- **ビジネスロジック**: Domain層に集約
- **データベースアクセス**: Repository経由のみ
- **エラーハンドリング**: ドメインエラーを定義し、Handler層でHTTPステータスに変換

### よくある作業パターン

#### 新しいエンティティ追加
1. Value Object定義（`internal/domain/valueobject/`）
2. Entity定義（`internal/domain/model/`）
3. Repository Interface定義（`internal/domain/repository/`）
4. Repository Implementation（`internal/infrastructure/persistence/`）
5. Use Case実装（`internal/application/usecase/`）
6. DTO定義（`internal/interfaces/http/dto/`）
7. Handler実装（`internal/interfaces/http/handler/`）
8. ルーティング追加（`cmd/server/main.go`）
9. マイグレーションファイル作成（`migrations/`）

#### デバッグ方法
```bash
# サーバーログ確認
make run

# MySQLデータ確認
mysql -h localhost -u user -p kotti_game
SELECT * FROM sessions;

# gRPC疎通確認（画像認識）
grpcurl -plaintext localhost:50051 list
```

### 参考資料
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://www.domainlanguage.com/ddd/)
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- [golang-migrate](https://github.com/golang-migrate/migrate)

---

## 更新履歴

- **2025-11-07**: 初版作成（PostgreSQLからMySQL 8.0への移行完了時点）

---

## 連絡先・サポート

質問や問題が発生した場合は、以下を参照してください：

- **SETUP.md**: セットアップ手順の詳細
- **CLAUDE.md**: プロジェクト全体の仕様
- **migrations/**: データベーススキーマ定義

---

**このドキュメントは、後任者が迷わず開発を継続できるよう、現在の実装状況を正確に記録したものです。不明点があれば、各ファイルのコメントも参照してください。**
