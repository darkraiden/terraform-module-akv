resource "azurerm_key_vault" "this" {
  tenant_id                   = var.tenant_id
  resource_group_name         = var.resource_group_name
  location                    = var.location
  
  name                        = var.akv_name
  sku_name = var.sku_name

  enabled_for_disk_encryption = var.enabled_for_disk_encryption
  soft_delete_retention_days  = var.soft_delete_retention_days
  purge_protection_enabled    = var.purge_protection_enabled
}
