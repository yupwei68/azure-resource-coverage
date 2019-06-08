package tfprovider

import (
	"path/filepath"
)

func LoadConfig(rootPath string) (*TerraformConfig, error) {
	path := filepath.Join(rootPath, "azurerm", "config.go")

	imports, err := parsePackages(path)
	if err != nil {
		return nil, err
	}

	clients, err := parseClients(path)
	if err != nil {
		return nil, err
	}

	return &TerraformConfig{
		path,
		imports,
		clients,
	}, nil
}
