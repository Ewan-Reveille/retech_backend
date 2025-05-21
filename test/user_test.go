// test/models/user_test.go
package test

import (
	"testing"

	"github.com/Ewan-Reveille/retech/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserModel_Create(t *testing.T) {
	db := SetupTestDB(t)
	defer TeardownTestDB(db)

	userModel := models.UserModel{DB: db}

	user := models.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	err := userModel.Create(&user)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
}

// Add tests for GetByID, GetByUsername, etc.
