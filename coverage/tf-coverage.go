package coverage

import (
	"fmt"
	"os"
	"strings"

	"github.com/JunyiYi/azure-resource-coverage/tfprovider"
)

func (cov *ResourceCoverage) AnalyzeTerraformCoverage(tf *tfprovider.TerraformConfig) error {
	for _, client := range *tf.Clients {
		if !cov.configuration.Terraform.Excludes.isExcluded(client) {
			_, err := cov.findExactOneEntry(client)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warn: %v\n", err)
			}
		}
	}
	return nil
}

func (cov ResourceCoverage) findExactOneEntry(client tfprovider.ReferencedClient) (*CoverageEntry, error) {
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
		} else if nsClient == "resources" && resClient == "groups" {
			// resources.Groups <=> resourceGroups
			resClientAlt2 = strings.TrimRight(nsClient, "s") + resClient
		}

		for resName :=  range(entry.ResourceName){
			if resClient == strings.ToLower(resName) || resClientAlt == strings.ToLower(resName) || resClientAlt2 == strings.ToLower(resName) {
				found = entry
				entry.InTerraform = true
			}
		}

		//if found != nil {
		//	resClientApiVersion := client.Package.ApiVersion
		//	isVersionMatch := false
		//	for _, apiVer := range entry.Resource.Versions {
		//		if apiVer.SDKVersion == resClientApiVersion {
		//			isVersionMatch = true
		//		}
		//	}
		//
		//	if isVersionMatch {
		//		return nil, fmt.Errorf("Found more than one client (%s).%s in coverage", client.Package.ReferenceName(), client.GoSDKClient)
		//	} else {
		//		continue
		//	}
		//}

	}
	if found == nil {
		return nil, fmt.Errorf("Cannot find client (%s).%s in coverage", client.Package.ReferenceName(), client.GoSDKClient)
	}
	return found, nil
}
