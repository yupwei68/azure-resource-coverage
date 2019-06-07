package apispec

func (spec *ApiSpec) Namespaces() []*NamespaceDefinition {
	defs := make([]*NamespaceDefinition, 0)
	for _, ns := range spec.namespaces {
		defs = append(defs, ns)
	}
	return defs
}
