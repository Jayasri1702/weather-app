package auth

import (
	"bufio"
	"bytes"
	"errors"
	"testing"
	"weatherapp/internal/storage"
	"weatherapp/models"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// hash: Hashes the given password using bcrypt and returns the hashed password.
func hash(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

// TestLogin: Tests the login functionality for different scenarios like successful login, invalid password, and user not found.
func TestLogin(t *testing.T) {
	originalLoadUsers := storage.LoadUsers
	defer func() { storage.LoadUsers = originalLoadUsers }()

	// Mocking the LoadUsers function for testing.
	storage.LoadUsers = func() []models.User {
		return []models.User{
			{UserID: "1", Name: "deepak", Password: hash("123")},
		}
	}

	// Test case: Successful login.
	t.Run("Successful login", func(t *testing.T) {
		input := "deepak\n123\n"
		reader := bufio.NewReader(bytes.NewBufferString(input))
		userID := Login(reader)
		assert.Equal(t, "1", userID)
	})

	// Test case: Invalid password.
	t.Run("Invalid password", func(t *testing.T) {
		input := "deepak\nwrong\n"
		reader := bufio.NewReader(bytes.NewBufferString(input))
		userID := Login(reader)
		assert.Empty(t, userID)
	})

	// Test case: User not found.
	t.Run("User not found", func(t *testing.T) {
		input := "dev\n123\n"
		reader := bufio.NewReader(bytes.NewBufferString(input))
		userID := Login(reader)
		assert.Empty(t, userID)
	})
}

// TestRegister: Tests the registration functionality, including successful registration and error handling while saving the user.
func TestRegister(t *testing.T) {
	originalSaveUser := storage.SaveUser
	defer func() { storage.SaveUser = originalSaveUser }()

	// Test case: Successful registration.
	t.Run("Successful registration", func(t *testing.T) {
		input := "testid\nTest User\npassword123\n"
		reader := bufio.NewReader(bytes.NewBufferString(input))

		var savedUser models.User
		// Mocking SaveUser function to test the registration flow.
		storage.SaveUser = func(user models.User) error {
			savedUser = user
			return nil
		}

		Register(reader)
		assert.Equal(t, "testid", savedUser.UserID)
		assert.Equal(t, "Test User", savedUser.Name)
		assert.NotEmpty(t, savedUser.Password)
	})

	// Test case: Error saving user due to database issues.
	t.Run("SaveUser error", func(t *testing.T) {
		input := "testid2\nUser2\npass2\n"
		reader := bufio.NewReader(bytes.NewBufferString(input))

		// Mocking SaveUser to simulate a database error.
		storage.SaveUser = func(user models.User) error {
			return errors.New("db error")
		}

		Register(reader)
	})
}
