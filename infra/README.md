# PoC インフラ (Terraform)

このディレクトリには、ECS(Fargate)上で API(Go) と WEB(Next.js) を動作させるための最小構成 Terraform コードが含まれます。コストを最優先し、単一の ALB や最小限のリソースで構成されています。

## 前提ツール

- [Terraform](https://www.terraform.io/) >= 1.6
- [AWS CLI](https://aws.amazon.com/cli/)
- [Docker](https://www.docker.com/)

## 初期セットアップ

1. **ECR へログインしイメージをビルド & プッシュ**

   ```bash
   # 例: 東京リージョン
   aws ecr get-login-password --region ap-northeast-1 \
     | docker login --username AWS --password-stdin <ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com

   # API
   docker build -t api:latest ./api
   docker tag api:latest <ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com/api:latest
   docker push <ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com/api:latest

   # WEB
   docker build -t web:latest ./web
   docker tag web:latest <ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com/web:latest
   docker push <ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com/web:latest
   ```

2. **`terraform.tfvars` を作成**

   ```hcl
   db_password        = "your-secret"
   api_container_image = "<ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com/api:latest"
   web_container_image = "<ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com/web:latest"
   ```

3. **Terraform 実行**

   ```bash
   terraform init
   terraform apply -auto-approve
   ```

4. **動作確認**

   出力された `alb_dns_name` にブラウザでアクセスして動作を確認します。
   - `/` : WEB サービス
   - `/api/healthz` : API サービス

## コスト注意点

- ALB ×1, Fargate タスク ×2, RDS ×1, S3(公開) などが課金対象です。
- 利用後は必ずリソースを削除してください。

## 片付け

```bash
terraform destroy -auto-approve
```
