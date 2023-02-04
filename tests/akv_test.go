package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/darkraiden/tthelper"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

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

func TestAKV(t *testing.T) {
	var keyVault KeyVault
	uniqueId := random.UniqueId()
	resourceGroupName := fmt.Sprintf("%s-%s", resourceGroupNamePrefix, uniqueId)
	keyVault.Name.Expected = fmt.Sprintf("%s-%s", akvNamePrefix, uniqueId)

	tt := tthelper.New(t)
	terraformOptions := tt.TerraformOptions(terraformDir)
	terraformOptions.Vars["resource_group_name"] = resourceGroupName
	terraformOptions.Vars["akv_name"] = keyVault.Name.Expected
	defer terraform.Destroy(t, terraformOptions)

	_, err := terraform.InitAndApplyAndIdempotentE(t, terraformOptions)
	assert.Nil(t, err, "idempotency test failed")

	keyVault.ID.Expected = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s", tt.SubscriptionID, resourceGroupName, keyVault.Name.Expected)
	keyVault.URI.Expected = strings.ToLower(fmt.Sprintf("https://%s.vault.azure.net/", keyVault.Name.Expected))

	keyVault.Name.Actual = terraform.Output(t, terraformOptions, "akv_name")
	keyVault.ID.Actual = terraform.Output(t, terraformOptions, "akv_id")
	keyVault.URI.Actual = terraform.Output(t, terraformOptions, "akv_uri")

	assert.Equal(t, keyVault.Name.Expected, keyVault.Name.Actual, fmt.Sprintf("invalid Key Vault Name. Expected %s, got %s\n", keyVault.Name.Expected, keyVault.Name.Actual))
	assert.Equal(t, keyVault.ID.Expected, keyVault.ID.Actual, fmt.Sprintf("invalid Key Vault ID. Expected %s, got %s\n", keyVault.ID.Expected, keyVault.ID.Actual))
	assert.Equal(t, keyVault.URI.Expected, keyVault.URI.Actual, fmt.Sprintf("invalid Key Vault URI. Expected %s, got %s\n", keyVault.URI.Expected, keyVault.URI.Actual))
}
