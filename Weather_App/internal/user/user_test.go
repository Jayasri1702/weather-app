package user

import (
	"bufio"
	"bytes"
	"testing"
	"weatherapp/internal/storage"
	"weatherapp/models"

	"github.com/stretchr/testify/assert"
)


// TestChangePreferences checks if user preferences are correctly updated via CLI input.
func TestChangePreferences(t *testing.T) {
	originalLoadUsers := storage.LoadUsers
	originalUpdateUser := storage.UpdateUser
	defer func() {
		storage.LoadUsers = originalLoadUsers
		storage.UpdateUser = originalUpdateUser
	}()

	sampleUser := models.User{UserID: "u1", Preferences: models.Preferences{}}
	storage.LoadUsers = func() []models.User {
		return []models.User{sampleUser}
	}

	var updatedUser models.User
	storage.UpdateUser = func(user models.User) error {
		updatedUser = user
		return nil
	}

	input := "Mumbai\ncelsius\nverbose\nweek\n"
	reader := bufio.NewReader(bytes.NewBufferString(input))
	ChangePreferences(reader, "u1")

	assert.Equal(t, "Mumbai", updatedUser.Preferences.Location)
	assert.Equal(t, "celsius", updatedUser.Preferences.Unit)
	assert.Equal(t, "verbose", updatedUser.Preferences.Verbosity)
	assert.Equal(t, "week", updatedUser.Preferences.Forecast)
}

// TestEnsurePreferences_PrefSet verifies EnsurePreferences does nothing if preferences already exist.
func TestEnsurePreferences_PrefSet(t *testing.T) {
	originalLoadUsers := storage.LoadUsers
	originalUpdateUser := storage.UpdateUser
	defer func() {
		storage.LoadUsers = originalLoadUsers
		storage.UpdateUser = originalUpdateUser
	}()

	existing := models.User{
		UserID: "u1",
		Preferences: models.Preferences{
			Location: "SetAlready",
		},
	}

	storage.LoadUsers = func() []models.User {
		return []models.User{existing}
	}

	storage.UpdateUser = func(user models.User) error {
		t.FailNow()
		return nil
	}

	reader := bufio.NewReader(bytes.NewBuffer(nil))
	EnsurePreferences(reader, "u1")
}

// TestEnsurePreferences_PrefNotSet checks if EnsurePreferences updates preferences when not already set.
func TestEnsurePreferences_PrefNotSet(t *testing.T) {
	originalLoadUsers := storage.LoadUsers
	originalUpdateUser := storage.UpdateUser
	defer func() {
		storage.LoadUsers = originalLoadUsers
		storage.UpdateUser = originalUpdateUser
	}()

	sampleUser := models.User{UserID: "u1", Preferences: models.Preferences{}}
	storage.LoadUsers = func() []models.User {
		return []models.User{sampleUser}
	}

	var updatedUser models.User
	storage.UpdateUser = func(user models.User) error {
		updatedUser = user
		return nil
	}

	input := "Delhi\nfahrenheit\nbrief\nday\n"
	reader := bufio.NewReader(bytes.NewBufferString(input))
	EnsurePreferences(reader, "u1")

	assert.Equal(t, "Delhi", updatedUser.Preferences.Location)
	assert.Equal(t, "fahrenheit", updatedUser.Preferences.Unit)
	assert.Equal(t, "brief", updatedUser.Preferences.Verbosity)
	assert.Equal(t, "day", updatedUser.Preferences.Forecast)
}