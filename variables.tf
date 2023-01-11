variable "akv_name" {
  type        = string
  description = "The name of the Key Vault"
  validation {
    condition     = (length(var.akv_name) >= 3 && length(var.akv_name) <= 24) && can(regex("[A-Za-z0-9-]", var.akv_name))
    error_message = "akv_name must contain alphanumeric characters and dashes and must be between 3-24 chars"
  }
}

variable "resource_group_name" {
  type        = string
  description = "The name of the resource group the Key Vault is deployed in"
}

variable "tenant_id" {
  type        = string
  description = "The ID of the tenant the Key Vault is deployed in"
}

variable "location" {
  type        = string
  description = "The Azure region the Key Vault is deployed in"
}

variable "enabled_for_disk_encryption" {
  type        = bool
  description = "[Optional] Whether Key Vault has disk encryption enabled or not"
  default     = true
}

variable "soft_delete_retention_days" {
  type        = number
  description = "[Optional] The number of days that Key Vault items should be retained for once soft-deleted"
  default     = 7
  validation {
    condition     = var.soft_delete_retention_days >= 7 && var.soft_delete_retention_days <= 90
    error_message = "`soft_delete_retention_days` must be a valye between 7 and 90"
  }
}

variable "purge_protection_enabled" {
  type        = bool
  description = "[Optional] Whether Key Vault's purge protection is enabled or not"
  default     = false
}

variable "sku_name" {
  type        = string
  description = "[Optional] The name of the SKU used for this Key Vault"
  default     = "standard"
  validation {
    condition     = can(regex(("standard|premium"), var.sku_name))
    error_message = "only 'standard' and 'premium' are valid values for `sku_name`"
  }
}

variable "access_policies" {
  type = list(object({
    object_id               = string
    certificate_permissions = optional(list(string))
    key_permissions         = optional(list(string))
    secret_permissions      = optional(list(string))
  }))
  description = "[Optional] List of objects to create access policies to be applied to the Key Vault with the object_id as the key"
  default     = []
}
