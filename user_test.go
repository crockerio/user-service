package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crockerio/cservice"
	"github.com/crockerio/cservice/test"
)

func initUserController(t *testing.T) *userController {
	config := &cservice.DatabaseConfig{
		Driver: "sqlite",
		File:   "file::memory:?cache=shared",
		Models: []interface{}{User{}, Role{}, PasswordResets{}},
	}

	err := cservice.InitDatabase(config)
	if err != nil {
		panic(err)
	}

	controller := &userController{}
	controller.SetDB(cservice.GetDB())

	t.Cleanup(func() {
		controller.DB.Exec("DELETE FROM users;")
	})

	return controller
}

func TestUserControllerIndex(t *testing.T) {
	controller := initUserController(t)

	// Set up Data
	expectedLength := 2

	user1 := User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@example.com",
	}

	user2 := User{
		Username: "user2",
		Password: "password2",
		Email:    "user2@example.com",
	}

	controller.DB.Create(&user1)
	controller.DB.Create(&user2)

	// Make the Request
	request := httptest.NewRequest(http.MethodGet, "/user", nil)
	response, err := controller.Index(request, make(map[string]string))

	if err != nil {
		t.Error(err)
	}

	if users, ok := response.([]User); ok {
		if len(users) != expectedLength {
			t.Errorf("expected %d users, got %d", expectedLength, len(users))
		}

		result1 := users[0]
		result2 := users[1]

		// Assert User 1
		test.AssertStringEquals(t, user1.Username, result1.Username)
		test.AssertStringEquals(t, user1.Password, result1.Password)
		test.AssertStringEquals(t, user1.Email, result1.Email)

		// Assert User 2
		test.AssertStringEquals(t, user2.Username, result2.Username)
		test.AssertStringEquals(t, user2.Password, result2.Password)
		test.AssertStringEquals(t, user2.Email, result2.Email)

		// Pass the test
		return
	}

	t.Errorf("response was not of of type []User")
}

func TestUserCreate(t *testing.T) {
	controller := initUserController(t)

	user := User{
		Username: "test1",
		Password: "password1",
		Email:    "test1@example.com",
	}

	payload := make(map[string]string)
	payload["username"] = user.Username
	payload["password"] = user.Password
	payload["email"] = user.Email

	// Make the Request
	request := httptest.NewRequest(http.MethodPost, "/user", nil)
	response, err := controller.Create(request, payload)

	if err != nil {
		t.Error(err)
	}

	if result, ok := response.(User); ok {
		// Assert User
		test.AssertStringEquals(t, user.Username, result.Username)
		test.AssertStringEquals(t, user.Password, result.Password)
		test.AssertStringEquals(t, user.Email, result.Email)

		// Pass the test
		return
	}

	t.Errorf("response was not of of type User")
}

func TestUserCreate_UsernameIsRequired(t *testing.T) {
	controller := initUserController(t)

	user := User{
		Username: "test1",
		Password: "password1",
		Email:    "test1@example.com",
	}

	payload := make(map[string]string)
	payload["password"] = user.Password
	payload["email"] = user.Email

	// Make the Request
	request := httptest.NewRequest(http.MethodPost, "/user", nil)
	_, err := controller.Create(request, payload)

	test.AssertErrorThrown(t, err, "username is required")
}

func TestUserCreate_PasswordIsRequired(t *testing.T) {
	controller := initUserController(t)

	user := User{
		Username: "test1",
		Password: "password1",
		Email:    "test1@example.com",
	}

	payload := make(map[string]string)
	payload["username"] = user.Username
	payload["email"] = user.Email

	// Make the Request
	request := httptest.NewRequest(http.MethodPost, "/user", nil)
	_, err := controller.Create(request, payload)

	test.AssertErrorThrown(t, err, "password is required")
}

