# convstruct
a tool to generate a struct and a convert method from a struct.

## Todos:
* [x] parse struct and source pkg
* [x] generate struct and load `func`
* [x] add imports automatic
* [x] support source pkg structs
* [ ] auto import source pkg


## Cli commands

```sh
➜  go-tools git:(master) ✗ ./bin/convstruct -h
Usage of ./bin/convstruct:
  -debug
        open debug mode
  -in string
        Filename to be parsed
  -out string
        Filename to keep result
  -outPkgName string
        Pkg name to output, only will be used when output file is not exist
  -structName string
        StructName to be parsed
```

## source
```go
package testdata

import (
	"time"

	"github.com/jinzhu/gorm"
)

// CustomModel .
type CustomModel struct {
	Model       gorm.Model
	Name        string
	Age         int
	Birthday    time.Time
	SubModel    SubModel
	SubModelPtr *SubModel
}

// SubModel .
type SubModel struct {
	SubName string
	SubAge  int
}
```

### dst[generated]

```go
package pkgname

import (
	"time"

	"path/to/testdata"
	"github.com/jinzhu/gorm"
)

// CustomModel .
type CustomModel struct {
	Model       gorm.Model         `json:"model"`
	Name        string             `json:"name"`
	Age         int                `json:"age"`
	Birthday    time.Time          `json:"birthday"`
	SubModel    testdata.SubModel  `json:"sub_model"`
	SubModelPtr *testdata.SubModel `json:"sub_model_ptr"`
}

// LoadCustomModel .
func LoadCustomModel(m *testdata.CustomModel) *CustomModel {
	return &CustomModel{
		Model:       m.Model,
		Name:        m.Name,
		Age:         m.Age,
		Birthday:    m.Birthday,
		SubModel:    m.SubModel,
		SubModelPtr: m.SubModelPtr,
	}
}
```