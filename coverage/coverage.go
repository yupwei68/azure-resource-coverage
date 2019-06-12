package coverage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func NewCoverage(configPath string) (*ResourceCoverage, error) {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("Error reading resource configuration file %q: %v", configPath, err)
	}
	var config config
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("Error parsing resource configuration file %q: %v", configPath, err)
	}
	return &ResourceCoverage{
		make([]*CoverageEntry, 0),
		config,
	}, nil
}
