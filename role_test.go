package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crockerio/cservice"
	"github.com/crockerio/cservice/test"
)

func initRoleController(t *testing.T) *roleController {
	config := &cservice.DatabaseConfig{
		Driver: "sqlite",
		File:   "file::memory:?cache=shared",
		Models: []interface{}{Role{}, Role{}, PasswordResets{}},
	}

	err := cservice.InitDatabase(config)
	if err != nil {
		panic(err)
	}

	controller := &roleController{}
	controller.SetDB(cservice.GetDB())

	t.Cleanup(func() {
		controller.DB.Exec("DELETE FROM roles;")
	})

	return controller
}

func TestRoleControllerIndex(t *testing.T) {
	controller := initRoleController(t)

	// Set up Data
	expectedLength := 2

	role1 := Role{
		Name: "Admin",
	}

	role2 := Role{
		Name: "role",
	}

	controller.DB.Create(&role1)
	controller.DB.Create(&role2)

	// Make the Request
	request := httptest.NewRequest(http.MethodGet, "/role", nil)
	response, err := controller.Index(request, make(map[string]string))

	if err != nil {
		t.Error(err)
	}

	if roles, ok := response.([]Role); ok {
		if len(roles) != expectedLength {
			t.Errorf("expected %d roles, got %d", expectedLength, len(roles))
		}

		result1 := roles[0]
		result2 := roles[1]

		// Assert role 1
		test.AssertStringEquals(t, role1.Name, result1.Name)

		// Assert role 2
		test.AssertStringEquals(t, role2.Name, result2.Name)

		// Pass the test
		return
	}

	t.Errorf("response was not of of type []Role")
}

func TestRoleCreate(t *testing.T) {
	controller := initRoleController(t)

	role := Role{
		Name: "Test Role",
	}

	payload := make(map[string]string)
	payload["name"] = role.Name

	// Make the Request
	request := httptest.NewRequest(http.MethodPost, "/role", nil)
	response, err := controller.Create(request, payload)

	if err != nil {
		t.Error(err)
	}

	if result, ok := response.(Role); ok {
		// Assert role
		test.AssertStringEquals(t, role.Name, result.Name)

		// Pass the test
		return
	}

	t.Errorf("response was not of type Role")
}

func TestRoleCreate_NameIsRequired(t *testing.T) {
	controller := initRoleController(t)

	payload := make(map[string]string)

	// Make the Request
	request := httptest.NewRequest(http.MethodPost, "/role", nil)
	_, err := controller.Create(request, payload)

	test.AssertErrorThrown(t, err, "name is required")
}

func TestRoleCreate_DuplicateName(t *testing.T) {
	controller := initRoleController(t)

	dupe := Role{
		Name: "dupe-role",
	}
	controller.DB.Create(&dupe)

	role := Role{
		Name: "dupe-role",
	}

	payload := make(map[string]string)
	payload["name"] = role.Name

	// Make the Request
	request := httptest.NewRequest(http.MethodPost, "/role", nil)
	_, err := controller.Create(request, payload)

	test.AssertErrorThrown(t, err, "name is already in use")
}

func TestRoleRead(t *testing.T) {
	controller := initRoleController(t)

	// Set up Data
	role := Role{
		Name: "Test Role",
	}

	controller.DB.Create(&role)

	params := make(map[string]string)
	params["id"] = fmt.Sprint(role.ID)

	// Make the Request
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/role/%d", role.ID), nil)
	response, err := controller.Read(request, params)

	if err != nil {
		t.Error(err)
	}

	if result, ok := response.(Role); ok {
		test.AssertEquals(t, role.ID, result.ID)
		test.AssertStringEquals(t, role.Name, result.Name)

		// Pass the test
		return
	}

	t.Errorf("response was not of type Role")
}

func TestRoleRead_IdIsRequired(t *testing.T) {
	controller := initRoleController(t)

	// Set up Data
	role := Role{
		Name: "Test Role",
	}

	controller.DB.Create(&role)

	params := make(map[string]string)

	// Make the Request
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/role/%d", role.ID), nil)
	_, err := controller.Read(request, params)

	test.AssertErrorThrown(t, err, "id is required")
}

