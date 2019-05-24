package tfprovider

import (
	"fmt"
	"strings"
)

func parseClient(def string) (*ReferencedClient, error) {
	def = strings.TrimSpace(def)
	if def == "" || strings.HasPrefix(def, "//") {
		return nil, nil
	}
	segs := strings.Fields(def)
	if (len(segs) < 2) || (len(segs) > 2 && !strings.HasPrefix(segs[2], "//")) {
		return nil, fmt.Errorf("Invalid client definition %q", def)
	}
	field := segs[0]
	parts := strings.Split(segs[1], ".")
	if len(parts) != 1 && len(parts) != 2 {
		return nil, fmt.Errorf("Invalid client type definition %q", segs[1])
	}
	ns := ""
	sdk := parts[0]
	if len(parts) == 2 {
		ns = parts[0]
		sdk = parts[1]
	}
	return &ReferencedClient{
		field,
		ns,
		sdk,
	}, nil
}

func (c *ReferencedClient) IsGoClient() bool {
	if c.GoSDKNamespace == "" {
		return false
	}
	if c.GoSDKNamespace == "az" && c.GoSDKClient == "Environment" {
		return false
	}
	if c.GoSDKNamespace == "context" && c.GoSDKClient == "Context" {
		return false
	}
	return true
}
