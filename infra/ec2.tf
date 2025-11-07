# EC2用のセキュリティグループ（ECSからのHTTPアクセスを許可）
resource "aws_security_group" "voicevox_ec2" {
  name        = "${var.name_prefix}-voicevox-ec2-sg"
  description = "Security group for VOICEVOX EC2 instance"
  vpc_id      = aws_vpc.this.id

  # ECSからのHTTPアクセスを許可
  ingress {
    from_port       = var.voicevox_port
    to_port         = var.voicevox_port
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs.id]
  }

  # SSH接続用（オプション、必要に応じて削除可能）
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # アウトバウンドトラフィックを許可
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.tags, { "Name" = "${var.name_prefix}-voicevox-ec2-sg" })
}

# EC2用のIAMロール
resource "aws_iam_role" "voicevox_ec2" {
  name = "${var.name_prefix}-voicevox-ec2-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })

  tags = local.tags
}

# SSM Session Manager用のポリシーをアタッチ
resource "aws_iam_role_policy_attachment" "voicevox_ec2_ssm" {
  role       = aws_iam_role.voicevox_ec2.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

# S3へのアクセス用ポリシー（VOICEVOX生成音声をS3に保存する場合）
resource "aws_iam_role_policy_attachment" "voicevox_ec2_s3" {
  role       = aws_iam_role.voicevox_ec2.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}

# EC2インスタンスプロファイル
resource "aws_iam_instance_profile" "voicevox_ec2" {
  name = "${var.name_prefix}-voicevox-ec2-profile"
  role = aws_iam_role.voicevox_ec2.name

  tags = local.tags
}

# EC2インスタンス（t3.large + gp3 20GB）
resource "aws_instance" "voicevox" {
  ami                    = "ami-0244ef75e95122fd9"
  instance_type          = "t3.large"
  subnet_id              = aws_subnet.public[0].id
  vpc_security_group_ids = [aws_security_group.voicevox_ec2.id]
  iam_instance_profile   = aws_iam_instance_profile.voicevox_ec2.name

  # gp3 20GBのEBSボリューム設定
  root_block_device {
    volume_type           = "gp3"
    volume_size           = 20
    delete_on_termination = true
    encrypted             = true

    tags = merge(local.tags, { "Name" = "${var.name_prefix}-voicevox-root" })
  }

  # パブリックIPアドレスを自動割り当て
  associate_public_ip_address = true

  # ユーザーデータ（オプション：VOICEVOX起動スクリプトなど）
  user_data = <<-EOF
              #!/bin/bash
              # VOICEVOXのセットアップスクリプトをここに記述
              # 例：
              # yum update -y
              # yum install -y docker
              # systemctl start docker
              # docker run -d -p ${var.voicevox_port}:50021 voicevox/voicevox_engine:cpu-ubuntu20.04-latest
              EOF

  tags = merge(local.tags, { "Name" = "${var.name_prefix}-voicevox-ec2" })
}

# ECSセキュリティグループからEC2へのアウトバウンドルールを追加
resource "aws_security_group_rule" "ecs_to_voicevox" {
  type                     = "egress"
  from_port                = var.voicevox_port
  to_port                  = var.voicevox_port
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.voicevox_ec2.id
  security_group_id        = aws_security_group.ecs.id
  description              = "Allow ECS to access VOICEVOX EC2"
}
