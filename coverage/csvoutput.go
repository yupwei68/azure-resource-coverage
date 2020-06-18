package coverage

import (
	"fmt"
	"github.com/JunyiYi/azure-resource-coverage/apispec"
	"github.com/JunyiYi/azure-resource-coverage/sqlconnect"
	"log"
	"time"
)

func (cov *ResourceCoverage) OutputCsv() {
	fmt.Println("Namespace,Type,Provider,Resource,OperationReqPath, Versions,Operations,Terraform Support")
	for _, entry := range cov.Entries {
		nsType := nsTypeToOutputString(entry.Namespace.Type)
		ops := cov.configuration.APISpec.Operations.supportedOperations(entry)
		tfStatus := ""
		if entry.InTerraform {
			tfStatus = "yes"
		}
		fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s\n", entry.Namespace.Name, nsType, entry.ProviderName, entry.ResourceName,entry.Resource.OperationReqPath, entry.Resource.Versions,ops, tfStatus)
	}
}

func (cov *ResourceCoverage) OutputSqlServer() {
	fmt.Println("Namespace,Type,Provider,Resource,Operations,Terraform Support")
	if sqlconnect.Connect(){
		d := time.Now().Format("2006-01-02")
		for _, entry := range cov.Entries {
			nsType := nsTypeToOutputString(entry.Namespace.Type)
			ops := cov.configuration.APISpec.Operations.supportedOperations(entry)
			tfStatus := ""
			if entry.InTerraform {
				tfStatus = "yes"
			}
			versions := fmt.Sprintf("%s",entry.Resource.Versions)
			createID, err := sqlconnect.CreateCoverage(entry.Namespace.Name,nsType,entry.ProviderName,entry.ResourceName,entry.Resource.OperationReqPath,versions,ops,d,tfStatus=="yes")
			if err != nil {
				log.Fatal("Error creating Employee: ", err.Error())
			}
			fmt.Printf("Inserted ID: %d successfully.\n", createID)
		}
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