func TestUserCreate_EmailIsRequired(t *testing.T) {
	controller := initUserController(t)

	user := User{
		Username: "test1",
		Password: "password1",
		Email:    "test1@example.com",
	}

	payload := make(map[string]string)
	payload["username"] = user.Username
	payload["password"] = user.Password

	// Make the Request
	request := httptest.NewRequest(http.MethodPost, "/user", nil)
	_, err := controller.Create(request, payload)

	test.AssertErrorThrown(t, err, "email is required")
}

func TestUserCreate_DuplicateUsername(t *testing.T) {
	controller := initUserController(t)

	dupe := User{
		Username: "popular-username",
		Password: "password",
		Email:    "popular-email@example.com",
	}
	controller.DB.Create(&dupe)

	user := User{
		Username: "popular-username",
		Password: "password1",
		Email:    "test1@example.com",
	}

	payload := make(map[string]string)
	payload["username"] = user.Username
	payload["password"] = user.Password
	payload["email"] = user.Email

	// Make the Request
	request := httptest.NewRequest(http.MethodPost, "/user", nil)
	_, err := controller.Create(request, payload)

	test.AssertErrorThrown(t, err, "username is already in use")
}

func TestUserCreate_DuplicateEmail(t *testing.T) {
	controller := initUserController(t)

	dupe := User{
		Username: "popular-username",
		Password: "password",
		Email:    "popular-email@example.com",
	}
	controller.DB.Create(&dupe)

	user := User{
		Username: "test1",
		Password: "password1",
		Email:    "popular-email@example.com",
	}

	payload := make(map[string]string)
	payload["username"] = user.Username
	payload["password"] = user.Password
	payload["email"] = user.Email

	// Make the Request
	request := httptest.NewRequest(http.MethodPost, "/user", nil)
	_, err := controller.Create(request, payload)

	test.AssertErrorThrown(t, err, "email is already in use")
}

func TestUserRead(t *testing.T) {
	controller := initUserController(t)

	// Set up Data
	user := User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@example.com",
	}

	controller.DB.Create(&user)

	params := make(map[string]string)
	params["id"] = fmt.Sprint(user.ID)

	// Make the Request
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d", user.ID), nil)
	response, err := controller.Read(request, params)

	if err != nil {
		t.Error(err)
	}

	if result, ok := response.(User); ok {
		test.AssertEquals(t, user.ID, result.ID)
		test.AssertStringEquals(t, user.Username, result.Username)
		test.AssertStringEquals(t, user.Password, result.Password)
		test.AssertStringEquals(t, user.Email, result.Email)

		// Pass the test
		return
	}

	t.Errorf("response was not of of type User")
}

func TestUserRead_IdIsRequired(t *testing.T) {
	controller := initUserController(t)

	// Set up Data
	user := User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@example.com",
	}

	controller.DB.Create(&user)

	params := make(map[string]string)

	// Make the Request
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d", user.ID), nil)
	_, err := controller.Read(request, params)

	test.AssertErrorThrown(t, err, "id is required")
}

func TestUserRead_NotFound(t *testing.T) {
	controller := initUserController(t)

	// Set up Data
	user := User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@example.com",
	}

	controller.DB.Create(&user)

	params := make(map[string]string)
	params["id"] = "999"

	// Make the Request
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d", user.ID), nil)
	_, err := controller.Read(request, params)

	test.AssertErrorThrown(t, err, "resource not found")
}

func TestUserUpdate(t *testing.T) {
	controller := initUserController(t)

	// Set up Data
	user := User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@example.com",
	}

	controller.DB.Create(&user)

	params := make(map[string]string)
	params["id"] = fmt.Sprint(user.ID)
	params["username"] = "newusername"
	params["password"] = "newpassword"
	params["email"] = "new.email@example.com"

	// Make the Request
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/user/%d", user.ID), nil)
	response, err := controller.Update(request, params)

	if err != nil {
		t.Error(err)
	}

	if result, ok := response.(User); ok {
		test.AssertEquals(t, user.ID, result.ID)
		test.AssertStringEquals(t, "newusername", result.Username)
		test.AssertStringEquals(t, "newpassword", result.Password)
		test.AssertStringEquals(t, "new.email@example.com", result.Email)

		// Pass the test
		return
	}

	t.Errorf("response was not of of type User")
}

