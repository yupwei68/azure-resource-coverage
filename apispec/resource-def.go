package apispec

import "strings"

func (resource *ResourceDefinition) Operations() []string {
	operations := make([]string, 0)
	for _, op := range resource.operations {
		operations = append(operations, op)
	}
	return operations
}

func (res *ResourceDefinition) SupportOperation(operation string) bool {
	_, ok := res.operations[strings.ToLower(operation)]
	return ok
}

func (res *ResourceDefinition) SupportAnyOperation(operations []string) bool {
	for _, op := range operations {
		if res.SupportOperation(op) {
			return true
		}
	}
	return false
}
