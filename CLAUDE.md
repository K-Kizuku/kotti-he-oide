# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリのコードを扱う際のガイダンスを提供します。

**重要**: このドキュメントの内容は `.kiro/steering/` ディレクトリ内のファイル（product.md、tech.md、structure.md）を常に参照し、最新の仕様に準拠してください。steeringファイルの内容が最も正確で最新の情報源です。

## プロジェクト概要

**赤煉瓦文化館 〜こっちにおいで〜**

福岡市赤煉瓦文化館（現エンジニアカフェ）を舞台とした、体験型Webホラーゲームです。

### プロジェクトの目的
- 赤煉瓦文化館の歴史的連続性（1909→資料館→文化館→現在）を「体験」として理解させる
- 恐怖→解放→内省の感情の揺れを使い、プレイヤー自身の「理想」「生きる意義」の記憶を呼び戻す
- プレイヤーがメッセージを刻むことで、「この建物は受け継がれてきた/自分も受け渡す側になる」という構図を体感させる
- 匿名メッセージの継承により、「誰かがここで生きた」痕跡を蓄積していく

### 想定プレイヤー
- エンジニアカフェの来館者
- 赤煉瓦文化館の歴史紹介イベント参加者
- デジタル/ARコンテンツに抵抗のない20〜40代社会人
- プレイ時間：20〜30分（館内移動込み）
- セッション有効時間：開始から60分

## アーキテクチャ

フルスタックハッカソンプロジェクトで、4つの主要コンポーネントから構成されます：

- **Frontend**: Next.js 15.5.2 React アプリケーション（TypeScript + Turbopack）
- **Server**: Go REST API バックエンド（レイヤードアーキテクチャ + DDD原則）
- **Infra**: AWS ECS Fargate インフラストラクチャ（Terraform IaC）
- **VOICEVOX Server**: AWS EC2上で動作する音声生成サーバー

フロントエンドはNext.js App Routerを使用し、TypeScriptパス設定（`@/*` → `./src/*`）が構成されています。

## ゲームの主要機能

### 時系列体験システム
プレイヤーは複数の時代を行き来します：
- **1942年**: 生命保険会社の診査としての来館、担当者との初回会話
- **1972年**: 資料館時代、7分間の探索パート
- **2002年**: 文化館時代、担当者のメッセージ受け取り
- **2025年**: 現在、自分のメッセージを刻む

### シーン構造（S0〜S9）

#### S0: 起動・注意書き・許可
- カメラ・音声・イヤホンの許可取得
- ホラー演出の注意喚起
- 言語選択（日本語/English）

#### S1: 1942年パート - 担当者との初回会話
- VOICEVOX（青山龍星(しっとり)）による音声生成
- プレイヤー情報の入力（来館方法、普段の活動など）
- カメラプレビュー + 暗いオーバーレイ

#### S2: お気に入りの場所（固定5箇所）
1. 螺旋階段を見上げる高い天井
2. メインホールの暖炉のレンガ
3. 裏玄関の扉の蝶番
4. 入口エントランスの扉
5. 階上応接室のピアノ

#### S3: 移動指示 + ホラー演出
- 2階診査室への導線
- カメラにホラー用フィルター（色味+ノイズ）を重ねる
- ランダムで影・視線・虚な人のシルエットを表示
- SE：低音の環境音 + 「こっちに来てはならない」

#### S4: 診査室 - 内省質問パート
- 固定10問の質問（小学生時代の夢中だったこと、尊敬した人、人生の願望など）
- "なし"や"特にない"などの無回答はバリデーションで弾く
- 回答は逐次サーバーに保存（セッション再開対応）
- 必須質問：「人生の最期に達成したいこと」「名前」

#### S5: 死亡届受理 + 7分制限開始（1972年パート）
- 「死亡届が受理されました」の通知
- 7分カウントダウン開始
- 担当者の存在証明書との交換を要求

#### S6: 存在証明書探索パート（最重要）
- **目標**: 5つの場所すべてで1枚ずつピースを取得（合計5枚）
- **制限時間**: 7分以内
- **到達判定**:
  - 撮影ボタン → カメラ起動 → サーバーで画像類似度判定（閾値0.5〜0.6）
  - 撮影できない場合は「この場所にいることにする」で代替可能
- **クイズシステム**:
  - S4の回答を元に5問のクイズを生成
  - 4択（正解：プレイヤー自身の回答、ダミー：別回答/過去プレイヤー回答/システム汎用）
  - 正解でピース取得、不正解は即再挑戦可能（ホラー演出あり）
- **タイマー猶予**: 最後のピース回答中に時間切れの場合のみ許可

