package test

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"net/http/httptest"
	"testing"

	"github.com/Ewan-Reveille/retech/models"
	"github.com/stretchr/testify/assert"
)

func TestUserRoutes(t *testing.T) {
	app := SetupApp()

	user := models.User{
		Username: "testuser",
		Email:    "test@testing.com",
		Password: "Password1234*",
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, res.StatusCode)

	var createdUser models.User
	json.NewDecoder(res.Body).Decode(&createdUser)
	assert.NotEmpty(t, createdUser.ID)

	// // Test login
	// loginBody, _ := json.Marshal(map[string]string{
	// 	"username": user.Username,
	// 	"password": user.Password,
	// })

	// loginReq := httptest.NewRequest("POST", "/login", bytes.NewReader(loginBody))
	// loginReq.Header.Set("Content-Type", "application/json")
	// loginRes, err := app.Test(loginReq)
	// assert.NoError(t, err)
	// fmt.Println("Login response:", loginRes)
	// assert.Equal(t, 200, loginRes.StatusCode)
	// var loggedInUser models.User
	// json.NewDecoder(loginRes.Body).Decode(&loggedInUser)
	// assert.Equal(t, createdUser.Username, loggedInUser.Username)
	// assert.Equal(t, createdUser.Email, loggedInUser.Email)
	// assert.Empty(t, loggedInUser.Password)

	// Test get all users
	getReq := httptest.NewRequest("GET", "/user", nil)
	getRes, err := app.Test(getReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, getRes.StatusCode)
	var users []models.User
	json.NewDecoder(getRes.Body).Decode(&users)
	assert.NotEmpty(t, users)
	assert.Greater(t, len(users), 0)
	assert.Equal(t, createdUser.Username, users[0].Username)
	assert.Equal(t, createdUser.Email, users[0].Email)
	assert.Empty(t, users[0].Password)
	// Test duplicate user
	// duplicateBody, _ := json.Marshal(user)
	// duplicateReq := httptest.NewRequest("POST", "/register", bytes.NewReader(duplicateBody))
	// duplicateReq.Header.Set("Content-Type", "application/json")
	// duplicateRes, err := app.Test(duplicateReq)
	// assert.NoError(t, err)
	// assert.Equal(t, 409, duplicateRes.StatusCode)
	// var duplicateResponse map[string]string
	// json.NewDecoder(duplicateRes.Body).Decode(&duplicateResponse)
	// assert.Equal(t, "username ou email déjà utilisé", duplicateResponse["error"])
	// // Test invalid login
	// invalidLoginBody, _ := json.Marshal(map[string]string{
	// 	"username": "invaliduser",
	// 	"password": "invalidpassword",
	// })

	// invalidLoginReq := httptest.NewRequest("POST", "/login", bytes.NewReader(invalidLoginBody))
	// invalidLoginReq.Header.Set("Content-Type", "application/json")
	// invalidLoginRes, err := app.Test(invalidLoginReq)
	// assert.NoError(t, err)
	// assert.Equal(t, 401, invalidLoginRes.StatusCode)
	// var invalidLoginResponse map[string]string
	// json.NewDecoder(invalidLoginRes.Body).Decode(&invalidLoginResponse)
	// assert.Equal(t, "Invalid credentials", invalidLoginResponse["error"])

}