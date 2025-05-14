package storage
 
import (
    "context"
    "errors"
    "testing"
    "weatherapp/models"
 
    "cloud.google.com/go/firestore"
    "github.com/stretchr/testify/assert"
)
 
// Mocked Firestore Client Components
 
type mockDocRef struct{}
 

type mockClient struct{}


 // Set mocks setting data in Firestore (used for testing).
func (m *mockDocRef) Set(_ context.Context, _ interface{}, _ ...firestore.SetOption) (*firestore.WriteResult, error) {
    return nil, nil
}

// Collection mocks retrieving a Firestore collection (used for testing).
func (m *mockClient) Collection(name string) *firestore.CollectionRef {
    return &firestore.CollectionRef{}
}

// TestSaveUser_Fake verifies SaveUser behavior for both success and failure cases.
func TestSaveUser_Fake(t *testing.T) {
    Client = &firestore.Client{}
    SaveUser = func(user models.User) error {
        if user.UserID == "fail" {
            return errors.New("simulated Firestore error")
        }
        return nil
    }

    t.Run("Successful save", func(t *testing.T) {
        err := SaveUser(models.User{UserID: "123"})
        assert.NoError(t, err)
    })

    t.Run("Simulated Firestore error", func(t *testing.T) {
        err := SaveUser(models.User{UserID: "fail"})
        assert.Error(t, err)
    })
}

// TestLoadUsers_Fake verifies mock LoadUsers returns expected user data.
func TestLoadUsers_Fake(t *testing.T) {
    LoadUsers = func() []models.User {
        return []models.User{
            {UserID: "u1", Name: "Test User"},
        }
    }

    users := LoadUsers()
    assert.Len(t, users, 1)
    assert.Equal(t, "u1", users[0].UserID)
}

// TestUpdateUser_Fake checks UpdateUser handling for success and failure.
func TestUpdateUser_Fake(t *testing.T) {
    UpdateUser = func(user models.User) error {
        if user.UserID == "fail" {
            return errors.New("update failed")
        }
        return nil
    }

    t.Run("Update success", func(t *testing.T) {
        err := UpdateUser(models.User{UserID: "ok"})
        assert.NoError(t, err)
    })

    t.Run("Update failure", func(t *testing.T) {
        err := UpdateUser(models.User{UserID: "fail"})
        assert.Error(t, err)
    })
}