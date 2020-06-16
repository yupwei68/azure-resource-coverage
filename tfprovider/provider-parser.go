package tfprovider

import (
	"os"
	"path/filepath"
)

func LoadConfig(rootPath string) (*TerraformConfig, error) {
	paths, err := getServicePaths(rootPath)
	if err != nil {
		return nil, err
	}

	var clients []ReferencedClient

	for _, path := range paths {
		imports, pkgErr := parsePackages(path)
		if pkgErr != nil {
			return nil, pkgErr
		}

		clientErr := parseClients(path, imports, &clients)
		if clientErr != nil {
			return nil, clientErr
		}
	}
	return &TerraformConfig{
		&clients,
	}, nil
}

func getServicePaths(rootPath string) ([]string, error) {
	serviceListPath := filepath.Join(rootPath, "azurerm", "internal", "services")
	paths := make([]string, 0)
	err := filepath.Walk(serviceListPath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.Name() == "client.go" {
			paths = append(paths, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return paths, err
}