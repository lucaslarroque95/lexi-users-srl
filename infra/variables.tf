variable "region" {
  description = "Must match lexi-infra-aws's region"
  type        = string
  default     = "us-east-1"
}

variable "project" {
  description = "Must match lexi-infra-aws's var.project"
  type        = string
  default     = "lexi-serverless"
}

# Fed by routes.auto.tfvars.json, which ../generate_routes.py regenerates
# from every endpoint folder's own endpoint.json. Required (no default): a
# missing var-file should fail loudly, not silently plan to destroy every
# function.
variable "users_routes" {
  type = map(object({
    method = string
    path   = string
  }))
}
