resource "aws_security_group" "alb" {
  name        = "${var.name_prefix}-alb-sg"
  description = "Allow HTTP/HTTPS inbound"
  vpc_id      = aws_vpc.this.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.tags
}

resource "aws_security_group" "ecs" {
  name        = "${var.name_prefix}-ecs-sg"
  description = "Allow from ALB"
  vpc_id      = aws_vpc.this.id

  ingress {
    from_port       = var.api_container_port
    to_port         = var.api_container_port
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  ingress {
    from_port       = var.web_container_port
    to_port         = var.web_container_port
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.tags
}

resource "aws_security_group" "rds" {
  name        = "${var.name_prefix}-rds-sg"
  description = "Allow MySQL from ECS"
  vpc_id      = aws_vpc.this.id

  ingress {
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.tags
}

# マイクロサービス用セキュリティグループ
resource "aws_security_group" "microservice" {
  name        = "${var.name_prefix}-microservice-sg"
  description = "Security group for microservice"
  vpc_id      = aws_vpc.this.id

  # gRPC通信用のポートをAPIサーバーからのアクセスのみ許可
  ingress {
    from_port       = var.microservice_container_port
    to_port         = var.microservice_container_port
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs.id]
  }

  # S3とインターネットへの HTTPS アクセスを許可
  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # APIサーバーとのgRPC通信用
  egress {
    from_port       = var.api_container_port
    to_port         = var.api_container_port
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs.id]
  }

  tags = local.tags
}

# 既存のECSセキュリティグループにマイクロサービスとの通信を追加
resource "aws_security_group_rule" "ecs_to_microservice" {
  type                     = "egress"
  from_port                = var.microservice_container_port
  to_port                  = var.microservice_container_port
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.microservice.id
  security_group_id        = aws_security_group.ecs.id
}

resource "aws_security_group_rule" "microservice_from_ecs" {
  type                     = "ingress"
  from_port                = var.api_container_port
  to_port                  = var.api_container_port
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.ecs.id
  security_group_id        = aws_security_group.microservice.id
}
