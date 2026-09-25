# Routes onto the shared API Gateway owned by lexi-infra-aws — this repo
# never creates the API itself, just its own /api/users/* slice of it.
locals {
  base_path = "/api/users"
}

resource "aws_apigatewayv2_integration" "users" {
  for_each = var.users_routes

  api_id                 = local.foundation.api_id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.users[each.key].invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "users" {
  for_each = var.users_routes

  api_id    = local.foundation.api_id
  route_key = "${each.value.method} ${local.base_path}${each.value.path}"
  target    = "integrations/${aws_apigatewayv2_integration.users[each.key].id}"
}

resource "aws_lambda_permission" "users" {
  for_each = var.users_routes

  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.users[each.key].function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${local.foundation.api_execution_arn}/*/*"
}
