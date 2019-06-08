package tfprovider

import "path/filepath"

func (pkg *GoPackage) BaseName() string {
	return filepath.Base(pkg.Package)
}

func (pkg *GoPackage) ReferenceName() string {
	if pkg.Alias != "" {
		return pkg.Alias
	}
	return pkg.BaseName()
}
