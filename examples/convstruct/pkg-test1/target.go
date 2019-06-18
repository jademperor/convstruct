package test1

import (
	"time"

	test2 "github.com/jademperor/go-tools/examples/convstruct/pkg-test2"
	"github.com/jinzhu/gorm"
)

// CustomModel .
type CustomModel struct {
	Model       gorm.Model
	Name        string
	Age         int
	Birthday    time.Time
	SubModel    test2.SubModel
	SubModelPtr *test2.SubModel
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
