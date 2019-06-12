package coverage

import (
	"fmt"
	"os"
	"strings"

	"github.com/JunyiYi/azure-resource-coverage/tfprovider"
)

func (cov *ResourceCoverage) AnalyzeTerraformCoverage(tf *tfprovider.TerraformConfig) error {
	for _, client := range tf.Clients {
		if !cov.configuration.Terraform.Excludes.isExcluded(client) {
			entry, err := cov.findExactOneEntry(client)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warn: %v\n", err)
			}
			if entry != nil {
				entry.InTerraform = true
			}
		}
	}
	return nil
}

func (cov ResourceCoverage) findExactOneEntry(client *tfprovider.ReferencedClient) (*CoverageEntry, error) {
	var found *CoverageEntry = nil
	for _, entry := range cov.Entries {
		nsClient := cov.configuration.Terraform.Mappings.getNamespace(client.Package)
		if nsClient != entry.Namespace.Name {
			continue
		}

		resClient := strings.ToLower(cov.configuration.Terraform.Mappings.getResource(client.GoSDKClient, client.Package))

		// storage.Accounts <=> StorageAccounts
		resClientAlt := strings.ToLower(client.Package.BaseName()) + resClient

		resClientAlt2 := ""
		if strings.HasSuffix(resClient, "sgroup") {
			// DeploymentsGroup <=> Deployments
			resClientAlt2 = resClient[:len(resClient)-len("group")]
		} else if resClient == "group" {
			// subscriptions.Group <=> Subscriptions
			resClientAlt2 = strings.ToLower(client.Package.BaseName())
		}

		resName := strings.ToLower(entry.ResourceName)
		if resClient != resName && resClientAlt != resName && resClientAlt2 != resName {
			continue
		}

		if found != nil {
			return nil, fmt.Errorf("Found more than one client (%s).%s in coverage", client.Package.ReferenceName(), client.GoSDKClient)
		}
		found = entry
	}
	if found == nil {
		return nil, fmt.Errorf("Cannot find client (%s).%s in coverage", client.Package.ReferenceName(), client.GoSDKClient)
	}
	return found, nil
}
