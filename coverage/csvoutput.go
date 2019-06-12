package coverage

import (
	"fmt"

	"github.com/JunyiYi/azure-resource-coverage/apispec"
)

func (cov *ResourceCoverage) OutputCsv() {
	fmt.Println("Namespace,Type,Provider,Resource,Operations,Terraform Support")
	for _, entry := range cov.Entries {
		nsType := nsTypeToOutputString(entry.Namespace.Type)
		ops := cov.configuration.APISpec.Operations.supportedOperations(entry)
		tfStatus := ""
		if entry.InTerraform {
			tfStatus = "yes"
		}
		fmt.Printf("%s,%s,%s,%s,%s,%s\n", entry.Namespace.Name, nsType, entry.ProviderName, entry.ResourceName, ops, tfStatus)
	}
}

func nsTypeToOutputString(typ apispec.NamespaceType) string {
	switch typ {
	case apispec.Management:
		return "Management"
	case apispec.DataPlane:
		return "Data Plane"
	case apispec.ControlPlane:
		return "Control Plane"
	default:
		return "<unknown>"
	}
}
