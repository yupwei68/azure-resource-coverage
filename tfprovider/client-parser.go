package tfprovider

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var clientRe = regexp.MustCompile(`:=\s*(?P<package>[^.:=]+)\.New(?P<client>[a-zA-Z_0-9]+)WithBaseURI`)

func parseClients(path string) ([]*ReferencedClient, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	clients := make([]*ReferencedClient, 0)
	for _, match := range clientRe.FindAllStringSubmatch(string(buf), -1) {
		client, err := parseClientReference(match)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func parseClientReference(def []string) (*ReferencedClient, error) {
	match, captures, err := toNamedCaptures(def, clientRe)
	if err != nil {
		return nil, err
	}

	pkg, ok := captures["package"]
	if !ok || pkg == "" {
		return nil, fmt.Errorf("Cannot parse Go SDK package in %q", match)
	}

	client, ok := captures["client"]
	if !ok || client == "" {
		return nil, fmt.Errorf("Cannot parse Go SDK client in %q", match)
	}

	if !strings.HasSuffix(client, "Client") {
		return nil, fmt.Errorf("Go SDK client %q does not end with 'Client'", client)
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
