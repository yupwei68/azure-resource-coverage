package tfprovider

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

// Regular expression should match the following:
//   * <package>.New<client>ClientWithBaseURI(
//   * <package>.NewClientWithBaseURI(
//   * <package>.New<client>Client(
//   * <package>.NewClient(
var clientRe = regexp.MustCompile(`:=\s*(?P<package>[^.:=]+)\.New(?P<client>[a-zA-Z_0-9]*?)Client(?:WithBaseURI)?\(`)

func parseClients(path string, packages []*GoPackage) ([]*ReferencedClient, error) {
	refs, err := ToReferenceMap(packages)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure Go client references in %q: %v", path, err)
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure Go client references in %q: %v", path, err)
	}

	clients := make([]*ReferencedClient, 0)
	for _, match := range clientRe.FindAllStringSubmatch(string(buf), -1) {
		client, err := parseClientReference(match, refs)
		if err != nil {
			return nil, fmt.Errorf("Cannot parse Azure Go client references in %q: %v", path, err)
		}
		clients = append(clients, client)
	}
	return clients, nil
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
