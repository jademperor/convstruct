package testdata

import (
	"errors"
	"time"
)

type UserModel struct {
	Name            string    `gorm:"colunm:name"`
	Password        string    `gorm:"column:password"`
	CreateTime      time.Time `gorm:"colunm:create_time"`
	UpdateTime      time.Time `gorm:"colunm:update_time"`
	unexported      int       `gorm:"column:unexported"`
	AnotherNewField string    `gorm:"column:another_new_field"`
}

type userStruct struct {
	Name     string `json:"typeStructName"`
	Password string `json:"typeStructPassword"`
}

type AliasString string

var (
	A UserModel
	B int64
	C string
)

func f(a string) error {
	return errors.New(a)
}
