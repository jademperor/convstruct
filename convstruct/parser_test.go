package convstruct

import (
	"testing"
)

func Test_loadGoFile(t *testing.T) {
	type args struct {
		dir      string
		filename string
		path     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "case 1",
			args: args{
				dir:      "/Users/yeqiang/go/src/github.com/yeqown/server-common/dbs/convstruct/testdata",
				path:     "github.com/yeqown/server-common/dbs/convstruct/testdata",
				filename: "type_model.go",
			},
			wantErr: false,
		},
	}

	isDebug = true

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ises, err := loadGoFiles(tt.args.dir, tt.args.path, tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("loadGoFile() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				for _, is := range ises {
					t.Log(is.name, is.pkgName, is.content)
					for _, fld := range is.fields {
						t.Log(*fld)
					}
				}
			}
		})
	}
}
