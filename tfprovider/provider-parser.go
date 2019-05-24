package tfprovider

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

func LoadConfig(rootPath string) (*TerraformConfig, error) {
	path := filepath.Join(rootPath, "azurerm", "config.go")
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := &TerraformConfig{
		path,
		make([]*ReferencedClient, 0),
	}

	config := string(buf)
	re := regexp.MustCompile(`(?sU)type ArmClient struct \{.*\}`)
	clients := re.FindString(config)
	lines := strings.Split(clients, "\n")
	for _, l := range lines[1 : len(lines)-1] {
		client, err := parseClient(l)
		if err != nil {
			return nil, err
		}
		if client != nil && client.IsGoClient() {
			if !strings.HasSuffix(client.GoSDKClient, "Client") {
				return nil, fmt.Errorf("Go SDK client %q does not end with 'Client'", client.GoSDKClient)
			}
			result.Clients = append(result.Clients, client)
		}
	}

	return result, nil
}
