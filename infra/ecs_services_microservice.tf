resource "aws_cloudwatch_log_group" "microservice" {
  name              = "/${var.name_prefix}/microservice"
  retention_in_days = 7
  tags              = local.tags
}

resource "aws_ecs_task_definition" "microservice" {
  family                   = "${var.name_prefix}-microservice"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "512"
  memory                   = "1024"
  execution_role_arn       = aws_iam_role.task_execution.arn
  task_role_arn            = aws_iam_role.microservice_task.arn

  container_definitions = jsonencode([
    {
      name  = "microservice"
      image = "${aws_ecr_repository.microservice.repository_url}:latest"
      portMappings = [{
        containerPort = var.microservice_container_port
        protocol      = "tcp"
      }]
      environment = [
        {
          name  = "S3_BUCKET_NAME"
          value = aws_s3_bucket.microservice.bucket
        },
        {
          name  = "AWS_REGION"
          value = var.region
        },
        {
          name  = "API_SERVER_ENDPOINT"
          value = "http://api:${var.api_container_port}"
        },
        {
          name  = "GRPC_PORT"
          value = tostring(var.microservice_container_port)
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.microservice.name
          awslogs-region        = var.region
          awslogs-stream-prefix = "microservice"
        }
      }
    }
  ])
  tags = local.tags
}

resource "aws_ecs_service" "microservice" {
  name            = "${var.name_prefix}-microservice"
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.microservice.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.private[*].id
    security_groups  = [aws_security_group.microservice.id]
    assign_public_ip = false
  }

  # Cloud Map でサービスディスカバリを有効化
  service_registries {
    registry_arn = aws_service_discovery_service.microservice.arn
  }

  tags = local.tags

  # Service discovery (オプション: 必要に応じて有効化)
  # service_registries {
  #   registry_arn = aws_service_discovery_service.microservice.arn
  # }
}

# Service Discovery (gRPC サービス発見)
resource "aws_service_discovery_private_dns_namespace" "this" {
  name        = "${var.name_prefix}.local"
  description = "Private DNS namespace for microservices"
  vpc         = aws_vpc.this.id
  tags        = local.tags
}

resource "aws_service_discovery_service" "microservice" {
  name = "microservice"

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.this.id
    dns_records {
      ttl  = 10
      type = "A"
    }
  }

  health_check_grace_period_seconds = 30
  tags                             = local.tags
}
