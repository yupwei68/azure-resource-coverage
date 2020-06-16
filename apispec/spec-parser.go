package apispec

import (
	"fmt"
	"os"
	"path/filepath"
)

func LoadFrom(rootPath string) (*ApiSpec, error) {
	specPath := filepath.Join(rootPath, "specification")
	potentialJsons := make([]string, 0)

	err := filepath.Walk(specPath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".json" && filepath.VerifyFilePath(path) {
				potentialJsons = append(potentialJsons, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	result := &ApiSpec{
		specPath,
		make(map[namespaceLocator]*NamespaceDefinition),
	}
	for _, json := range potentialJsons {
		if err := result.parseJson(json); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// jsonPath = <spec>/<namespace>/<mng|data|control>/../<provider>/<stable|preview>/<ver>/<file>.json
func (spec *ApiSpec) parseJson(jsonPath string) error {
	rel, err := filepath.Rel(spec.FullPath, jsonPath)
	if err != nil {
		return err
	}

	for p := rel; p != "." && p != ""; p = filepath.Dir(p) {
		if filepath.Base(p) == "examples" {
			return nil
		}
	}

	ancestor := filepath.Dir(rel)
	if ancestor == "." || ancestor == "" {
		return fmt.Errorf("Invalid json path %q, cannot parse SDKVersion", jsonPath)
	}
	ver := filepath.Base(ancestor)

	ancestor = filepath.Dir(ancestor)
	if ancestor == "." || ancestor == "" {
		return fmt.Errorf("Invalid json path %q, cannot parse IsPreview", jsonPath)
	}
	name := filepath.Base(ancestor)
	if name != "preview" && name != "stable" && name != "common" {
		return fmt.Errorf("Invalid json path %q, cannot parse IsPreview", jsonPath)
	}
	isPreview := name == "preview"

	ancestor = filepath.Dir(ancestor)
	if ancestor == "." || ancestor == "" {
		return fmt.Errorf("Invalid json path %q, cannot parse Provider", jsonPath)
	}
	provider := filepath.Base(ancestor)
	pvdPath := ancestor
	if provider == "data-plane" {
		ancestor = filepath.Join(ancestor, provider)
		provider = "<unknown>"
	}

	nsType := Unknown
	for nsType == Unknown && ancestor != "." && ancestor != "" {
		ancestor = filepath.Dir(ancestor)
		name = filepath.Base(ancestor)
		if name == "resource-manager" {
			nsType = Management
		} else if name == "data-plane" {
			nsType = DataPlane
		} else if name == "control-plane" {
			nsType = ControlPlane
		}
	}
	if nsType == Unknown {
		return fmt.Errorf("Invalid json path %q, cannot parse NamespaceType", jsonPath)
	}

	ancestor = filepath.Dir(ancestor)
	if ancestor == "." || ancestor == "" {
		return fmt.Errorf("Invalid json path %q, cannot parse Namespace", jsonPath)
	}
	namespace := filepath.Base(ancestor)
	nsPath := ancestor

	ancestor = filepath.Dir(ancestor)
	if ancestor != "." && ancestor != "" {
		return fmt.Errorf("Invalid json path %q, contain unparsed elements", jsonPath)
	}

	ns := spec.getOrCreateNamespace(namespace, nsPath, nsType)
	pvd := ns.getOrCreateProvider(provider, pvdPath)
	return pvd.parseResourcesFromJSON(jsonPath, VersionDefinition{isPreview, ver})
}
