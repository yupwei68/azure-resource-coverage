package apispec

func (resource *ResourceDefinition) AdditionalOperations() []string {
	operations := make([]string, 0)
	for _, op := range resource.operations {
		operations = append(operations, op)
	}
	return operations
}

func (res *ResourceDefinition) SupportCreate() bool {
	return res.supportOperation("create") || res.supportOperation("createorupdate")
}

func (res *ResourceDefinition) SupportRead() bool {
	return res.supportOperation("get")
}

func (res *ResourceDefinition) SupportUpdate() bool {
	return res.supportOperation("update") || res.supportOperation("createorupdate")
}

func (res *ResourceDefinition) SupportDelete() bool {
	return res.supportOperation("delete")
}

func (res *ResourceDefinition) SupportList() bool {
	return res.supportOperation("list")
}

func (res *ResourceDefinition) supportOperation(operation string) bool {
	if _, ok := res.operations[operation]; ok {
		return true
	}
	return false
}
