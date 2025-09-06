output "alb_dns_name" {
  value = aws_lb.this.dns_name
}

output "web_target_group_arn" {
  value = aws_lb_target_group.web.arn
}

output "api_target_group_arn" {
  value = aws_lb_target_group.api.arn
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