func TestUserUpdate_IdIsRequired(t *testing.T) {
	controller := initUserController(t)

	// Set up Data
	user := User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@example.com",
	}

	controller.DB.Create(&user)

	params := make(map[string]string)
	params["username"] = "newusername"
	params["password"] = "newpassword"
	params["email"] = "new.email@example.com"

	// Make the Request
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/user/%d", user.ID), nil)
	_, err := controller.Update(request, params)

	test.AssertErrorThrown(t, err, "id is required")
}

func TestUserUpdate_NotFound(t *testing.T) {
	controller := initUserController(t)

	// Set up Data
	user := User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@example.com",
	}

	controller.DB.Create(&user)

	params := make(map[string]string)
	params["id"] = "999"
	params["username"] = "newusername"
	params["password"] = "newpassword"
	params["email"] = "new.email@example.com"

	// Make the Request
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/user/%d", user.ID), nil)
	_, err := controller.Update(request, params)

	test.AssertErrorThrown(t, err, "resource not found")
}

func TestUserUpdate_EmptyUsernameKeepsOldValue(t *testing.T) {
	controller := initUserController(t)

	// Set up Data
	user := User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@example.com",
	}

	controller.DB.Create(&user)

	params := make(map[string]string)
	params["id"] = fmt.Sprint(user.ID)
	params["password"] = "newpassword"
	params["email"] = "new.email@example.com"

	// Make the Request
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/user/%d", user.ID), nil)
	response, err := controller.Update(request, params)

	if err != nil {
		t.Error(err)
	}

	if result, ok := response.(User); ok {
		test.AssertEquals(t, user.ID, result.ID)
		test.AssertStringEquals(t, user.Username, result.Username)
		test.AssertStringEquals(t, "newpassword", result.Password)
		test.AssertStringEquals(t, "new.email@example.com", result.Email)

		// Pass the test
		return
	}

	t.Errorf("response was not of of type User")
}

func TestUserUpdate_EmptyPasswordKeepsOldValue(t *testing.T) {
	controller := initUserController(t)

	// Set up Data
	user := User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@example.com",
	}

	controller.DB.Create(&user)

	params := make(map[string]string)
	params["id"] = fmt.Sprint(user.ID)
	params["username"] = "newusername"
	params["email"] = "new.email@example.com"

	// Make the Request
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/user/%d", user.ID), nil)
	response, err := controller.Update(request, params)

	if err != nil {
		t.Error(err)
	}

	if result, ok := response.(User); ok {
		test.AssertEquals(t, user.ID, result.ID)
		test.AssertStringEquals(t, "newusername", result.Username)
		test.AssertStringEquals(t, user.Password, result.Password)
		test.AssertStringEquals(t, "new.email@example.com", result.Email)

		// Pass the test
		return
	}

	t.Errorf("response was not of of type User")
}

func TestUserUpdate_EmptyEmailKeepsOldValue(t *testing.T) {
	controller := initUserController(t)

	// Set up Data
	user := User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@example.com",
	}

	controller.DB.Create(&user)

	params := make(map[string]string)
	params["id"] = fmt.Sprint(user.ID)
	params["username"] = "newusername"
	params["password"] = "newpassword"

	// Make the Request
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/user/%d", user.ID), nil)
	response, err := controller.Update(request, params)

	if err != nil {
		t.Error(err)
	}

	if result, ok := response.(User); ok {
		test.AssertEquals(t, user.ID, result.ID)
		test.AssertStringEquals(t, "newusername", result.Username)
		test.AssertStringEquals(t, "newpassword", result.Password)
		test.AssertStringEquals(t, user.Email, result.Email)

		// Pass the test
		return
	}

	t.Errorf("response was not of of type User")
}

