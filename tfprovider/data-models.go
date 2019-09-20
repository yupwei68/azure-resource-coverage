package tfprovider

type TerraformConfig struct {
	Clients  *[]ReferencedClient
}

type ReferencedClient struct {
	Package     *GoPackage
	GoSDKClient string
}

type GoPackage struct {
	Alias   string
	Package string
	ApiVersion string
}
