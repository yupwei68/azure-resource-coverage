package apispec

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

func (pvd *ProviderDefinition) parseResourcesFromJSON(jsonPath string, ver VersionDefinition) error {
	content, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("Error reading json file %q: %v", jsonPath, err)
	}
	var swagger swagger
	if err := json.Unmarshal(content, &swagger); err != nil {
		return fmt.Errorf("Error parsing json file %q: %v", jsonPath, err)
	}

	allPaths := make(map[string]map[string]interface{})
	for k, v := range swagger.Paths {
		allPaths[k] = v
	}
	for k, v := range swagger.XMsPaths {
		if _, ok := allPaths[k]; ok {
			return fmt.Errorf("Error merging x-ms-paths into paths due to duplicated path %q", k)
		}
		allPaths[k] = v
	}

	for pn, p := range allPaths {
		for k, v := range p {
			if k == "parameters" {
				continue
			}
			op := v.(map[string]interface{})
			if opID, ok := op["operationId"]; ok {
				ids := strings.SplitN(opID.(string), "_", 2)
				if len(ids) != 1 && len(ids) != 2 {
					return fmt.Errorf("Operation ID %q is invalid in %q (%q -> %q)", opID, jsonPath, pn, k)
				}

				resName := "<unknown>"
				opName := ids[0]
				if len(ids) == 2 {
					resName = strings.ToLower(ids[0])
					opName = ids[1]
				}
				// operations doesn't have meaning
				if resName == "operations" || resName == "<unknown>"{
					continue
				}
				re := regexp.MustCompile(`\/{(\w)+}`)
				opsPath := k+re.ReplaceAllString(pn,``)
				if !strings.HasSuffix(opsPath,"/"){
					opsPath = opsPath+"/"
				}

				var resId int32
				if resId,ok =pvd.OpsPathMap[opsPath]; !ok{
					if resId,ok = pvd.ResourceNameMap[resName];!ok{
						resId = pvd.createResource()
						pvd.ResourceNameMap[resName] = resId
					}
					pvd.OpsPathMap[opsPath] = resId
				}

				pvd.ResourceList.Resources[resId].appendVersion(ver)
				// append resourceName to the set
				pvd.ResourceList.Resources[resId].Name[resName] = resName
				// append ops path to the set
				pvd.ResourceList.Resources[resId].OperationReqPath[opsPath] = opsPath
				pvd.ResourceList.Resources[resId].operations[strings.ToLower(opName)] = opName
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
