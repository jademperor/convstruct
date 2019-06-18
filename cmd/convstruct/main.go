package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"io"
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
		flagDebug      = flag.Bool("debug", false, "open debug mode")
	)

	flag.Parse()

	if *flagDebug == false {
		// true: debug == false
		pkg.CloseDebug()
	}

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

	pkg.DebugF("generated struct: %v, %v", genStruct.GenStruct, genStruct.FromPkg)
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
	w                         io.ReadWriteCloser
}

func (c *config) valid() bool {
	if c.targetStructName == "" {
		fmt.Println("empty target struct name")
		return false
	}

	if c.outputFile == "" {
		c.outputFile = "out.go"
	}

	if c.inputPkgDir == "" {
		fmt.Println("empty input package dir")
		return false
	}

	// if c.outPkgName == "" {
	// }

	return true
}

// parse .
// parse input package and output file
func (c *config) parse() (err error) {
	if !c.valid() {
		return errors.New("invalid config")
	}

	filenames := fs.ListFiles(c.inputPkgDir, fs.IgnoreDirFilter())
	if c.inPkg, err = pkg.ParsePkg(c.inputPkgDir, filenames); err != nil {
		return err
	}

	if _, err := os.Stat(c.outputFile); os.IsNotExist(err) {
		c.flagOutputFileBeGenerated = true
	} else {
		if c.outPkg, err = pkg.ParseFile(c.outputFile); err != nil {
			return err
		}
		c.flagOutputFileBeGenerated = false
	}

	c.w, err = os.OpenFile(c.outputFile,
		os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)

	return err
}

// process to
// select the target struct which is wanted to generate from
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

	// process all fields in target struct
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
	if f.Names == nil && fieldName != "" {
		ident, ok := f.Type.(*ast.Ident)
		if !ok {
			// t := f.Type.(type)
			pkg.DebugF("[processField] could not convert Type(%T) to ast.Ident", f.Type)
		}

		fieldName = ident.Name
	}
	genField.Name = fieldName
	genField.Expr = c.processExpr(f.Type)
	genField.Tag = quote(processFieldTag(fieldName))

	return genField
}

func (c *config) generate(genStruct *genConvStruct) error {
	// byts, _ := ioutil.ReadAll(c.w)
	buf := bytes.NewBuffer(nil)

	if c.flagOutputFileBeGenerated {
		// true: output file not exist
		tmpl := template.Must(
			template.New("header").Parse(TmplHeader))

		type genPkgHeaderWrap struct {
			PkgName string
		}

		if err := tmpl.Execute(buf, genPkgHeaderWrap{PkgName: c.outPkgName}); err != nil {
			return err
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

		pkg.DebugF("[generate] will add imports: %v", _imports)
		tmpl = template.Must(template.New("imports").Parse(TmplImport))
		if err := tmpl.Execute(buf, &importsWrap{_imports}); err != nil {
			return err
		}
	}

	// TO: add struct type decl
	tmpl := template.Must(template.New("struct").Parse(TmplStruct))
	if err := tmpl.Execute(buf, genStruct); err != nil {
		return err
	}

	// TO: add convert func
	tmpl = template.Must(template.New("method").Parse(TmplMethod))
	if err := tmpl.Execute(buf, genStruct); err != nil {
		return err
	}

	formated, err := format.Source(buf.Bytes())
	pkg.DebugF("formatted: %s", formated)
	_, err = c.w.Write(formated)

	return err
}

func (c *config) processExpr(expr ast.Expr) string {
	var (
		exprStr string
	)

	switch t := expr.(type) {
	case *ast.SelectorExpr:
		exprStr = c.processExpr(t.X)
		// auto-add import
		c.addImport(exprStr)
		exprStr = exprStr + "." + t.Sel.Name
	case *ast.StarExpr:
		pkg.DebugF("[processExpr] get ast.StarExpr: %v", t)
		exprStr = "*" + c.processExpr(t.X)
	case *ast.Ident:
		pkg.DebugF("[processExpr] get ast.Ident: %v", t)
		exprStr = c.selfPkgExpr(t.Name)
	default:
		pkg.DebugF("[processExpr] get expr: %v", t)
	}

	return exprStr
}

func (c *config) selfPkgExpr(typName string) string {
	if _, ok := c.inPkg.Types[typName]; ok {
		c.addImport(c.inPkg.Name)
		return c.inPkg.Name + "." + typName
	}
	return typName
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
