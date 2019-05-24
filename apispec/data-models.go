package apispec

type ApiSpec struct {
	FullPath   string
	Namespaces map[string]*NamespaceDefinition
}

type NamespaceDefinition struct {
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

type swagger struct {
	Paths    map[string]swaggerPath
	xMsPaths map[string]swaggerPath
}

type swaggerPath map[string]interface{}
