package apispec

func (spec *ApiSpec) getOrCreateNamespace(name string, rel string, typ NamespaceType) *NamespaceDefinition {
	locator := namespaceLocator{name, typ}
	if _, ok := spec.namespaces[locator]; !ok {
		spec.namespaces[locator] = &NamespaceDefinition{
			name,
			rel,
			typ,
			make(map[string]*ProviderDefinition),
		}
	}
	return spec.namespaces[locator]
}

func (ns *NamespaceDefinition) getOrCreateProvider(name string, rel string) *ProviderDefinition {
	if _, ok := ns.Providers[name]; !ok {
		ns.Providers[name] = &ProviderDefinition{
			RelativePath:rel,
			ResourceList:ResourceList{
				Index:0,
				Resources: make(map[int32]*ResourceDefinition),
			},
			ResourceNameMap:make(map[string]int32),
			OpsPathMap:make(map[string]int32),
		}
	}
	return ns.Providers[name]
}

func (pvd *ProviderDefinition) createResource() int32 {
	pvd.ResourceList.Index = pvd.ResourceList.Index+1
	pvd.ResourceList.Resources[pvd.ResourceList.Index] = &ResourceDefinition{
		Versions: make([]VersionDefinition,0),
		Name:     make(map[string]string),
		operations: make(map[string]string),
		OperationReqPath:make(map[string]string),
		}
	return pvd.ResourceList.Index
}
