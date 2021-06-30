package main

import (
	"github.com/crockerio/cservice"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string
	Email    string
	Roles    []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
	gorm.Model
	Name string
}

func init() {
	err := cservice.InitDatabase("root:root@tcp(localhost:3306)/user-service?charset=utf8&parseTime=True&loc=Local", &gorm.Config{})
	if err != nil {
		panic(err)
	}

	cservice.MigrateModels(&User{}, &Role{})
}

func main() {
	//
}
