package main

import (
	"time"

	"github.com/crockerio/cservice"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username        string
	Password        string
	Email           string
	EmailVerifiedAt time.Time
	Roles           []Role `gorm:"many2many:user_roles;"`
	PasswordResets  []PasswordResets
}

type PasswordResets struct {
	UserId     int
	User       User
	Token      string
	CreatedAt  time.Time
	ValidUntil time.Time
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

	cservice.MigrateModels(&User{}, &Role{}, &PasswordResets{})
}

func main() {
	//
}