func TestUserUpdate_DuplicateUsername(t *testing.T) {
	// 	controller := initUserController(t)

	// 	// Set up Data
	// 	user := User{
	// 		Username: "user1",
	// 		Password: "password1",
	// 		Email:    "user1@example.com",
	// 	}

	// 	controller.DB.Create(&user)

	// 	params := make(map[string]string)
	// 	params["id"] = fmt.Sprint(user.ID)
	// 	params["username"] = "newusername"
	// 	params["password"] = "newpassword"

	// 	// Make the Request
	// 	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/user/%d", user.ID), nil)
	// 	response, err := controller.Update(request, params)

	// 	if err != nil {
	// 		t.Error(err)
	// 	}

	// 	if result, ok := response.(User); ok {
	// 		test.AssertEquals(t, user.ID, result.ID)
	// 		test.AssertStringEquals(t, "newusdername", result.Username)
	// 		test.AssertStringEquals(t, "newpassword", result.Password)
	// 		test.AssertStringEquals(t, user.Email, result.Email)

	// 		// Pass the test
	// 		return
	// 	}

	t.Errorf("NYI; response was not of of type User")
}

func TestUserUpdate_DuplicateEmail(t *testing.T) {
	// 	controller := initUserController(t)

	// 	// Set up Data
	// 	user := User{
	// 		Username: "user1",
	// 		Password: "password1",
	// 		Email:    "user1@example.com",
	// 	}

	// 	controller.DB.Create(&user)

	// 	params := make(map[string]string)
	// 	params["id"] = fmt.Sprint(user.ID)
	// 	params["username"] = "newusername"
	// 	params["password"] = "newpassword"

	// 	// Make the Request
	// 	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/user/%d", user.ID), nil)
	// 	response, err := controller.Update(request, params)

	// 	if err != nil {
	// 		t.Error(err)
	// 	}

	// 	if result, ok := response.(User); ok {
	// 		test.AssertEquals(t, user.ID, result.ID)
	// 		test.AssertStringEquals(t, "newusername", result.Username)
	// 		test.AssertStringEquals(t, "newpdassword", result.Password)
	// 		test.AssertStringEquals(t, user.Email, result.Email)

	// 		// Pass the test
	// 		return
	// 	}

	t.Errorf("NYI; response was not of of type User")
}

func TestUserDelete(t *testing.T) {
	controller := initUserController(t)

	// Set up Data
	user := User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@example.com",
	}

	controller.DB.Create(&user)

	params := make(map[string]string)
	params["id"] = fmt.Sprint(user.ID)

	// Make the Request
	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/user/%d", user.ID), nil)
	response, err := controller.Delete(request, params)

	if err != nil {
		t.Error(err)
	}

	if result, ok := response.(string); ok {
		test.AssertEquals(t, "user deleted", result)

		// Pass the test
		return
	}

	t.Errorf("response was not of of type User")
}

func TestUserDelete_IdIsRequired(t *testing.T) {
	controller := initUserController(t)

	// Set up Data
	user := User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@example.com",
	}

	controller.DB.Create(&user)

	params := make(map[string]string)

	// Make the Request
	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/user/%d", user.ID), nil)
	_, err := controller.Delete(request, params)

	test.AssertErrorThrown(t, err, "id is required")
}

func TestUserDelete_NotFound(t *testing.T) {
	controller := initUserController(t)

	// Set up Data
	user := User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@example.com",
	}

	controller.DB.Create(&user)

	params := make(map[string]string)
	params["id"] = "999"

	// Make the Request
	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/user/%d", user.ID), nil)
	_, err := controller.Delete(request, params)

	test.AssertErrorThrown(t, err, "resource not found")
}
