package convstruct

// 解析一个go文件获取其中之类特征的struct并解析结构
import (
	"go/types"
	"log"
	"strings"

	"golang.org/x/tools/go/loader"
)

// filed of struct
type field struct {
	// field name
	name string
	// field type, TODO:save as types.Type
	typ string
	// tag you got by yourself, call `SetCustomParseTagFunc` to set.
	// while origin tag is `json:"json" xml:"xml_tag"`
	// but you only want to got json, you must get this by yourself
	tag string
}

type innerStruct struct {
	// all fields
	fields []*field
	// struct origin string define
	content string
	// struct name
	name string
	// struct owned by which package
	pkgName string
}

// Exported, and specified type
func loadGoFiles(dir string, importPath string, filenames ...string) ([]*innerStruct, error) {
	var conf loader.Config

	conf.Cwd = dir
	conf.CreateFromFilenames(importPath, filenames...)

	prog, err := conf.Load()
	if err != nil {
		log.Println("load program err:", err)
		return nil, err
	}

	if isDebug {
		log.Println("dir:", dir)
		log.Println("importPath:", importPath)
		log.Println("filename:", filenames)
		// log.Println("len of prog imported", len(prog.Imported))
	}

	return loopProgramCreated(prog.Created), nil
}

// loopProgramCreated to loo and filter:
// 1. unexported type
// 2. bultin types
// 3. only specified style struct name
func loopProgramCreated(
	created []*loader.PackageInfo,
) (innerStructs []*innerStruct) {
	for _, pkgInfo := range created {
		pkgName := pkgInfo.Pkg.Name()
		defs := pkgInfo.Defs

		// imports := pkgInfo.Pkg.Imports()
		// for _, imp := range imports {
		// 	log.Println(imp.Path(), imp.Name())
		// }

		for indent, obj := range defs {
			if !indent.IsExported() ||
				obj == nil ||
				!strings.HasSuffix(indent.Name, specifiedStructTypeSuffix) {
				continue
			}

			st, ok := obj.Type().Underlying().(*types.Struct)
			if !ok {
				log.Println("not a struct, skip this")
				continue
			}
			is := new(innerStruct)

			is.content = st.String()
			is.pkgName = pkgName
			is.name = obj.Name()
			is.fields = parseStructFields(st)

			if isDebug {
				log.Println("parse one Model: ", is.name, is.pkgName, is.content)
			}

			innerStructs = append(innerStructs, is)
		}
	}
	return
}

// parseStructFields parse fields
func parseStructFields(st *types.Struct) []*field {
	flds := make([]*field, 0, st.NumFields())

	for i := 0; i < st.NumFields(); i++ {
		fld := st.Field(i)
		// skip unexported field
		if !fld.Exported() {
			continue
		}
		isField := new(field)

		isField.name = fld.Name()
		isField.tag = parseTag(st.Tag(i))
		isField.typ = fld.Type().String()

		flds = append(flds, isField)
	}
	return flds
}

// input: gorm:"colunm:name"
// output: gorm:column:name
func defaultParseTagFunc(s string) string {
	s = strings.Replace(s, `"`, "", -1)
	splited := strings.Split(s, ":")
	return splited[len(splited)-1]
}
