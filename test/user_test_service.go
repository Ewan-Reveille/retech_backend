// test/services/user_service_test.go
package test

import (
	"testing"

	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/services"
	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(db)

	userModel := models.UserModel{DB: db}
	userService := services.UserService{Repo: &userModel}

	req := services.CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "securepassword",
	}

	user, err := userService.Create(&req)
	assert.NoError(t, err)
	assert.Equal(t, "newuser", user.Username)
	assert.Empty(t, user.Password) // Password should be hidden
}

// Test duplicate user creation, login, etc.