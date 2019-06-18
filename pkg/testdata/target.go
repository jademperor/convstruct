package testdata

import (
	"time"

	"github.com/jademperor/go-tools/convstruct/testdata2"
	"github.com/jinzhu/gorm"
)

// CustomModel .
type CustomModel struct {
	gorm.Model
	Name        string
	Age         int
	Birthday    time.Time
	SubModel    testdata2.SubModel
	SubModelPtr *testdata2.SubModel
	Int         SelfInt
	M           SelfModel
}

// SelfInt . for test convstruct
type SelfInt uint8

// SelfModel . for test convstruct
type SelfModel struct {
	Name string
	Int  SelfInt
}