func TestRoleRead_NotFound(t *testing.T) {
	controller := initRoleController(t)

	// Set up Data
	role := Role{
		Name: "Test Role",
	}

	controller.DB.Create(&role)

	params := make(map[string]string)
	params["id"] = "999"

	// Make the Request
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/role/%d", role.ID), nil)
	_, err := controller.Read(request, params)

	test.AssertErrorThrown(t, err, "resource not found")
}

func TestRoleUpdate(t *testing.T) {
	controller := initRoleController(t)

	// Set up Data
	role := Role{
		Name: "Test Role",
	}

	controller.DB.Create(&role)

	params := make(map[string]string)
	params["id"] = fmt.Sprint(role.ID)
	params["name"] = "Updated Role"

	// Make the Request
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/role/%d", role.ID), nil)
	response, err := controller.Update(request, params)

	if err != nil {
		t.Error(err)
	}

	if result, ok := response.(Role); ok {
		test.AssertEquals(t, role.ID, result.ID)
		test.AssertStringEquals(t, "Updated Role", result.Name)

		// Pass the test
		return
	}

	t.Errorf("response was not of type Role")
}

func TestRoleUpdate_IdIsRequired(t *testing.T) {
	controller := initRoleController(t)

	// Set up Data
	role := Role{
		Name: "Test Role",
	}

	controller.DB.Create(&role)

	params := make(map[string]string)
	params["rolename"] = "newrolename"

	// Make the Request
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/role/%d", role.ID), nil)
	_, err := controller.Update(request, params)

	test.AssertErrorThrown(t, err, "id is required")
}

func TestRoleUpdate_NotFound(t *testing.T) {
	controller := initRoleController(t)

	// Set up Data
	role := Role{
		Name: "Test Role",
	}

	controller.DB.Create(&role)

	params := make(map[string]string)
	params["id"] = "999"
	params["rolename"] = "newrolename"

	// Make the Request
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/role/%d", role.ID), nil)
	_, err := controller.Update(request, params)

	test.AssertErrorThrown(t, err, "resource not found")
}

func TestRoleUpdate_EmptyNameKeepsOldValue(t *testing.T) {
	controller := initRoleController(t)

	// Set up Data
	role := Role{
		Name: "Test Role",
	}

	controller.DB.Create(&role)

	params := make(map[string]string)
	params["id"] = fmt.Sprint(role.ID)

	// Make the Request
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/role/%d", role.ID), nil)
	response, err := controller.Update(request, params)

	if err != nil {
		t.Error(err)
	}

	if result, ok := response.(Role); ok {
		test.AssertEquals(t, role.ID, result.ID)
		test.AssertStringEquals(t, role.Name, result.Name)

		// Pass the test
		return
	}

	t.Errorf("response was not of type Role")
}

func TestRoleDelete(t *testing.T) {
	controller := initRoleController(t)

	// Set up Data
	role := Role{
		Name: "Test Role",
	}

	controller.DB.Create(&role)

	params := make(map[string]string)
	params["id"] = fmt.Sprint(role.ID)

	// Make the Request
	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/role/%d", role.ID), nil)
	response, err := controller.Delete(request, params)

	if err != nil {
		t.Error(err)
	}

	if result, ok := response.(string); ok {
		test.AssertEquals(t, "role deleted", result)

		// Pass the test
		return
	}

	t.Errorf("response was not of type Role")
}

func TestRoleDelete_IdIsRequired(t *testing.T) {
	controller := initRoleController(t)

	// Set up Data
	role := Role{
		Name: "Test Role",
	}

	controller.DB.Create(&role)

	params := make(map[string]string)

	// Make the Request
	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/role/%d", role.ID), nil)
	_, err := controller.Delete(request, params)

	test.AssertErrorThrown(t, err, "id is required")
}

func TestRoleDelete_NotFound(t *testing.T) {
	controller := initRoleController(t)

	// Set up Data
	role := Role{
		Name: "Test Role",
	}

	controller.DB.Create(&role)

	params := make(map[string]string)
	params["id"] = "999"

	// Make the Request
	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/role/%d", role.ID), nil)
	_, err := controller.Delete(request, params)

	test.AssertErrorThrown(t, err, "resource not found")
}
