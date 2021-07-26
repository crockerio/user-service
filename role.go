package main

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name string
}

type roleController struct {
	DB *gorm.DB
}

func (c *roleController) SetDB(db *gorm.DB) {
	c.DB = db
}

func (c *roleController) Index(r *http.Request, params map[string]string) (interface{}, error) {
	roles := []Role{}
	result := c.DB.Find(&roles)

	if result.Error != nil {
		return nil, result.Error
	}

	return roles, nil
}

func (c *roleController) Create(r *http.Request, params map[string]string) (interface{}, error) {
	role := Role{}

	if _, ok := params["name"]; !ok {
		return nil, errors.New("name is required")
	}

	// check for duplicate name
	duplicateRole := Role{}
	result := c.DB.Where("name=?", params["name"]).Find(&duplicateRole)

	if result.Error == nil && result.RowsAffected > 0 {
		return nil, errors.New("name is already in use")
	}

	role.Name = params["name"]

	result = c.DB.Create(&role)

	if result.Error != nil {
		// TODO probably want to hide the actual error - log it instead and return a "server error"
		return nil, result.Error
	}

	return role, nil
}

func (c *roleController) Read(r *http.Request, params map[string]string) (interface{}, error) {
	role := Role{}

	if _, ok := params["id"]; !ok {
		return nil, errors.New("id is required")
	}

	result := c.DB.Find(&role, params["id"])

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("resource not found")
	}

	return role, nil
}

func (c *roleController) Update(r *http.Request, params map[string]string) (interface{}, error) {
	role := Role{}

	if _, ok := params["id"]; !ok {
		return nil, errors.New("id is required")
	}

	result := c.DB.Find(&role, params["id"])

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("resource not found")
	}

	if val, ok := params["name"]; ok {
		role.Name = val
	}

	c.DB.Save(&role)

	return role, nil
}

func (c *roleController) Delete(r *http.Request, params map[string]string) (interface{}, error) {
	role := Role{}

	if _, ok := params["id"]; !ok {
		return nil, errors.New("id is required")
	}

	result := c.DB.Find(&role, params["id"])

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("resource not found")
	}

	result = c.DB.Delete(&role)

	if result.Error != nil {
		return nil, result.Error
	}

	return "role deleted", nil
}
