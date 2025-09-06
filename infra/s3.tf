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
