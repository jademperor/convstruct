package pkg

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// ParsePkg .
func ParsePkg(pkgPath string, filenames []string) (*Package, error) {
	pkgPath, _ = filepath.Abs(pkgPath)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// load all types
	allTyps := make(map[string]*TypeSpec)
	oriImports := make([]*ast.ImportSpec, 0)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			oriImports = append(oriImports, file.Imports...)

			typs := CollectTypes(file)
			for k, typ := range typs {
				allTyps[k] = typ
			}
		}
	}

	allImports := make([]*GenImport, len(oriImports))
	for idx, v := range oriImports {
		allImports[idx] = ParseImportSpec(v)
	}

	return &Package{
		Imports: allImports,
		Types:   allTyps,
		// Name: TODO:
		// Path: TODO:
	}, nil
}

// ParseFile .
func ParseFile(filename string) (*Package, error) {
	filename, _ = filepath.Abs(filename)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	obj := file.Scope.Lookup("package")
	DebugF("find package: %v", obj)

	// load all types
	allTyps := make(map[string]*TypeSpec)
	typs := CollectTypes(file)
	for k, typ := range typs {
		allTyps[k] = typ
	}
	// all imports
	oriImports := make([]*ast.ImportSpec, 0)
	oriImports = append(oriImports, file.Imports...)
	allImports := make([]*GenImport, len(oriImports))
	for idx, v := range oriImports {
		allImports[idx] = ParseImportSpec(v)
	}

	return &Package{
		Imports: allImports,
		Types:   allTyps,
		// Name: TODO:
		// Path: TODO:
	}, nil
}

// ParseImportSpec .
func ParseImportSpec(spec *ast.ImportSpec) *GenImport {
	if spec == nil {
		return nil
	}

	var (
		name string
		path = spec.Path.Value
	)
	if spec.Name == nil {
		sl := strings.Split(path, "/")
		name = sl[len(sl)-1]
		name = strings.Trim(name, `"`)
	}

	return &GenImport{
		Name: name,
		Path: path,
	}
}
