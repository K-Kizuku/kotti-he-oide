# AWS デプロイメントガイド

このドキュメントでは、AWS上にデプロイされている各リソースの環境変数設定方法と、ローカル開発環境からの接続方法を説明します。

## 目次

1. [AWSリソース概要](#awsリソース概要)
2. [環境変数設定一覧](#環境変数設定一覧)
3. [ローカル開発環境からRDSへの接続](#ローカル開発環境からrdsへの接続)
4. [マイグレーション実行方法](#マイグレーション実行方法)
5. [デプロイ後の確認手順](#デプロイ後の確認手順)

---

## AWSリソース概要

このプロジェクトは以下のAWSリソースで構成されています：

| リソース | サービス | 用途 |
|---------|---------|------|
| ECS API Service | ECS Fargate | Go APIバックエンド |
| ECS Web Service | ECS Fargate | Next.jsフロントエンド |
| ECS Microservice | ECS Fargate | 画像認識サービス（gRPC） |
| RDS MySQL | RDS | データベース（MySQL 8.0） |
| EC2 VOICEVOX | EC2 | 音声生成サーバー |
| S3 Assets Bucket | S3 | 静的アセット（パブリック） |
| S3 Microservice Bucket | S3 | 画像認識データ（プライベート） |
| ALB | ALB | ロードバランサー |
| ECR | ECR | コンテナイメージレジストリ |

---

## 環境変数設定一覧

### 1. ECS API Service（Go Backend）

**設定場所**: `infra/ecs_services_api.tf` の `container_definitions` セクション

現在設定されている環境変数：
```hcl
environment = [
  {
    name  = "IMAGE_RECOGNITION_GRPC_ADDR"
    value = "microservice.${var.name_prefix}.local:${var.microservice_container_port}"
  }
]
```

**追加すべき環境変数**（データベース接続、VOICEVOX、S3用）：

```hcl
environment = [
  {
    name  = "IMAGE_RECOGNITION_GRPC_ADDR"
    value = "microservice.${var.name_prefix}.local:${var.microservice_container_port}"
  },
  # データベース接続設定
  {
    name  = "DB_HOST"
    value = aws_db_instance.this.address
  },
  {
    name  = "DB_PORT"
    value = "3306"
  },
  {
    name  = "DB_NAME"
    value = var.db_name
  },
  {
    name  = "DB_USER"
    value = var.db_username
  },
  {
    name  = "DB_PASSWORD"
    value = var.db_password
  },
  # または統合されたDATABASE_URL
  {
    name  = "DATABASE_URL"
    value = "${var.db_username}:${var.db_password}@tcp(${aws_db_instance.this.address}:3306)/${var.db_name}?charset=utf8mb4&parseTime=True&loc=Local"
  },
  # VOICEVOX接続設定
  {
    name  = "VOICEVOX_ENDPOINT"
    value = "http://${aws_instance.voicevox.private_ip}:${var.voicevox_port}"
  },
  # S3設定（音声ファイル保存用）
  {
    name  = "S3_BUCKET_NAME"
    value = aws_s3_bucket.assets.bucket
  },
  {
    name  = "AWS_REGION"
    value = var.region
  },
  # サーバー設定
  {
    name  = "PORT"
    value = tostring(var.api_container_port)
  }
]
```

**重要**: セキュリティのため、本番環境では `DB_PASSWORD` などの機密情報は **AWS Secrets Manager** または **SSM Parameter Store** を使用することを推奨します。

例（Secrets Manager使用）：
```hcl
secrets = [
  {
    name      = "DB_PASSWORD"
    valueFrom = aws_secretsmanager_secret.db_password.arn
  }
]
```

---

### 2. ECS Web Service（Next.js Frontend）

**設定場所**: `infra/ecs_services_web.tf` の `container_definitions` セクション

現在環境変数が設定されていないため、以下を追加：

```hcl
environment = [
  # APIバックエンドのエンドポイント
  {
    name  = "NEXT_PUBLIC_API_URL"
    value = "https://${var.custom_domain_name}/api"
  },
  # または ALB DNS名を使用する場合
  # {
  #   name  = "NEXT_PUBLIC_API_URL"
  #   value = "http://${aws_lb.this.dns_name}/api"
  # },
  # Next.jsのビルド設定
  {
    name  = "NODE_ENV"
    value = "production"
  }
]
```

**注意**: `NEXT_PUBLIC_` プレフィックスは、Next.jsでクライアント側に公開される環境変数に必要です。

---

### 3. ECS Microservice（画像認識サービス）

**設定場所**: `infra/ecs_services_microservice.tf` の `container_definitions` セクション

現在設定されている環境変数：
```hcl
environment = [
  {
    name  = "S3_BUCKET_NAME"
    value = aws_s3_bucket.microservice.bucket
  },
  {
    name  = "AWS_REGION"
    value = var.region
  },
  {
    name  = "API_SERVER_ENDPOINT"
    value = "http://api:${var.api_container_port}"
  },
  {
    name  = "GRPC_PORT"
    value = tostring(var.microservice_container_port)
  }
]
```

**問題**: `API_SERVER_ENDPOINT` の値が `http://api:${var.api_container_port}` となっていますが、Service Discoveryを使用している場合は正しいDNS名に変更が必要です。

**修正案**：
```hcl
{
  name  = "API_SERVER_ENDPOINT"
  value = "http://api.${var.name_prefix}.local:${var.api_container_port}"
}
```

または、APIサービスもService Discoveryに登録する必要があります。

---

### 4. EC2 VOICEVOX Server

**設定場所**: EC2インスタンスのユーザーデータまたはSSM経由で設定

VOICEVOX Dockerコンテナの起動スクリプトに環境変数を追加：

```bash
#!/bin/bash
# VOICEVOXのセットアップ

# システムアップデート
yum update -y

# Dockerインストール
yum install -y docker
systemctl start docker
systemctl enable docker

# VOICEVOXコンテナ起動
docker run -d \
  --name voicevox \
  --restart always \
  -p 50021:50021 \
  -e VOICEVOX_PORT=50021 \
  voicevox/voicevox_engine:cpu-ubuntu20.04-latest

# ヘルスチェック
echo "VOICEVOX started on port 50021"
```

S3への音声ファイルアップロードが必要な場合：
```bash
# 環境変数設定
export S3_BUCKET_NAME="${aws_s3_bucket.assets.bucket}"
export AWS_REGION="${var.region}"
```

**注意**: EC2のIAMロールには既に `AmazonS3FullAccess` が付与されているため、追加のクレデンシャル設定は不要です。

---

### 5. Terraform変数（terraform.tfvars）

**設定場所**: `infra/terraform.tfvars`（このファイルは `.gitignore` に追加すべき）

```hcl
# リージョン設定
region = "ap-northeast-1"

# プロジェクト名プレフィックス
name_prefix = "kotti"

# データベース設定
db_name     = "kotti_game"
db_username = "admin"
db_password = "YOUR_SECURE_PASSWORD_HERE"  # 必ず変更してください

# カスタムドメイン
custom_domain_name = "kotti.kizuku-hackathon.work"

# コンテナポート設定（デフォルトから変更する場合）
# api_container_port         = 8080
# web_container_port         = 3000
# microservice_container_port = 50051
# voicevox_port              = 50021
```

**重要**: `terraform.tfvars` ファイルには機密情報が含まれるため、**絶対にGitにコミットしない**でください。

---

## ローカル開発環境からRDSへの接続

### 前提条件

1. RDSが既にデプロイされている
2. AWS CLIがインストールされ、設定されている
3. 適切なIAM権限がある
4. RDSのセキュリティグループがローカルIPからの接続を許可している

### 方法1: RDSエンドポイントへの直接接続（パブリックアクセス有効時）

**注意**: 現在の設定では `publicly_accessible = false` のため、この方法は使用できません。

### 方法2: Session Manager Port Forwardingを使用（推奨）

RDSがプライベートサブネットにある場合、Session Managerを使ってポートフォワーディングを行います。

#### ステップ1: Bastion EC2インスタンスへの接続

```bash
# VOICEVOXインスタンスをBastionとして使用（またはBastion専用インスタンスを作成）
INSTANCE_ID=$(aws ec2 describe-instances \
  --filters "Name=tag:Name,Values=poc-voicevox-ec2" \
  --query "Reservations[0].Instances[0].InstanceId" \
  --output text)

echo "Instance ID: $INSTANCE_ID"
```

#### ステップ2: ポートフォワーディングセッションの開始

```bash
# RDSエンドポイントを取得
RDS_ENDPOINT=$(cd infra && terraform output -raw rds_endpoint)
RDS_HOST=$(echo $RDS_ENDPOINT | cut -d':' -f1)

# ポートフォワーディング開始
aws ssm start-session \
  --target $INSTANCE_ID \
  --document-name AWS-StartPortForwardingSessionToRemoteHost \
  --parameters "{\"host\":[\"$RDS_HOST\"],\"portNumber\":[\"3306\"],\"localPortNumber\":[\"13306\"]}"
```

これにより、ローカルの `localhost:13306` がRDSの `3306` に転送されます。

#### ステップ3: MySQLクライアントで接続

別のターミナルで：

```bash
# terraform.tfvarsから認証情報を取得（またはSecrets Managerから）
mysql -h 127.0.0.1 -P 13306 -u admin -p
# パスワードを入力
```

### 方法3: VPN接続（長期運用向け）

AWS Client VPNを設定して、VPC内のリソースに直接アクセスできるようにします（このプロジェクトでは未設定）。

---

## マイグレーション実行方法

### ローカルからRDSへのマイグレーション

#### 前提条件

1. `golang-migrate` ツールがインストールされている
2. Session Manager Port Forwardingが有効（上記参照）

#### ステップ1: Terraform出力から接続情報を取得

```bash
cd infra
terraform output rds_endpoint
# 出力例: kotti-db.abc123.ap-northeast-1.rds.amazonaws.com:3306
```

#### ステップ2: 環境変数を設定

```bash
# 方法A: 統合URL形式
export DATABASE_URL="admin:YOUR_PASSWORD@tcp(localhost:13306)/kotti_game"

# 方法B: 個別の環境変数
export RDS_ENDPOINT="localhost"  # Port Forwarding経由
export RDS_PORT="13306"
export RDS_DATABASE="kotti_game"
export RDS_USER="admin"
export RDS_PASSWORD="YOUR_PASSWORD"
```

**注意**: Port Forwarding使用時は `localhost:13306` を指定します。

#### ステップ3: マイグレーション実行

```bash
cd server

# 方法A: DATABASE_URLを使用
make migrate-up

# 方法B: RDS変数を使用
make migrate-up-rds
```

#### マイグレーションのロールバック

```bash
# 最新のマイグレーションを1つロールバック
make migrate-down-rds
```

#### 新しいマイグレーションファイルの作成

```bash
make migrate-create NAME=add_user_profile_table
```

---

### EC2からRDSへのマイグレーション（本番環境）

Session Manager経由でEC2インスタンスに接続し、そこからマイグレーションを実行することもできます。

#### ステップ1: EC2にSSH接続

```bash
# Session Manager経由
aws ssm start-session --target $INSTANCE_ID

# または通常のSSH（EC2がパブリックIPを持つ場合）
ssh ec2-user@EC2_PUBLIC_IP
```

#### ステップ2: 必要なツールのインストール

```bash
# golang-migrateのインストール
sudo yum install -y wget
wget https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz
tar -xzf migrate.linux-amd64.tar.gz
sudo mv migrate /usr/local/bin/
chmod +x /usr/local/bin/migrate

# または Goからインストール
sudo yum install -y golang
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

#### ステップ3: マイグレーションファイルの配置

```bash
# S3からダウンロード、またはGitリポジトリからクローン
git clone https://github.com/YOUR_ORG/kotti-infra.git
cd kotti-infra/server
```

#### ステップ4: 環境変数設定とマイグレーション実行

```bash
# RDS情報を設定（EC2からは直接アクセス可能）
export RDS_ENDPOINT="kotti-db.abc123.ap-northeast-1.rds.amazonaws.com"
export RDS_PORT="3306"
export RDS_DATABASE="kotti_game"
export RDS_USER="admin"
export RDS_PASSWORD="YOUR_PASSWORD"

# マイグレーション実行
make migrate-up-rds
```

---

## デプロイ後の確認手順

### 1. RDS接続確認

```bash
# Port Forwarding経由で接続
mysql -h 127.0.0.1 -P 13306 -u admin -p -e "SHOW DATABASES;"
```

### 2. ECS API Serviceの動作確認

```bash
# ALB経由でヘルスチェック
curl https://kotti.kizuku-hackathon.work/api/healthz
```

または：

```bash
# ALB DNS名を使用
ALB_DNS=$(cd infra && terraform output -raw alb_dns_name)
curl http://$ALB_DNS/api/healthz
```

### 3. ECS Web Serviceの動作確認

```bash
# ブラウザで確認
open https://kotti.kizuku-hackathon.work

# またはcURLで
curl -I https://kotti.kizuku-hackathon.work
```

### 4. VOICEVOXサーバーの動作確認

```bash
# VOICEVOXのヘルスチェック（ECSタスクから実行）
VOICEVOX_IP=$(cd infra && terraform output -raw voicevox_ec2_private_ip)
curl http://$VOICEVOX_IP:50021/version
```

### 5. Microservice（画像認識）の動作確認

```bash
# ECS API ServiceからgRPCでヘルスチェック
# grpcurlを使用（事前にインストール必要）
grpcurl -plaintext microservice.poc.local:50051 list
```

### 6. CloudWatch Logsの確認

```bash
# API Serviceのログ
aws logs tail /poc/api --follow

# Web Serviceのログ
aws logs tail /poc/web --follow

# Microserviceのログ
aws logs tail /poc/microservice --follow
```

---

## トラブルシューティング

### RDSに接続できない

1. **セキュリティグループの確認**:
   ```bash
   aws ec2 describe-security-groups \
     --filters "Name=tag:Name,Values=poc-rds-sg" \
     --query "SecurityGroups[0].IpPermissions"
   ```

2. **RDSのステータス確認**:
   ```bash
   aws rds describe-db-instances \
     --db-instance-identifier poc-db \
     --query "DBInstances[0].DBInstanceStatus"
   ```

### ECSタスクが起動しない

1. **タスク定義の確認**:
   ```bash
   aws ecs describe-task-definition --task-definition poc-api
   ```

2. **サービスイベントの確認**:
   ```bash
   aws ecs describe-services \
     --cluster poc-cluster \
     --services poc-api \
     --query "services[0].events[0:5]"
   ```

### マイグレーションが失敗する

1. **データベース接続文字列の確認**:
   ```bash
   # 接続テスト
   mysql -h 127.0.0.1 -P 13306 -u admin -p -e "SELECT 1;"
   ```

2. **マイグレーションテーブルの確認**:
   ```sql
   USE kotti_game;
   SELECT * FROM schema_migrations;
   ```

---

## セキュリティベストプラクティス

1. **機密情報の管理**:
   - `terraform.tfvars` を `.gitignore` に追加
   - Secrets ManagerまたはSSM Parameter Storeを使用
   - 環境変数をハードコードしない

2. **ネットワークセキュリティ**:
   - RDSは常にプライベートサブネットに配置
   - セキュリティグループで最小限の権限を付与
   - VPCフローログを有効化

3. **アクセス管理**:
   - IAMロールを適切に設定
   - Session Managerを使用してSSHキーを避ける
   - CloudTrailでAPI呼び出しを監査

4. **データベースセキュリティ**:
   - 定期的にパスワードをローテーション
   - 暗号化を有効化（at-rest and in-transit）
   - 自動バックアップを有効化

---

## 参考リンク

- [AWS ECS Documentation](https://docs.aws.amazon.com/ecs/)
- [AWS RDS MySQL Documentation](https://docs.aws.amazon.com/rds/mysql/)
- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [AWS Session Manager](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager.html)

---

## 更新履歴

- 2025-01-XX: 初版作成
