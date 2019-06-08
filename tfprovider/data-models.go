package tfprovider

type TerraformConfig struct {
	FullPath string
	Imports  []*GoPackage
	Clients  []*ReferencedClient
}

type ReferencedClient struct {
	Package     *GoPackage
	GoSDKClient string
}

type GoPackage struct {
	Alias   string
	Package string
}
