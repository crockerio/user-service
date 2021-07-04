package main

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string
	Password       string
	Email          string
	Roles          []Role `gorm:"many2many:user_roles;"`
	PasswordResets []PasswordResets
}

type userController struct {
	DB *gorm.DB
}

func (c *userController) SetDB(db *gorm.DB) {
	c.DB = db
}

func (c *userController) Index(r *http.Request, params map[string]string) (interface{}, error) {
	users := []User{}
	result := c.DB.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (c *userController) Create(r *http.Request, params map[string]string) (interface{}, error) {
	user := User{}

	if _, ok := params["username"]; !ok {
		return nil, errors.New("username is required")
	}

	if _, ok := params["email"]; !ok {
		return nil, errors.New("email is required")
	}

	if _, ok := params["password"]; !ok {
		return nil, errors.New("password is required")
	}

	// TODO check for duplicate username
	// TODO check for duplicate email

	user.Username = params["username"]
	user.Email = params["email"]
	user.Password = params["password"]
	// user.EmailVerifiedAt = nil

	result := c.DB.Create(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (c *userController) Read(r *http.Request, params map[string]string) (interface{}, error) {
	return "read", nil
}

func (c *userController) Update(r *http.Request, params map[string]string) (interface{}, error) {
	return "update", nil
}

func (c *userController) Delete(r *http.Request, params map[string]string) (interface{}, error) {
	return "delete", nil
}
