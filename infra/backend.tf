# Terraform状態管理用のS3バケットとDynamoDBテーブル
# このファイルは手動でapplyしてからmain構成にbackend設定を追加する

# S3バケット（Terraform状態ファイル用）
resource "aws_s3_bucket" "terraform_state" {
  bucket        = "kotti-he-oide-terraform-state-${random_id.bucket_suffix.hex}"
  force_destroy = false

  tags = {
    Name        = "TerraformStateBucket"
    Project     = "kotti-he-oide"
  }
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
}

# バケットのバージョニングを有効化
resource "aws_s3_bucket_versioning" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id
  versioning_configuration {
    status = "Enabled"
  }
}

# バケットの暗号化を有効化
resource "aws_s3_bucket_server_side_encryption_configuration" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# パブリックアクセスをブロック
resource "aws_s3_bucket_public_access_block" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# 出力値（backend設定で使用）
output "terraform_state_bucket" {
  value       = aws_s3_bucket.terraform_state.bucket
  description = "Name of the S3 bucket for Terraform state"
}