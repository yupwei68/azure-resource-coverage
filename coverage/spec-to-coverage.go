package coverage

import (
	"sort"

	"github.com/JunyiYi/azure-resource-coverage/apispec"
)

func ToCoverage(spec *apispec.ApiSpec) ResourceCoverage {
	coverage := make(ResourceCoverage, 0)
	for nsn, ns := range spec.Namespaces {
		for pvdn, pvd := range ns.Providers {
			for resn, res := range pvd.Resources {
				coverage = append(coverage, &CoverageEntry{
					nsn,
					ns,
					pvdn,
					pvd,
					resn,
					res,
				})
			}
		}
	}

	// Sort by namespace, then by provider, finally by resource
	sort.SliceStable(coverage, func(i, j int) bool {
		return coverage[i].ResourceName < coverage[j].ResourceName
	})

	sort.SliceStable(coverage, func(i, j int) bool {
		return coverage[i].ProviderName < coverage[j].ProviderName
	})

	sort.SliceStable(coverage, func(i, j int) bool {
		return coverage[i].NamespaceName < coverage[j].NamespaceName
	})

	return coverage
}
