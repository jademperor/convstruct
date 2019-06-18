# convstruct
a tool to generate a struct and a convert method from a struct.

## Todos:
* [ ] parse struct and source pkg
* [ ] generate struct and load `func`
* [ ] add imports automatic
* [ ] support source pkg structs


### source
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