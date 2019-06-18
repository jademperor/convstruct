package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/jademperor/go-tools/pkg"
	"github.com/yeqown/infrastructure/pkg/fs"
)

/*
 1. 使用模版来生成
 2. Append到文件中
 2.1 如果不存在文件则创建且自动import
 2.2 如果存在则注释原有的结构体，并自动更新import
 3. 对于任意类型来说，如果该类型没有selector 那么则默认是被解析的模型的类型
*/

func main() {
	var (
		flagIn         = flag.String("in", "", "Filename to be parsed")
		flagStructName = flag.String("structName", "", "StructName to be parsed")
		flagOut        = flag.String("out", "", "Filename to keep result")
		flagPkgName    = flag.String("outPkgName", "", "Pkg name to output, only will be used when output file is not exist")
	)

	flag.Parse()

	c := &config{
		inputPkgDir:      *flagIn,
		outputFile:       *flagOut,
		targetStructName: *flagStructName,
		outPkgName:       *flagPkgName,

		imports: make(map[string]*pkg.GenImport),
	}

	if err := c.parse(); err != nil {
		panic(err)
	}
	defer c.w.Close()

	genStruct, err := c.process()
	if err != nil {
		panic(err)
	}

	if err := c.generate(genStruct); err != nil {
		panic(err)
	}
}

// config .
type config struct {
	inputPkgDir      string
	outputFile       string
	outPkgName       string
	targetStructName string

	flagOutputFileBeGenerated bool
	imports                   map[string]*pkg.GenImport
	inPkg                     *pkg.Package
	outPkg                    *pkg.Package
	w                         io.WriteCloser
}

func (c *config) valid() bool {
	return true
}

func (c *config) parse() (err error) {
	if !c.valid() {
		return errors.New("invalid config")
	}

	filenames := fs.ListFiles("./testdata", fs.IgnoreDirFilter())
	if c.inPkg, err = pkg.ParsePkg(c.inputPkgDir, filenames); err != nil {
		return err
	}

	if _, err := os.Stat(c.outputFile); os.IsNotExist(err) {
		// generate file
		if c.w, err = os.Create(c.outputFile); err != nil {
			return err
		}
		c.flagOutputFileBeGenerated = true
	} else {
		if c.outPkg, err = pkg.ParseFile(c.outputFile); err != nil {
			return err
		}
		if c.w, err = os.Open(c.outputFile); err != nil {
			return err
		}
		c.flagOutputFileBeGenerated = false
	}

	return nil
}

func (c *config) process() (*genConvStruct, error) {
	tarStruct := c.selectStruct()
	if tarStruct == nil {
		return nil, errors.New("Unable to find the struct: " + c.targetStructName)
	}

	genStruct := &genConvStruct{
		GenStruct: &pkg.GenStruct{
			Name:   c.targetStructName,
			Doc:    "this is Doc for test",
			Fields: make([]*pkg.GenField, 0),
		},
		FromPkg: &pkg.GenImport{
			Name: c.inPkg.Name,
			Path: c.inPkg.Path,
		},
	}

	for _, field := range tarStruct.Fields.List {
		genStruct.Fields = append(genStruct.Fields, c.processField(field))
	}

	return genStruct, nil
}

func (c *config) selectStruct() *ast.StructType {
	// var encStruct *ast.StructType
	for _, typ := range c.inPkg.Types {
		if typ.Name == c.targetStructName && typ.IsStruct {
			// encStruct =
			return typ.Expr.(*ast.StructType)
		}
	}

	return nil
}

func (c *config) processField(f *ast.Field) *pkg.GenField {
	genField := &pkg.GenField{}

	fieldName := ""
	if len(f.Names) != 0 {
		fieldName = f.Names[0].Name
	}

	// anonymous field
	if f.Names == nil {
		ident, ok := f.Type.(*ast.Ident)
		if !ok {
			pkg.DebugF("[DEBUG] could not convert Type(%v) to ast.Ident", f.Type)
		}

		fieldName = ident.Name
	}
	genField.Name = fieldName
	genField.Expr = c.processExpr(f.Type)
	genField.Tag = quote(processFieldTag(fieldName))

	return genField
}

func (c *config) generate(genStruct *genConvStruct) error {
	if c.flagOutputFileBeGenerated {
		// true: 新的文件
		tmpl := template.Must(
			template.New("header").Parse(TmplHeader))

		type genPkgHeader struct {
			PkgName string
		}

		if err := tmpl.Execute(c.w, genPkgHeader{PkgName: c.outPkgName}); err != nil {
			log.Fatal(err)
		}
	}

	// TO: add import
	type importsWrap struct {
		Imports []*pkg.GenImport
	}
	var (
		_imports []*pkg.GenImport
	)

	for _, imp := range c.imports {
		_imports = append(_imports, imp)
	}

	pkg.DebugF("will add imports: %v", _imports)
	tmpl := template.Must(template.New("imports").Parse(TmplImport))
	tmpl.Execute(c.w, &importsWrap{_imports})

	// TO: add struct
	tmpl = template.Must(template.New("struct").Parse(TmplStruct))
	tmpl.Execute(c.w, genStruct)

	// TO: add conv func
	tmpl = template.Must(template.New("method").Parse(TmplMethod))
	tmpl.Execute(c.w, genStruct)

	// format.Source()

	return nil
}

func (c *config) processExpr(expr ast.Expr) string {
	var (
		exprStr string
	)

	switch t := expr.(type) {
	case *ast.SelectorExpr:
		_pkgName := c.processExpr(t.X)
		exprStr = _pkgName + "." + t.Sel.Name
		c.addImport(_pkgName) // auto-add import
	case *ast.StarExpr:
		pkg.DebugF("get starExpr: %v", t)
		exprStr = c.processExpr(t.X)
	case *ast.Ident:
		pkg.DebugF("get ast.Ident: %v", t)
		exprStr = t.Name
		// TODO: add self Import
		// c.addImport("")
	default:
		pkg.DebugF("get expr: %v", t)
	}

	return exprStr
}

func (c *config) addImport(pkgName string) {
	pkg.DebugF("[addImport] want find pkg: %s", pkgName)

	for _, imp := range c.inPkg.Imports {
		pkg.DebugF("[addImport] compare pkgName: %s with %s", imp.Name, pkgName)
		if imp.Name == pkgName {
			c.imports[pkgName] = &pkg.GenImport{
				Name: pkgName,
				Path: imp.Path,
			}
			break
		}
	}

}

func processFieldTag(fieldName string) string {
	return fmt.Sprintf(`json:"%s,omitempty"`, strcase.ToSnake(fieldName))
}

func quote(tag string) string {
	return "`" + tag + "`"
}
