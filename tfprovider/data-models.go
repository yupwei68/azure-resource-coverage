package tfprovider

type TerraformConfig struct {
	FullPath string
	Imports  []*GoPackage
	Clients  []*ReferencedClient
}

type ReferencedClient struct {
	GoSDKNamespace string
	GoSDKClient    string
}

type GoPackage struct {
	Alias   string
	Package string
}
