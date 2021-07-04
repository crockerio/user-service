package main

import (
	"time"

	"github.com/crockerio/cservice"

	"gorm.io/gorm"
)

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

func main() {
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

	server := cservice.NewServer()
	server.Resource("/user", &userController{})

	server.Start(5000)
}
