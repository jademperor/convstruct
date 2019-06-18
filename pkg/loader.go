package pkg

import (
	"errors"
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

	if len(pkgs) > 1 {
		return nil, errors.New("package depth must be equal to 1")
	}

	// load all types
	allTyps := make(map[string]*TypeSpec)
	oriImports := make([]*ast.ImportSpec, 0)
	var pkgName string
	for k, pkg := range pkgs {
		// pkgName = pkg.Name
		pkgName = k
		DebugF("[ParsePkg] range package: %s", k)
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

	// allImports = append(allImports, &GenImport{
	// 	Name: pkgName,
	// 	Path: pkgName,
	// })

	return &Package{
		Imports: allImports,
		Types:   allTyps,
		Name:    pkgName,
		Path:    pkgName,
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

	// obj := file.Scope.Lookup("package")
	// DebugF("[ParseFile] find package: %v", obj)

	// load all types
	allTyps := make(map[string]*TypeSpec)
	typs := CollectTypes(file)
	for k, typ := range typs {
		allTyps[k] = typ
	}

	// load all imports
	oriImports := make([]*ast.ImportSpec, 0)
	oriImports = append(oriImports, file.Imports...)
	allImports := make([]*GenImport, len(oriImports))
	for idx, v := range oriImports {
		allImports[idx] = ParseImportSpec(v)
	}

	return &Package{
		Imports: allImports,
		Types:   allTyps,
		// Name: "TODO",
		// Path: "TODO",
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
	} else {
		name = spec.Name.Name
	}

	return &GenImport{
		Name: name,
		Path: path,
	}
}