#### S7: 2002年パート - メッセージ受け取り
- S4で答えた「人生の最期に達成したいこと」を一字一句で再入力
- 不一致の場合は再入力を要求

#### S8: メインホール3分探索
- マップ画像から該当箇所を探す
- 3分カウントダウン
- 成功時に過去プレイヤーのメッセージ一覧を表示

#### S9: 2025年 - メッセージ刻み
- 最大120文字のメッセージ入力
- S2の5箇所から刻む場所を選択
- 匿名 + セッションIDのみで保存
- 保存したメッセージは後続プレイヤーに提示される

### 匿名継承システム
- 過去プレイヤーの回答が匿名化され、次のプレイヤーのクイズのダミー選択肢になる
- プレイが続くほど「誰かがここで生きた」痕跡が増える
- IP・端末IDは保持せず、完全匿名

### VOICEVOX統合
- EC2上でVOICEVOXサーバーを常駐
- アプリケーションサーバーからHTTPで音声生成リクエスト
- 音声ファイルURLをクライアントに返却
- 使用キャラクター：青山龍星（しっとり）

### 画像類似度判定システム
- S6での場所到達判定に使用
- クライアントから撮影画像 + place_idを送信
- サーバー側で場所ごとの基準画像と類似度を計算
- 閾値以上で到達判定成功（verified_by: "photo"）
- 失敗時はWeb選択で代替可能（verified_by: "manual"）

### セッション管理
- 初回アクセス時にUUID形式のsession_idを発行
- サーバー側で1時間保持
- S4の回答は逐次保存し、Wi-Fi切断やリロード後も続行可能
- 1時間経過後は新規セッションとしてS0から開始

### ローカライズ対応
- 全テキストをキー化（例：`text.intro.warning.ja`, `text.intro.warning.en`）
- デフォルトは日本語
- 英語リソースは後日実装予定（ファイル構造のみ先行準備）

## New Features

### Camera Filter Feature
- Real-time browser camera video processing
- 5 filter types (retro/horror/serious/VHS/comic)
- Canvas 2D API-based high-speed pixel processing
- Front/back camera switching support

## Backend Architecture (DDD + Layered)

The Go server follows Domain-Driven Design (DDD) and Clean Architecture principles:

```
server/
├── cmd/server/                    # Application entry point
├── internal/
│   ├── interfaces/                # Interface Layer (外部インターフェース層)
│   │   └── http/
│   │       ├── handler/           # HTTP handlers
│   │       └── dto/               # Data Transfer Objects
│   ├── application/               # Application Layer (アプリケーション層)
│   │   └── usecase/               # Use cases
│   ├── domain/                    # Domain Layer (ドメイン層)
│   │   ├── model/                 # Domain entities
│   │   ├── repository/            # Repository interfaces
│   │   ├── service/               # Domain services
│   │   └── valueobject/           # Value objects
│   └── infrastructure/            # Infrastructure Layer (インフラストラクチャ層)
│       └── persistence/           # Data persistence implementations
└── pkg/
    └── errors/                    # Custom error types
```

### Dependency Flow
- **Interfaces** → **Application** → **Domain** ← **Infrastructure**
- Dependencies point inward (Clean Architecture)
- Infrastructure implements domain interfaces

## Development Commands

### Frontend Development
All frontend commands should be run from the `frontend/` directory:

```bash
cd frontend
pnpm dev          # Start development server with Turbopack
pnpm build        # Build for production with Turbopack  
pnpm start        # Start production server
pnpm lint         # Run ESLint (eslint command only)
```

The frontend runs on http://localhost:3000 by default.

### Server Development
All server commands should be run from the `server/` directory:

```bash
cd server
make run          # Run server from root main.go (port 8080)
make run-cmd      # Run server from cmd/server/main.go
make build        # Build binary to bin/server
make test         # Run tests (PoCのため新規テスト作成は不要)
make fmt          # Format code
make lint         # Run linter (必須)
```

The server runs on http://localhost:8080 by default.

### Infrastructure Development
All infrastructure commands should be run from the `infra/` directory:

```bash
cd infra
terraform init      # Initialize Terraform
terraform plan       # Preview infrastructure changes
terraform apply      # Apply infrastructure changes
terraform destroy    # Destroy infrastructure resources
```

### API Endpoints

