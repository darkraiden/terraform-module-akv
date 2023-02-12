output "akv_name" {
  description = "The name of the Key Vault"
  value       = azurerm_key_vault.this.name
}

output "akv_resource_group_name" {
  description = "The name of the resource group the Key Vault is created in"
  value       = azurerm_key_vault.this.resource_group_name
}

output "akv_id" {
  description = "The ID of the Key Vault"
  value       = azurerm_key_vault.this.id
}

output "akv_uri" {
  description = "The URI of the Key Vault"
  value       = azurerm_key_vault.this.vault_uri
}

output "access_policies_object_ids" {
  description = "The Object IDs of the access policies attached to the Key Vault"
  value       = [for policy in azurerm_key_vault_access_policy.this : policy.object_id]
}
