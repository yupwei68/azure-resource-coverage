package apispec

import (
	"strings"
)

func (resource *ResourceDefinition) Operations() []string {
	operations := make([]string, 0)
	for _, op := range resource.operations {
		operations = append(operations, op)
	}
	return operations
}

func (res *ResourceDefinition) SupportOperation(operation string) bool {
	// k is lowercase swagger ops name, v is ops name
	// k may be `addsServices_listMetricsAverage`, startswtih `list`
	for k,_:= range(res.operations){
		if strings.HasPrefix(k,strings.ToLower(operation)){
			return true
		}
	}
	return false
}

func (res *ResourceDefinition) SupportAnyOperation(operations []string) bool {
	for _, op := range operations {
		if res.SupportOperation(op) {
			return true
		}
	}
	return false
}

func (v VersionDefinition) String() string{
   return v.SDKVersion
}
