variable "subscription_id" {
  type = string
}

variable "tenant_id" {
  type = string
}

variable "resource_group_name" {
}

variable "akv_name" {
}

variable "is_rbac_auth_enabled" {
  default = false
}
