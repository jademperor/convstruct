package pkg_test

import (
	"testing"

	"github.com/jademperor/go-tools/pkg"
	"github.com/yeqown/infrastructure/pkg/fs"
)

func Test_ParsePkg(t *testing.T) {
	filenames := fs.ListFiles("./testdata", fs.IgnoreDirFilter())
	// t.Logf("filenames to ParsePkg: %v", filenames)

	pkg, err := pkg.ParsePkg("./testdata", filenames)
	_ = pkg
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for _, v := range pkg.Imports {
		t.Log(v.Path, v.Name)
	}
}
