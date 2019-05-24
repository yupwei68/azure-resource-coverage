package tfprovider

type TerraformConfig struct {
	FullPath string
	Clients  []*ReferencedClient
}

type ReferencedClient struct {
	Name           string
	GoSDKNamespace string
	GoSDKClient    string
}
