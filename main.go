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
	db := &cservice.DatabaseConfig{
		User:     "root",
		Password: "root",
		Host:     "localhost",
		Port:     3306,
		Database: "user-service",
		Models:   []interface{}{User{}, Role{}, PasswordResets{}},
	}

	err := cservice.InitDatabase(db)
	if err != nil {
		panic(err)
	}
}

func main() {
	//
}
