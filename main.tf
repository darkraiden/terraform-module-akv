resource "azurerm_key_vault" "this" {
  name = var.akv_name

  tenant_id           = var.tenant_id
  location            = var.location
  resource_group_name = var.resource_group_name
  sku_name            = var.sku_name

  enabled_for_disk_encryption = var.enabled_for_disk_encryption
  soft_delete_retention_days  = var.soft_delete_retention_days
  purge_protection_enabled    = var.purge_protection_enabled
}

resource "azurerm_key_vault_access_policy" "this" {
  # convert list(object({})) into a map having the `object_id` as the key of each element
  for_each = { for policy in var.access_policies : policy.object_id => policy }

  key_vault_id = azurerm_key_vault.this.id
  tenant_id    = var.tenant_id
  object_id    = each.key

  certificate_permissions = try(each.value.certificate_permissions, [])
  secret_permissions      = try(each.value.secret_permissions, [])
  key_permissions         = try(each.value.key_permissions, [])
}
