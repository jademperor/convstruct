package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/jademperor/common/pkg/utils"
	"github.com/jademperor/go-tools/convstruct"
)

var (
	offsetFlag = flag.String("offset", "", "file and byte offset of identifier to be renamed, e.g. 'file.go:#123'.  For use by editors.")
	fromFlag   = flag.String("from", "", "identifier to be renamed; see -help for formats")
	toFlag     = flag.String("to", "", "new name for identifier")
	helpFlag   = flag.Bool("help", false, "show usage message")

	filenames utils.StringArray

	debug = flag.Bool("debug", false, "debug mode, if open this will output info")
	dir   = flag.String("dir", ".", "model directory from where")
	// generate cfg
	generateFilename     = flag.String("generateFilename", "types.go", "generate file name will be use this, default is (types.go)")
	generateDir          = flag.String("generateDir", ".", "generate file will be saved here, default is (current dir)")
	generatePkgName      = flag.String("generatePkgName", "types", "generate package name, default is (types)")
	generateStructSuffix = flag.String("generateStructSuffix", "", "replace model struct name suffix, like: (UserSuffix => User)")
	// model cfg
	modelImportPath   = flag.String("modelImportPath", "", "model package path, cannot be empty, like (my-server/models)")
	modelStructSuffix = flag.String("modelStructSuffix", "Model", "specified in which Model name style can be generate")
)

// go run tool.main.go -dir=./convstruct/testdata -filename=type_model.go -generatePkgName=testdata -modelImportPath -generateDir=./convstruct/testdata
func main() {
	flag.Var(&filenames, "filename", "specified filename those you want to generate, if no filename be set, will parse all files under ($dir)")
	flag.Parse()

	exportDir, _ := filepath.Abs(*generateDir)
	*dir, _ = filepath.Abs(*dir)
	if *debug {
		log.Println("fromDir:", *dir)
		log.Println("generateFilename:", *generateFilename)
		log.Println("exportDir:", exportDir)
	}

	// set custom funcs
	// convstruct.SetCustomGenTagFunc(CustomGenerateTagFunc)
	// convstruct.SetCustomParseTagFunc(CustomParseTagFunc)

	if len(filenames) == 0 {
		files, _ := ioutil.ReadDir(*dir)
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
				continue
			}
			filenames = append(filenames, file.Name())
		}
	}

	cfg := &convstruct.UsageCfg{
		ExportDir:          exportDir,
		ExportFilename:     *generateFilename,
		ExportPkgName:      *generatePkgName,
		ExportStructSuffix: *generateStructSuffix,
		ModelImportPath:    *modelImportPath,
		StructSuffix:       *modelStructSuffix,
		Debug:              *debug,
		Filenames:          filenames,
		Dir:                *dir,
	}

	if *debug {
		log.Println("following filenames will be parsed", filenames)
	}

	if err := convstruct.ParseAndGenerate(cfg); err != nil {
		panic(err)
	}

	println("done!")
}

// CustomParseTagFunc to custom implment yourself parseTagFunc
// @param (gorm:"colunm:name")
// return ("name")
func CustomParseTagFunc(s string) string {
	log.Println("calling  CustomParseTagFunc", s)

	s = strings.Replace(s, `"`, "", -1)
	splited := strings.Split(s, ":")
	return splited[len(splited)-1]
}

// CustomGenerateTagFunc to implment yourself generateTagFunc
// @param name fieldName (Age, Name, Year, CreateTime)
// @param typ fieldType (string, int64, time.Time)
// @param tag (CustomParseTagFunc) return value, default is gorm tag
// return (json:"name")
func CustomGenerateTagFunc(name, typ, tag string) string {
	return fmt.Sprintf("json:\"%s\"", tag)
}
