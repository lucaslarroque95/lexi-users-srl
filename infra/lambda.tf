locals {
  foundation = data.terraform_remote_state.foundation.outputs

  postgres_env = {
    POSTGRES_SERVER   = local.foundation.rds_address
    POSTGRES_PORT     = tostring(local.foundation.rds_port)
    POSTGRES_USER     = local.foundation.db_master_username
    POSTGRES_PASSWORD = local.foundation.rds_master_password
    POSTGRES_DB       = local.foundation.db_name
  }
}

resource "aws_lambda_function" "users" {
  for_each = var.users_routes

  function_name = "${var.project}-users-${each.key}"
  role          = local.foundation.users_lambda_role_arn
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  handler       = "bootstrap"
  timeout       = 10
  memory_size   = 256

  filename         = "${path.module}/../dist/${each.key}.zip"
  source_code_hash = filebase64sha256("${path.module}/../dist/${each.key}.zip")

  vpc_config {
    subnet_ids         = local.foundation.private_subnet_ids
    security_group_ids = [local.foundation.lambda_security_group_id]
  }

  environment {
    variables = merge(local.postgres_env, {
      KEYS_DIR    = "/var/task"
      ADMIN_EMAIL = local.foundation.admin_email
    })
  }
}
