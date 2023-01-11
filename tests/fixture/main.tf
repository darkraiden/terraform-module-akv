data "azurerm_client_config" "current" {}

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

  access_policies = [
    {
      object_id = data.azurerm_client_config.current.object_id
      secret_permissions = [
        "Get",
        "List"
      ]
    }
  ]
}
