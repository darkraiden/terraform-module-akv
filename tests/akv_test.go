package tests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type AzureOptions struct {
	SubscriptionID string
	TenantID       string
}

type TestParameters struct {
	Expected string
	Actual   string
}

type KeyVault struct {
	ID   TestParameters
	Name TestParameters
	URI  TestParameters
}

var terraformDir string = "./fixture/"
var resourceGroupNamePrefix string = "terratest-darkraiden-akv"
var akvNamePrefix string = "terratest-akv"
var azureOptions AzureOptions

func init() {
	godotenv.Load()
	if getSubscriptionID() == "" || getTenantID() == "" {
		panic("both `TEST_SUBSCRIPTION_ID` and `TEST_TENANT_ID` need to be set to run this test file")
	}

	azureOptions.SubscriptionID = getSubscriptionID()
	azureOptions.TenantID = getTenantID()
}

func TestAKV(t *testing.T) {
	var keyVault KeyVault
	uniqueId := random.UniqueId()
	resourceGroupName := fmt.Sprintf("%s-%s", resourceGroupNamePrefix, uniqueId)
	keyVault.Name.Expected = fmt.Sprintf("%s-%s", akvNamePrefix, uniqueId)

	terraformOptions := generateTerraformOptions(map[string]interface{}{
		"resourc e_group_name": resourceGroupName,
		"akv_name":             keyVault.Name.Expected,
	})
	defer terraform.Destroy(t, &terraformOptions)

	_, err := terraform.InitAndApplyAndIdempotentE(t, &terraformOptions)
	assert.Nil(t, err, "idempotency test failed")

	keyVault.ID.Expected = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s", getSubscriptionID(), resourceGroupName, keyVault.Name.Expected)
	keyVault.URI.Expected = strings.ToLower(fmt.Sprintf("https://%s.vault.azure.net/", keyVault.Name.Expected))

	keyVault.Name.Actual = terraform.Output(t, &terraformOptions, "akv_name")
	keyVault.ID.Actual = terraform.Output(t, &terraformOptions, "akv_id")
	keyVault.URI.Actual = terraform.Output(t, &terraformOptions, "akv_uri")

	assert.Equal(t, keyVault.Name.Expected, keyVault.Name.Actual, fmt.Sprintf("invalid Key Vault Name. Expected %s, got %s\n", keyVault.Name.Expected, keyVault.Name.Actual))
	assert.Equal(t, keyVault.ID.Expected, keyVault.ID.Actual, fmt.Sprintf("invalid Key Vault ID. Expected %s, got %s\n", keyVault.ID.Expected, keyVault.ID.Actual))
	assert.Equal(t, keyVault.URI.Expected, keyVault.URI.Actual, fmt.Sprintf("invalid Key Vault URI. Expected %s, got %s\n", keyVault.URI.Expected, keyVault.URI.Actual))
}

func generateTerraformOptions(vars map[string]interface{}) terraform.Options {
	vars["subscription_id"] = azureOptions.SubscriptionID
	vars["tenant_id"] = azureOptions.TenantID

	return terraform.Options{
		TerraformDir: terraformDir,
		Vars:         vars,
	}
}

func getSubscriptionID() string {
	return os.Getenv("TEST_SUBSCRIPTION_ID")
}

func getTenantID() string {
	return os.Getenv("TEST_TENANT_ID")
}
