resource "random_pet" "this" {
  length = 2
}

resource "aws_s3_bucket" "assets" {
  bucket        = "${var.name_prefix}-static-assets-${random_pet.this.id}"
  force_destroy = true
  tags          = local.tags
}

resource "aws_s3_bucket_ownership_controls" "assets" {
  bucket = aws_s3_bucket.assets.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_public_access_block" "assets" {
  bucket = aws_s3_bucket.assets.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_acl" "assets" {
  depends_on = [
    aws_s3_bucket_ownership_controls.assets,
    aws_s3_bucket_public_access_block.assets
  ]
  bucket = aws_s3_bucket.assets.id
  acl    = "public-read"
}

resource "aws_s3_bucket_website_configuration" "assets" {
  bucket = aws_s3_bucket.assets.id
  index_document {
    suffix = "index.html"
  }
}

resource "aws_s3_bucket_policy" "assets" {
  depends_on = [aws_s3_bucket_public_access_block.assets]
  bucket     = aws_s3_bucket.assets.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = "*"
      Action    = ["s3:GetObject"]
      Resource  = ["${aws_s3_bucket.assets.arn}/*"]
    }]
  })
}

# マイクロサービス専用のS3バケット
resource "aws_s3_bucket" "microservice" {
  bucket        = "${var.name_prefix}-microservice-data-${random_pet.this.id}"
  force_destroy = true
  tags          = local.tags
}

resource "aws_s3_bucket_ownership_controls" "microservice" {
  bucket = aws_s3_bucket.microservice.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_public_access_block" "microservice" {
  bucket = aws_s3_bucket.microservice.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_acl" "microservice" {
  depends_on = [
    aws_s3_bucket_ownership_controls.microservice,
    aws_s3_bucket_public_access_block.microservice
  ]
  bucket = aws_s3_bucket.microservice.id
  acl    = "private"
}
