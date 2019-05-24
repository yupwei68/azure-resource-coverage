package apispec

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func (pvd *ProviderDefinition) parseResourcesFromJson(jsonPath string, ver VersionDefinition) error {
	content, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("Error reading json file %q: %v", jsonPath, err)
	}
	var swagger swagger
	if err := json.Unmarshal(content, &swagger); err != nil {
		return fmt.Errorf("Error parsing json file %q: %v", jsonPath, err)
	}

	for pn, p := range swagger.Paths {
		for k, v := range p {
			if k == "parameters" {
				continue
			}
			op := v.(map[string]interface{})
			if opId, ok := op["operationId"]; ok {
				ids := strings.SplitN(opId.(string), "_", 2)
				if len(ids) != 1 && len(ids) != 2 {
					return fmt.Errorf("Operation ID %q is invalid in %q (%q -> %q)", opId, jsonPath, pn, k)
				}

				resName := "<unknown>"
				opName := ids[0]
				if len(ids) == 2 {
					resName = ids[0]
					opName = ids[1]
				}
				res := pvd.getOrCreateResource(resName)
				res.appendVersion(ver)
				res.operations[strings.ToLower(opName)] = opName
			} else {
				return fmt.Errorf("No operationId found in %q (%q -> %q)", jsonPath, pn, k)
			}
		}
	}
	return nil
}

func (res *ResourceDefinition) appendVersion(ver VersionDefinition) {
	for _, v := range res.Versions {
		if v.IsPreview == ver.IsPreview && v.SDKVersion == ver.SDKVersion {
			return
		}
	}
	res.Versions = append(res.Versions, ver)
}
