resource "aws_ecr_repository" "api" {
  name                 = "api"
  image_tag_mutability = "MUTABLE"
  tags                 = local.tags

  # lifecycle_policy {
  #   policy = jsonencode({
  #     rules = [{
  #       rulePriority = 1
  #       description  = "Keep last 5 images"
  #       selection = {
  #         tagStatus   = "any"
  #         countType   = "imageCountMoreThan"
  #         countNumber = 5
  #       }
  #       action = {
  #         type = "expire"
  #       }
  #     }]
  #   })
  # }
}

resource "aws_ecr_repository" "web" {
  name                 = "web"
  image_tag_mutability = "MUTABLE"
  tags                 = local.tags

  # lifecycle_policy {
  #   policy = jsonencode({
  #     rules = [{
  #       rulePriority = 1
  #       description  = "Keep last 5 images"
  #       selection = {
  #         tagStatus   = "any"
  #         countType   = "imageCountMoreThan"
  #         countNumber = 5
  #       }
  #       action = {
  #         type = "expire"
  #       }
  #     }]
  #   })
  # }
}

resource "aws_ecr_repository" "microservice" {
  name                 = "microservice"
  image_tag_mutability = "MUTABLE"
  tags                 = local.tags

  # lifecycle_policy {
  #   policy = jsonencode({
  #     rules = [{
  #       rulePriority = 1
  #       description  = "Keep last 5 images"
  #       selection = {
  #         tagStatus   = "any"
  #         countType   = "imageCountMoreThan"
  #         countNumber = 5
  #       }
  #       action = {
  #         type = "expire"
  #       }
  #     }]
  #   })
  # }
}
