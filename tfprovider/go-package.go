package tfprovider

import (
	"fmt"
	"path/filepath"
)

func (pkg *GoPackage) BaseName() string {
	return filepath.Base(pkg.Package)
}

func (pkg *GoPackage) ReferenceName() string {
	if pkg.Alias != "" {
		return pkg.Alias
	}
	return pkg.BaseName()
}

func ToReferenceMap(pkgs []*GoPackage) (map[string]*GoPackage, error) {
	refs := make(map[string]*GoPackage)
	for _, pkg := range pkgs {
		name := pkg.ReferenceName()
		if _, ok := refs[name]; ok {
			return nil, fmt.Errorf("Duplicated import reference name %q", name)
		}
		refs[name] = pkg
	}
	return refs, nil
}
