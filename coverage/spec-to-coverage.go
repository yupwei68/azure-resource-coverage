package coverage

import (
	"github.com/JunyiYi/azure-resource-coverage/apispec"
	"sort"
)

func (cov *ResourceCoverage) LoadFromSpec(spec *apispec.ApiSpec) {
	for _, ns := range spec.Namespaces() {
		for pvdn, pvd := range ns.Providers {
			for _, res := range pvd.ResourceList.Resources {
				entry := &CoverageEntry{
					ns,
					pvdn,
					pvd,
					res.Name,
					res,
					false,
				}
				if !cov.configuration.APISpec.Excludes.isExcluded(entry) {
					cov.Entries = append(cov.Entries, entry)
				}
			}
		}
	}

	// Sort by namespace, then by provider, finally by resource
	sort.SliceStable(cov.Entries, func(i, j int) bool {
		for firstResNamei:= range(cov.Entries[i].ResourceName){
			for firstResNamej:=range(cov.Entries[j].ResourceName){
				return firstResNamei<firstResNamej
			}
		}
		return false
	})

	sort.SliceStable(cov.Entries, func(i, j int) bool {
		return cov.Entries[i].ProviderName < cov.Entries[j].ProviderName
	})

	sort.SliceStable(cov.Entries, func(i, j int) bool {
		return cov.Entries[i].Namespace.Name < cov.Entries[j].Namespace.Name
	})
}
