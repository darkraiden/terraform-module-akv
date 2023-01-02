resource "azurerm_resource_group" "this" {
  name     = var.resource_group_name
  location = "northeurope"
}

module "akv" {
  source              = "../.."
  resource_group_name = azurerm_resource_group.this.name
  akv_name            = var.akv_name
  location            = azurerm_resource_group.this.location
  tenant_id           = var.tenant_id
}