#### ゲーム関連
- `POST /api/session` - セッション開始（session_id発行）
- `GET /api/session/{session_id}` - セッション情報取得
- `POST /api/session/{session_id}/answers` - S4の回答を逐次保存
- `GET /api/session/{session_id}/answers` - 保存された回答を取得
- `POST /api/session/{session_id}/s6/start` - S6開始（7分タイマー開始）
- `POST /api/session/{session_id}/s6/verify-location` - 場所到達判定（画像類似度）
- `GET /api/session/{session_id}/s6/quiz/{place_id}` - 特定場所のクイズ取得
- `POST /api/session/{session_id}/s6/answer` - クイズ回答送信
- `POST /api/session/{session_id}/message` - 最終メッセージ保存（S9）
- `GET /api/messages` - 過去プレイヤーのメッセージ一覧取得

#### VOICEVOX統合
- `POST /api/voice/generate` - テキストから音声生成（VOICEVOXへのプロキシ）
- レスポンス：音声ファイルのURL

#### 画像類似度判定
- エンドポイント：`POST /api/session/{session_id}/s6/verify-location`
- リクエスト：`multipart/form-data`（画像 + place_id）
- レスポンス：`{"verified": true/false, "similarity": 0.0-1.0}`

#### ヘルスチェック
- `GET /api/healthz` - ヘルスチェック

### Frontend Routes
- `/` - ゲーム起動ページ（QRコードからのランディング）
- `/game` - メインゲーム画面（シーン遷移を管理）
- `/game/s0` - 注意書き・許可取得
- `/game/s1` - 1942年パート（担当者との会話）
- `/game/s2` - お気に入りの場所説明
- `/game/s3` - 移動指示 + ホラー演出
- `/game/s4` - 診査室（内省10問）
- `/game/s5` - 死亡届受理通知
- `/game/s6` - 存在証明書探索（7分タイマー）
- `/game/s7` - 2002年パート（メッセージ受け取り）
- `/game/s8` - メインホール探索（3分）
- `/game/s9` - メッセージ刻み
- `/game/gameover` - ゲームオーバー画面
- `/camera-filters` - カメラフィルターデモページ（開発用）

## Infrastructure Architecture

The infrastructure is deployed on AWS using Terraform and follows a containerized approach with ECS Fargate:

```
infra/
├── alb.tf                  # Application Load Balancer configuration
├── ecr.tf                  # Elastic Container Registry
├── ecs_cluster.tf          # ECS Cluster setup
├── ecs_services_api.tf     # API service configuration
├── ecs_services_web.tf     # Web service configuration
├── outputs.tf              # Terraform outputs
├── providers.tf            # AWS and Random providers
├── rds.tf                  # RDS MySQL 8.0 database
├── s3.tf                   # S3 bucket configuration
├── security.tf             # Security groups and IAM roles
├── variables.tf            # Input variables
├── versions.tf             # Terraform version constraints
└── vpc.tf                  # VPC and networking
```

### Infrastructure Components

- **Compute**:
  - ECS Fargate cluster hosting containerized API and Web services
  - EC2 instance for VOICEVOX server (音声生成専用)
- **Load Balancing**: Application Load Balancer (ALB) with path-based routing
- **Container Registry**: ECR repositories for API and Web images
- **Database**: RDS MySQL 8.0 instance
- **Storage**:
  - S3 bucket for static assets
  - S3 bucket for generated voice files (VOICEVOX output)
- **Networking**: Custom VPC with public/private subnets
- **Security**: Security groups and IAM roles for least privilege access

### Deployment Process

1. Build and push container images to ECR
2. Configure `terraform.tfvars` with required variables
3. Deploy infrastructure using Terraform
4. Access services via ALB DNS name:
   - `/` - Web service (Next.js frontend)
   - `/api/*` - API service (Go backend)

## Technology Stack

- **Frontend**: React 19.1.0, Next.js 15.5.2, TypeScript 5+, ESLint 9
- **Backend**: Go 1.25.1, Standard library HTTP server, DDD + Clean Architecture
- **Infrastructure**: AWS ECS Fargate, ALB, RDS MySQL 8.0, ECR, S3, EC2 (VOICEVOX), Terraform ~> 5.0
- **Database**: MySQL 8.0 (go-sql-driver/mysql), golang-migrate for migrations
- **Package Manager**: pnpm (frontend)
- **Build Tool**: Turbopack (Next.js)
- **CI/CD**: GitHub Actions (ECR/ECS auto-deployment)
- **Containerization**: Docker (multi-stage builds)
- **音声合成**: VOICEVOX（青山龍星(しっとり)）- EC2上で動作
- **画像処理**: Canvas 2D API（フロントエンド）、画像類似度判定（バックエンド）

## Key Libraries

### Backend
- **MySQL Driver**: `github.com/go-sql-driver/mysql` v1.8.1
- **UUID Generation**: `github.com/google/uuid` v1.6.0
- **gRPC**: `google.golang.org/grpc` v1.65.0 (VOICEVOX/画像認識統合用)

