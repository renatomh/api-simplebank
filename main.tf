# Specify the provider (AWS)
provider "aws" {
  region = "us-east-2" # Update with your desired region
}

# Creating a IAM user to access the resources
resource "aws_iam_user" "user" {
  name = "simple-bank-user" # Replace with your desired username
}

# Creating an access key for the newly created IAM user
# Access key ID and secret will be available at "terraform.tfstate" file
resource "aws_iam_access_key" "user" {
  user = aws_iam_user.user.name
}

# Creating a private ECR repository
# URL and credentials will be available at "terraform.tfstate" file
resource "aws_ecr_repository" "ecr_repository" {
  name = "simplebank"
}

# Creating a policy to the new user, so it can access ECR resources
resource "aws_iam_policy" "ecr_policy" {
  name        = "simple-bank-ecr-policy" # Choose a name for the policy
  description = "Policy to allow managing ECR"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:DescribeImages",
          "ecr:ListImages",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload",
          "ecr:PutImage"
        ],
        Effect   = "Allow",
        Resource = "*"
      }
    ]
  })
}

# Attaching the new policy to the user
resource "aws_iam_policy_attachment" "ecr_policy" {
  name       = "ecr-policy-attachment"
  policy_arn = aws_iam_policy.ecr_policy.arn
  users      = [aws_iam_user.user.name]
}

# Defining a security group to access the EC2 instance to be created
resource "aws_security_group" "security_group" {
  name        = "simple-bank-security-group"
  description = "Simple Bank Security Group"

  # Allow PostgreSQL (port 5432) access from anywhere
  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "PostgreSQL"
  }

  # Outbound rule to allow all traffic to all destinations
  egress {
    from_port   = 0
    to_port     = 0    # Allow all ports
    protocol    = "-1" # Allow all protocols
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Defining variables for the RDS instance user password
variable "db_password" {
  description = "User's password for the RDS instance (write down this password, since it'll be the database user's password)"
  type        = string
}

# Creating a new RDS database server instance
resource "aws_db_instance" "rds" {
  identifier            = "simple-bank"
  allocated_storage     = 20
  storage_type          = "gp3"
  engine                = "postgres"
  engine_version        = "12.17"
  instance_class        = "db.t2.micro"
  db_name               = "simple_bank"
  username              = "root"
  password              = var.db_password
  parameter_group_name  = "default.postgres12"
  publicly_accessible   = true
  skip_final_snapshot   = true
  storage_encrypted     = false
  copy_tags_to_snapshot = true

  # Enabling automated backups
  backup_retention_period = 7
  backup_window           = "03:50-04:20"

  # Adding the security group to the instance
  vpc_security_group_ids = [aws_security_group.security_group.id]
}

# Creating the secrets
resource "aws_secretsmanager_secret" "env_secret" {
  name        = "simple_bank"
  description = "Environment variables and secrets for Simple Bank"
}

# Defining the fields and values for the secrets
resource "aws_secretsmanager_secret_version" "env_secret_version" {
  secret_id     = aws_secretsmanager_secret.env_secret.id
  secret_string = <<EOT
{
  "DB_DRIVER": "postgres",
  "DB_SOURCE": "postgresql://root:${var.db_password}@${aws_db_instance.rds.endpoint}:5432/simple_bank?sslmode=disable",
  "SERVER_ADDRESS": "0.0.0.0:8080",
  "TOKEN_SYMMETRIC_KEY": "34984392010eaaac519278b232d94506224de14db0c4d5d77af8499d1b4e8f5c8375d457aeee187d75cb11305c3a2cea31723ea03aba5bd910967a335d8dcfed",
  "ACCESS_TOKEN_DURATION": "15m",
  "REFRESH_TOKEN_DURATION": "24h"
}
EOT
}
