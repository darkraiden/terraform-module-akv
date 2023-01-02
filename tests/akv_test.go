package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

var terraformDir string = "./fixture/"
var resourceGroupName string = "terratest-darkraiden-akv"
var akvName string = "terratest-akv"

func init() {
	if os.Getenv("TEST_SUBSCRIPTION_ID") == "" || os.Getenv("TEST_TENANT_ID") == "" {
		panic("both `TEST_SUBSCRIPTION_ID` and `TEST_TENANT_ID` need to be set to run this test file")
	}
}

func TestAKV(t *testing.T) {
	uniqueId := random.UniqueId()
	fullResourceGroupName := fmt.Sprintf("%s-%s", resourceGroupName, uniqueId)
	fullAkvName := fmt.Sprintf("%s-%s", akvName, uniqueId)

	terraformOptions := generateTerraformOptions(map[string]interface{}{
		"resource_group_name": fullResourceGroupName,
		"akv_name":            fullAkvName,
	})
	defer terraform.Destroy(t, &terraformOptions)

	_, err := terraform.InitAndApplyAndIdempotentE(t, &terraformOptions)
	assert.Nil(t, err, "idempotency test failed")

	expectedAkvID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s", getSubscriptionID(), fullResourceGroupName, fullAkvName)
	expectedAkvURI := fmt.Sprintf("https://%s.vault.azure.net/", fullAkvName)

	actualAkvName := terraform.Output(t, &terraformOptions, "akv_name")
	actualAkvID := terraform.Output(t, &terraformOptions, "akv_id")
	actualAkvURI := terraform.Output(t, &terraformOptions, "akv_uri")

	assert.Equal(t, fullAkvName, actualAkvName, fmt.Sprintf("invalid Key Vault Name. Expected %s, got %s\n", fullAkvName, actualAkvID))
	assert.Equal(t, expectedAkvID, actualAkvID, fmt.Sprintf("invalid Key Vault ID. Expected %s, got %s\n", expectedAkvID, actualAkvID))
	assert.Equal(t, expectedAkvURI, actualAkvURI, fmt.Sprintf("invalid Key Vault URI. Expected %s, got %s\n", expectedAkvURI, actualAkvURI))
}

func generateTerraformOptions(vars map[string]interface{}) terraform.Options {
	vars["subscription_id"] = os.Getenv("TEST_SUBSCRIPTION_ID")
	vars["tenant_id"] = os.Getenv("TEST_TENANT_ID")

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