### Frontend
- **Camera Processing**: Canvas 2D API, getUserMedia
- **Image Filters**: Custom ImageData processing (Sobel, posterization, etc.)

## Development Guidelines

### アーキテクチャ原則
- DDD原則に従う：ビジネスロジックはドメイン層に配置
- Value Objectsを使ったプリミティブ値の検証（UserID, Email, SessionIDなど）
- リポジトリパターンでデータアクセスを抽象化
- カスタムドメインエラー型を使ったエラーハンドリング
- 依存性注入による疎結合
- 全てのコミュニケーションは日本語で行う
- コード内のコメントも可能な限り日本語で記述する
- 変数名や関数名は英語でも構わないが、説明やドキュメントは日本語とする

## テストポリシー

**重要: このプロジェクトはPoC（Proof of Concept）です。**

- **テストコードの作成は一切不要**です
- 開発速度とプロトタイピングを最優先します
- ユニットテスト、統合テスト、E2Eテストなど、あらゆる種類のテストコードの実装は求められません
- `make test` コマンドは存在しますが、新規テストファイルの作成は不要です

### ゲーム固有の開発ガイドライン

#### セッション管理
- セッションIDは必ずUUID v4を使用
- セッション有効期限は作成時刻 + 60分
- S4の回答は必ず逐次保存（リロード対応のため）
- セッション破棄後は新規セッションとして扱う

#### シーン遷移
- 各シーンは独立したコンポーネントとして実装
- シーン間のデータ受け渡しはセッションIDを通じて行う
- カメラなしモードでも全シーンをプレイ可能にする

#### 画像処理
- 画像類似度判定の閾値は0.5〜0.6（調整可能にする）
- カメラ撮影失敗時は必ず代替手段（Web選択）を提供
- 撮影画像は一時的に保存し、判定後は削除

#### クイズ生成
- S6入室時に5問すべてを事前生成（遅延なし）
- 正解は必ずプレイヤー自身の回答
- ダミー選択肢は以下の優先順位で選択：
  1. プレイヤーの別回答
  2. 過去プレイヤーの匿名回答
  3. システム汎用回答（最終手段）

#### タイマー処理
- S6: 7分カウントダウン（最後のピース回答中のみ猶予あり）
- S8: 3分カウントダウン（猶予なし）
- クライアント・サーバー両方でタイマーを管理
- サーバー側の時刻を正とする

#### ホラー演出
- ジャンプスケアは最大2回まで（S6到達時1回 + 誤答時1回）
- カメラフィルターは常時適用（S3以降）
- 音声SEは音量70%を上限とする

#### 匿名性の保持
- プレイヤーメッセージにはIPアドレス・端末IDを保存しない
- session_idのみを保存（将来的にハッシュ化も検討）
- 過去プレイヤー回答は完全匿名で提供

## Code Quality Requirements (MANDATORY)

**Always run lint checks before committing any code changes:**

### Frontend Lint Check
```bash
cd frontend
pnpm lint    # Must pass without errors
```

### Backend Lint Check
```bash
cd server
make fmt     # Format code
make lint    # Must pass without errors
```

**IMPORTANT**: Committing or pushing code with lint errors is strictly prohibited. The CI/CD pipeline also enforces these checks automatically.

## Database Schema

プロジェクトにはMySQL 8.0用のマイグレーションファイル（`server/migrations/`）が含まれています：

### ゲーム関連テーブル
- **sessions**: ゲームセッション管理（session_id CHAR(36), current_scene, s6_started_at, created_at, expires_at）
- **session_answers**: S4の内省質問への回答（session_id, question_id, answer_text, answered_at）
- **session_s6_progress**: S6の探索進捗（session_id, place_id, verified_by, quiz_id, answered, correct）
- **player_messages**: プレイヤーが刻んだメッセージ（session_id（匿名化）, message_text, place_id, created_at）
- **quiz_questions**: 生成されたクイズ（quiz_id, session_id, place_id, question_text, option_1-4, answer_index）
- **location_images**: 場所の基準画像（place_id, image_data, created_at）

### データベース特徴
- UUIDは CHAR(36) として保存（MySQL 8.0互換）
- タイムスタンプは DATETIME 型を使用
- 文字セット：utf8mb4（絵文字対応）
- エンジン：InnoDB（トランザクションサポート）
- マイグレーション管理：golang-migrate

## CI/CD Pipeline

- **Trigger**: Push to `main` branch
- **Frontend**: Builds and deploys when `frontend/**` changes
- **Backend**: Builds and deploys when `server/**` changes
- **Process**: Docker build → ECR push → ECS service update
- **Requirements**: AWS credentials in GitHub Secrets