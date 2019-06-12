package coverage

import (
	"github.com/JunyiYi/azure-resource-coverage/apispec"
	"github.com/JunyiYi/azure-resource-coverage/tfprovider"
)

func (operations apispecOperations) supportedOperations(entry *CoverageEntry) string {
	for _, op := range operations {
		matched := true
		if op.Namespace != "" && op.Namespace != "*" && op.Namespace != entry.Namespace.Name {
			matched = false
		}
		if op.Provider != "" && op.Provider != "*" && op.Provider != entry.ProviderName {
			matched = false
		}
		if op.Resource != "" && op.Resource != "*" && op.Resource != entry.ResourceName {
			matched = false
		}
		if matched {
			opStr := ""
			if entry.Resource.SupportAnyOperation(op.Create) {
				opStr += "C"
			}
			if entry.Resource.SupportAnyOperation(op.Read) {
				opStr += "R"
			}
			if entry.Resource.SupportAnyOperation(op.Update) {
				opStr += "U"
			}
			if entry.Resource.SupportAnyOperation(op.Delete) {
				opStr += "D"
			}
			if entry.Resource.SupportAnyOperation(op.List) {
				opStr += "L"
			}
			return opStr
		}
	}
	return ""
}

func (excludes apispecExcludes) isExcluded(entry *CoverageEntry) bool {
	for _, excl := range excludes {
		matched := true
		if excl.Namespace != "" && excl.Namespace != "*" && excl.Namespace != entry.Namespace.Name {
			matched = false
		}
		if excl.Type != "" && excl.Type != "*" && configStringToNsType(excl.Type) != entry.Namespace.Type {
			matched = false
		}
		if excl.Provider != "" && excl.Provider != "*" && excl.Provider != entry.ProviderName {
			matched = false
		}
		if excl.Resource != "" && excl.Resource != "*" && excl.Resource != entry.ResourceName {
			matched = false
		}
		if matched {
			return true
		}
	}
	return false
}

func (excludes goClientExcludes) isExcluded(client *tfprovider.ReferencedClient) bool {
	for _, excl := range excludes {
		matched := true
		if excl.Package != "" && excl.Package != "*" && excl.Package != client.Package.Package {
			matched = false
		}
		if matched {
			return true
		}
	}
	return false
}

func (mapping goSDKMappings) getNamespace(pkg *tfprovider.GoPackage) string {
	if ns, ok := mapping[pkg.Package]; ok && ns.Namespace != "" {
		return ns.Namespace
	}
	return pkg.BaseName()
}

func (mapping goSDKMappings) getResource(sdkClient string, pkg *tfprovider.GoPackage) string {
	if ns, ok := mapping[pkg.Package]; ok {
		if res, ok := ns.Clients[sdkClient]; ok {
			return res.Resource
		}
	}
	return sdkClient
}

func configStringToNsType(s string) apispec.NamespaceType {
	if s == "data" {
		return apispec.DataPlane
	} else if s == "control" {
		return apispec.ControlPlane
	} else if s == "management" {
		return apispec.Management
	}
	return apispec.Unknown
}
