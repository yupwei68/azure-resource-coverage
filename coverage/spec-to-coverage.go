package coverage

import (
	"sort"

	"github.com/JunyiYi/azure-resource-coverage/apispec"
)

func (cov *ResourceCoverage) LoadFromSpec(spec *apispec.ApiSpec) {
	for _, ns := range spec.Namespaces() {
		for pvdn, pvd := range ns.Providers {
			for resn, res := range pvd.Resources {
				entry := &CoverageEntry{
					ns,
					pvdn,
					pvd,
					resn,
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
		return cov.Entries[i].ResourceName < cov.Entries[j].ResourceName
	})

	sort.SliceStable(cov.Entries, func(i, j int) bool {
		return cov.Entries[i].ProviderName < cov.Entries[j].ProviderName
	})

	sort.SliceStable(cov.Entries, func(i, j int) bool {
		return cov.Entries[i].Namespace.Name < cov.Entries[j].Namespace.Name
	})
}
