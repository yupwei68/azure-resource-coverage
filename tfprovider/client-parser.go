package tfprovider

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
)

// Regular expression should match the following:
//   * <package>.New<client>ClientWithBaseURI(
//   * <package>.NewClientWithBaseURI(
//   * <package>.New<client>Client(
//   * <package>.NewClient(
var clientRe = regexp.MustCompile(`:=\s*(?P<package>[^.:=]+)\.New(?P<client>[a-zA-Z_0-9]*?)Client(?:WithBaseURI)?\(`)

func parseClients(path string, packages []*GoPackage, clients *[]ReferencedClient) (error) {
	refs, err := ToReferenceMap(packages)
	if err != nil {
		return fmt.Errorf("Cannot parse Azure go packages in %q: %v", path, err)
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Cannot parse Azure go client file in %q: %v", path, err)
	}

	for _, match := range clientRe.FindAllStringSubmatch(string(buf), -1) {
		newClient, err := parseClientReference(match, refs)
		if err != nil {
			return fmt.Errorf("Cannot parse Azure Go client references in %q: %v", path, err)
		}

		pkgNameSame := false
		pkgPathSame := false
		clientNameSame := false
		for _, existedClient := range *clients {
			if existedClient.Package == newClient.Package {
				pkgNameSame = true
				if existedClient.Package.Package == newClient.Package.Package {
					pkgPathSame = true
					if existedClient.GoSDKClient == newClient.GoSDKClient {
						clientNameSame = true
					}
				}
			}
		}

		i := 1
		if clientNameSame {
			if pkgNameSame {
				if pkgPathSame {
					continue
				} else {
					newClient.Package.Alias = newClient.Package.ReferenceName() + "_" + strconv.Itoa(i)
					*clients = append(*clients, *newClient)
					i += 1
				}
			} else {
				*clients = append(*clients, *newClient)
			}
		} else {
			*clients = append(*clients, *newClient)
		}
	}
	return nil
}

func parseClientReference(def []string, refs map[string]*GoPackage) (*ReferencedClient, error) {
	match, captures, err := toNamedCaptures(def, clientRe)
	if err != nil {
		return nil, err
	}

	pkgRef, ok := captures["package"]
	if !ok || pkgRef == "" {
		return nil, fmt.Errorf("Cannot parse Go SDK package in %q", match)
	}

	pkg, ok := refs[pkgRef]
	if !ok {
		return nil, fmt.Errorf("Package %q not imported", pkgRef)
	}

	client, ok := captures["client"]
	if !ok {
		return nil, fmt.Errorf("Cannot parse Go SDK client in %q", match)
	}

	return &ReferencedClient{
		pkg,
		client,
	}, nil
}

func toNamedCaptures(match []string, re *regexp.Regexp) (string, map[string]string, error) {
	if match == nil {
		return "", nil, fmt.Errorf("match must not be nil")
	}

	captures := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		captures[name] = match[i]
	}
	return match[0], captures, nil
}
