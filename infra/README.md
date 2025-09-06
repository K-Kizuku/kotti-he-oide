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

## カスタムドメイン + HTTPS（Cloudflare）

このリポジトリは、単一のALB上で `/${api}` パスをAPI、`/` をWEBに振り分けます。Cloudflareで管理しているドメインからサブドメイン（例: `kotti.kizuku-hackathon.work`）でHTTPSアクセスさせる手順は以下です。

1. `terraform.tfvars` にサブドメインを設定

   ```hcl
   custom_domain_name = "kotti.kizuku-hackathon.work"
   ```

2. ACM証明書を作成（検証用レコードを取得）

   ```bash
   # 先に証明書のみ作成して検証値を取得
   terraform apply -target=aws_acm_certificate.this
   terraform output -json acm_validation_records
   ```

3. CloudflareでACM検証レコードを追加（DNS Validation）

   - 出力された `acm_validation_records` の各要素 `{ name, type, value }` に従って、Cloudflareの対象ゾーンにDNSレコードを追加します。
   - 種別はほぼ `CNAME` です。複数出る場合はすべて追加してください。

4. 証明書の検証完了とALB HTTPSの作成

   ```bash
   terraform apply
   ```

   - これにより ALB の 443 リスナー（証明書つき）と 80→443 リダイレクトが構成されます。

5. Cloudflareで配信用のCNAMEを作成

   - レコード: `kotti`（= `kotti.kizuku-hackathon.work`）
   - 種別: `CNAME`
   - 値: Terraform出力の `alb_dns_name`（例: `poc-alb-xxxxxxxx.ap-northeast-1.elb.amazonaws.com`）
   - プロキシ: 有効（オレンジ雲）推奨
   - SSL/TLS 暗号化モード: 「フル（厳格）」推奨

6. 動作確認

   - `https://kotti.kizuku-hackathon.work/` でWEB
   - `https://kotti.kizuku-hackathon.work/api/healthz` でAPI

## コスト注意点

- ALB ×1, Fargate タスク ×2, RDS ×1, S3(公開) などが課金対象です。
- 利用後は必ずリソースを削除してください。

## 片付け

```bash
terraform destroy -auto-approve
```
