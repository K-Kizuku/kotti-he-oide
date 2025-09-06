resource "aws_acm_certificate" "this" {
  domain_name               = var.custom_domain_name
  validation_method         = "DNS"
  subject_alternative_names = []

  lifecycle {
    create_before_destroy = true
  }

  tags = local.tags
}

# Cloudflare 側に CNAME を作成後、このリソースが検証を完了します
resource "aws_acm_certificate_validation" "this" {
  certificate_arn = aws_acm_certificate.this.arn
  validation_record_fqdns = [
    for dvo in aws_acm_certificate.this.domain_validation_options : dvo.resource_record_name
  ]
}

