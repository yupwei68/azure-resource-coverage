package apispec

type ApiSpec struct {
	FullPath string

	namespaces map[namespaceLocator]*NamespaceDefinition
}

type NamespaceDefinition struct {
	Name         string
	RelativePath string
	Type         NamespaceType
	Providers    map[string]*ProviderDefinition
}

type NamespaceType int

const (
	Management NamespaceType = iota
	DataPlane
	ControlPlane

	unknown NamespaceType = -1
)

type ProviderDefinition struct {
	RelativePath string
	Resources    map[string]*ResourceDefinition
}

type ResourceDefinition struct {
	Versions []VersionDefinition

	operations map[string]string
}

type VersionDefinition struct {
	IsPreview  bool
	SDKVersion string
}

type namespaceLocator struct {
	name string
	typ  NamespaceType
}

type swagger struct {
	Paths    map[string]swaggerPath
	xMsPaths map[string]swaggerPath
}

type swaggerPath map[string]interface{}
