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
	Unknown NamespaceType = -1
)

type ResourceList struct {
	Index int32
	Resources map[int32]*ResourceDefinition
}

type ProviderDefinition struct {
	RelativePath string
	// a list of real Resources
	ResourceList  ResourceList
	// a map of ResourceName
	ResourceNameMap map[string]int32
	// a map of Operation Path
	OpsPathMap map[string]int32
}

type ResourceDefinition struct {
	Versions []VersionDefinition
	// a Set of Resource Names
	Name            map[string]string
	operations map[string]string
	// a Set of Operation Path
	OperationReqPath map[string]string
	InTerraform bool
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
	XMsPaths map[string]swaggerPath `json:"x-ms-paths"`
}

type swaggerPath map[string]interface{}

type OpsPathManagement struct {
	resName string
	version VersionDefinition
}
