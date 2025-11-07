output "alb_dns_name" {
  value = aws_lb.this.dns_name
}

output "web_target_group_arn" {
  value = aws_lb_target_group.web.arn
}

output "api_target_group_arn" {
  value = aws_lb_target_group.api.arn
}

output "acm_domain" {
  description = "ACM 証明書の対象ドメイン"
  value       = aws_acm_certificate.this.domain_name
}

output "acm_validation_records" {
  description = "Cloudflareに作成する必要があるDNSレコード一覧"
  value = [
    for dvo in aws_acm_certificate.this.domain_validation_options : {
      name  = dvo.resource_record_name
      type  = dvo.resource_record_type
      value = dvo.resource_record_value
    }
  ]
}

output "ecr_api_url" {
  value = aws_ecr_repository.api.repository_url
}

output "ecr_web_url" {
  value = aws_ecr_repository.web.repository_url
}

output "rds_endpoint" {
  value = aws_db_instance.this.endpoint
}

output "s3_website_endpoint" {
  value = aws_s3_bucket_website_configuration.assets.website_endpoint
}

output "voicevox_ec2_private_ip" {
  description = "VOICEVOX EC2インスタンスのプライベートIPアドレス"
  value       = aws_instance.voicevox.private_ip
}

output "voicevox_ec2_public_ip" {
  description = "VOICEVOX EC2インスタンスのパブリックIPアドレス"
  value       = aws_instance.voicevox.public_ip
}

output "voicevox_ec2_endpoint" {
  description = "ECSからアクセスするVOICEVOXエンドポイント"
  value       = "http://${aws_instance.voicevox.private_ip}:${var.voicevox_port}"
}
