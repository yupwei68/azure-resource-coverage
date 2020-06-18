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
			rel,
			make(map[string]*ResourceDefinition),
		}
	}
	return ns.Providers[name]
}

func (pvd *ProviderDefinition) getOrCreateResource(name string) *ResourceDefinition {
	if _, ok := pvd.Resources[name]; !ok {
		pvd.Resources[name] = &ResourceDefinition{
			make([]VersionDefinition, 0),
			make(map[string]string),
			"",
		}
	}
	return pvd.Resources[name]
}
