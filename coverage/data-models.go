package coverage

import (
	"github.com/JunyiYi/azure-resource-coverage/apispec"
)

type ResourceCoverage struct {
	Entries []*CoverageEntry

	configuration config
}

type CoverageEntry struct {
	Namespace    *apispec.NamespaceDefinition
	ProviderName string
	Provider     *apispec.ProviderDefinition
	ResourceName string
	Resource     *apispec.ResourceDefinition
	InTerraform  bool
}

type config struct {
	APISpec   apispecConfig `json:"api_spec"`
	Terraform terraformConfig
}

type apispecConfig struct {
	Excludes apispecExcludes
}

type apispecExcludes []apispecExclude

type apispecExclude struct {
	Namespace string
	Type      string
	Provider  string
	Resource  string
}

type terraformConfig struct {
	Excludes goClientExcludes
	Mappings goSDKMappings
}

type goSDKMappings map[string]goSDKMapping

type goSDKMapping struct {
	Namespace string `json:"api_spec_namespace"`
	Clients   map[string]goClientMapping
}

type goClientMapping struct {
	Resource string `json:"api_spec_resource"`
}

type goClientExcludes []goClientExclude

type goClientExclude struct {
	Package string
}
