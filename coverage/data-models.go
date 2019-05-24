package coverage

import (
	"github.com/JunyiYi/azure-resource-coverage/apispec"
)

type ResourceCoverage []*CoverageEntry

type CoverageEntry struct {
	NamespaceName string
	Namespace     *apispec.NamespaceDefinition
	ProviderName  string
	Provider      *apispec.ProviderDefinition
	ResourceName  string
	Resource      *apispec.ResourceDefinition
	InTerraform   bool
}
