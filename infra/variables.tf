variable "region" {
  type    = string
  default = "ap-northeast-1"
}

variable "name_prefix" {
  type    = string
  default = "poc"
}

variable "db_name" {
  type    = string
  default = "appdb"
}

variable "db_username" {
  type    = string
  default = "app"
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "api_container_port" {
  type    = number
  default = 8080
}

variable "web_container_port" {
  type    = number
  default = 3000
}

variable "api_container_image" {
  type    = string
  default = "api:latest"
}

variable "web_container_image" {
  type    = string
  default = "web:latest"
}

variable "microservice_container_port" {
  type    = number
  # gRPC 標準ポートに合わせる
  default = 50051
}

variable "microservice_container_image" {
  type    = string
  default = "microservice:latest"
}

locals {
  tags = {
    Project     = var.name_prefix
    Owner       = "dev"
    Environment = "dev"
  }
}

variable "custom_domain_name" {
  description = "Cloudflareで管理しているカスタムドメイン(サブドメイン含む)。例: kotti.kizuku-hackathon.work"
  type        = string
}
