package main

import (
	"time"

	"github.com/crockerio/cservice"
)

type PasswordResets struct {
	UserId     int
	User       User
	Token      string
	CreatedAt  time.Time
	ValidUntil time.Time
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

	server := cservice.NewServer(5000)
	server.Resource("/user", &userController{})
	server.Resource("/role", &roleController{})

	server.Start()
}
