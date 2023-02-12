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

var terraformDir string = "./fixture/"

type akv struct {
	namePrefix                  string
	expectedID                  string
	expectedName                string
	expectedURI                 string
	expectedRGName              string
	expectedAccessPoliciesCount int
}

func TestAKV(t *testing.T) {
	tests := []struct {
		rgPrefix string
		vars     map[string]interface{}
		akv
	}{
		{
			rgPrefix: "terratest-darkraiden-akv-1",
			vars:     nil,
			akv: akv{
				namePrefix:                  "terratest-akv-1",
				expectedAccessPoliciesCount: 1,
			},
		},
		{
			rgPrefix: "terratest-darkraiden-akv-2",
			vars: map[string]interface{}{
				"is_rbac_auth_enabled": true,
			},
			akv: akv{
				namePrefix:                  "terratest-akv-2",
				expectedAccessPoliciesCount: 0,
			},
		},
	}

	for _, test := range tests {
		uniqueID := strings.ToLower(random.UniqueId())
		test.akv.expectedRGName = fmt.Sprintf("%s-%s", test.rgPrefix, uniqueID)
		test.akv.expectedName = fmt.Sprintf("%s-%s", test.akv.namePrefix, uniqueID)
		test.akv.expectedURI = fmt.Sprintf("https://%s.vault.azure.net/", test.akv.expectedName)

		if test.vars == nil {
			test.vars = make(map[string]interface{})
		}
		test.vars["akv_name"] = test.akv.expectedName
		test.vars["resource_group_name"] = test.akv.expectedRGName

		tt := tthelper.New(t)
		tfOptions := tt.TerraformOptions(terraformDir, test.vars)
		defer terraform.Destroy(t, tfOptions)

		test.akv.expectedID = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s", tt.SubscriptionID, test.akv.expectedRGName, test.akv.expectedName)

		_, err := terraform.InitAndApplyAndIdempotentE(t, tfOptions)
		assert.Nil(t, err, "idempotency test failed")

		actualAKVID := terraform.Output(t, tfOptions, "akv_id")
		actualAKVName := terraform.Output(t, tfOptions, "akv_name")
		actualAKVURI := terraform.Output(t, tfOptions, "akv_uri")
		actualAccessPoliciesCount := terraform.OutputList(t, tfOptions, "access_policies_object_ids")

		assert.Equal(t, test.akv.expectedID, actualAKVID, fmt.Sprintf("invalid Key Vault ID. Expected %s, got %s\n", test.akv.expectedID, actualAKVID))
		assert.Equal(t, test.akv.expectedName, actualAKVName, fmt.Sprintf("invalid Key Vault Name. Expected %s, got %s\n", test.akv.expectedName, actualAKVName))
		assert.Equal(t, test.akv.expectedURI, actualAKVURI, fmt.Sprintf("invalid Key Vault URI. Expected %s, got %s\n", test.akv.expectedURI, actualAKVURI))
		assert.Equal(t, test.akv.expectedAccessPoliciesCount, len(actualAccessPoliciesCount))
	}
}
