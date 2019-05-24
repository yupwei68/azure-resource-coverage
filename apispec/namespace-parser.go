package apispec

func (spec *ApiSpec) getOrCreateNamespace(name string, rel string, typ NamespaceType) *NamespaceDefinition {
	if _, ok := spec.Namespaces[name]; !ok {
		spec.Namespaces[name] = &NamespaceDefinition{
			rel,
			typ,
			make(map[string]*ProviderDefinition),
		}
	}
	return spec.Namespaces[name]
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
		}
	}
	return pvd.Resources[name]
}
