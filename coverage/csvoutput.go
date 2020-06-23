package coverage

import (
	"fmt"
	"github.com/JunyiYi/azure-resource-coverage/apispec"
	"github.com/JunyiYi/azure-resource-coverage/sqlconnect"
	"log"
	"strings"
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
		resNames := []string{}
		opsPaths := []string{}
		for k,_ := range(entry.ResourceName){
			resNames = append(resNames, k)
		}
		for k,_ :=range(entry.Resource.OperationReqPath){
			opsPaths = append(opsPaths,k)
		}
		fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s\n", entry.Namespace.Name, nsType, entry.ProviderName, strings.Join(resNames," | "),strings.Join(opsPaths," | "), entry.Resource.Versions,ops, tfStatus)
	}
}

func (cov *ResourceCoverage) OutputSqlServer() {
	if sqlconnect.Connect(){
		d := time.Now().Format("2006-01-02")
		if _, err :=sqlconnect.DeleteAll();err!=nil{
			log.Fatal("Error deleting all from Coverage:", err.Error())
		}
		for _, entry := range cov.Entries {
			nsType := nsTypeToOutputString(entry.Namespace.Type)
			ops := cov.configuration.APISpec.Operations.supportedOperations(entry)
			tfStatus := ""
			if entry.InTerraform {
				tfStatus = "yes"
			}
			versions := fmt.Sprintf("%s",entry.Resource.Versions)
			resNames := []string{}
			for k,_ := range(entry.ResourceName){
				resNames = append(resNames, k)
			}
			for opsPath,_:=range(entry.Resource.OperationReqPath){
			createID, err := sqlconnect.CreateCoverage(entry.Namespace.Name,nsType,entry.ProviderName,strings.Join(resNames," | ") ,opsPath,versions,ops,d,tfStatus=="yes")
			if err != nil {
				log.Fatal("Error creating Coverage Entry: ", err.Error())
			}
			fmt.Printf("Inserted ID: %d successfully.\n", createID)
			}
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
