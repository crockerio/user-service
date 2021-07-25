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

	// check for duplicate username
	duplicateUser := User{}
	result := c.DB.Where("username=?", params["username"]).Find(&duplicateUser)

	if result.Error == nil && result.RowsAffected > 0 {
		return nil, errors.New("username is already in use")
	}

	// check for duplicate email
	result = c.DB.Where("email=?", params["email"]).Find(&duplicateUser)

	if result.Error == nil && result.RowsAffected > 0 {
		return nil, errors.New("email is already in use")
	}

	user.Username = params["username"]
	user.Email = params["email"]
	user.Password = params["password"]
	// user.EmailVerifiedAt = nil

	result = c.DB.Create(&user)

	if result.Error != nil {
		// TODO probably want to hide the actual error - log it instead and return a "server error"
		return nil, result.Error
	}

	return user, nil
}

func (c *userController) Read(r *http.Request, params map[string]string) (interface{}, error) {
	user := User{}

	if _, ok := params["id"]; !ok {
		return nil, errors.New("id is required")
	}

	result := c.DB.Find(&user, params["id"])

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("resource not found")
	}

	return user, nil
}

func (c *userController) Update(r *http.Request, params map[string]string) (interface{}, error) {
	user := User{}

	if _, ok := params["id"]; !ok {
		return nil, errors.New("id is required")
	}

	result := c.DB.Find(&user, params["id"])

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("resource not found")
	}

	if val, ok := params["username"]; ok {
		user.Username = val
	}

	if val, ok := params["password"]; ok {
		user.Password = val
	}

	if val, ok := params["email"]; ok {
		user.Email = val
	}

	c.DB.Save(&user)

	return user, nil
}

func (c *userController) Delete(r *http.Request, params map[string]string) (interface{}, error) {
	user := User{}

	if _, ok := params["id"]; !ok {
		return nil, errors.New("id is required")
	}

	result := c.DB.Find(&user, params["id"])

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("resource not found")
	}

	result = c.DB.Delete(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return "user deleted", nil
}
