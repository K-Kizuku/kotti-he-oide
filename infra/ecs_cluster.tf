resource "aws_ecs_cluster" "this" {
  name = "${var.name_prefix}-cluster"
  tags = local.tags
}

data "aws_iam_policy_document" "ecs_task_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "task_execution" {
  name               = "${var.name_prefix}-task-execution-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume.json
  tags               = local.tags
}

resource "aws_iam_role_policy_attachment" "task_execution" {
  role       = aws_iam_role.task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role" "task" {
  name               = "${var.name_prefix}-task-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume.json
  tags               = local.tags
}

# マイクロサービス用のタスクロール
resource "aws_iam_role" "microservice_task" {
  name               = "${var.name_prefix}-microservice-task-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume.json
  tags               = local.tags
}

# マイクロサービス専用S3バケットへのアクセスポリシー
data "aws_iam_policy_document" "microservice_s3_access" {
  statement {
    effect = "Allow"
    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:DeleteObject",
      "s3:ListBucket"
    ]
    resources = [
      aws_s3_bucket.microservice.arn,
      "${aws_s3_bucket.microservice.arn}/*"
    ]
  }
}

resource "aws_iam_policy" "microservice_s3_access" {
  name        = "${var.name_prefix}-microservice-s3-access"
  description = "Policy for microservice to access its dedicated S3 bucket"
  policy      = data.aws_iam_policy_document.microservice_s3_access.json
  tags        = local.tags
}

resource "aws_iam_role_policy_attachment" "microservice_s3_access" {
  role       = aws_iam_role.microservice_task.name
  policy_arn = aws_iam_policy.microservice_s3_access.arn
}
