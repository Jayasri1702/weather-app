package storage

import (
	"context"
	"log"
	"os"

	"weatherapp/models"

	"cloud.google.com/go/firestore"
)

var (
	Client *firestore.Client

	// SaveUser writes a User into the "users" collection
	SaveUser = func(u models.User) error {
		ctx := context.Background()
		_, err := Client.Collection("users").Doc(u.UserID).Set(ctx, u)
		return err
	}

	// LoadUsers reads all users from Firestore
	LoadUsers = func() []models.User {
		ctx := context.Background()
		iter := Client.Collection("users").Documents(ctx)
		docs, err := iter.GetAll()
		if err != nil {
			return []models.User{}
		}
		var users []models.User
		for _, doc := range docs {
			var u models.User
			doc.DataTo(&u)
			users = append(users, u)
		}
		return users
	}

	// UpdateUser overwrites a single user document in Firestore
	UpdateUser = func(u models.User) error {
		ctx := context.Background()
		_, err := Client.Collection("users").Doc(u.UserID).Set(ctx, u)
		return err
	}
)

// InitFirestore initializes the Firestore client
func InitFirestore() {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal("GOOGLE_CLOUD_PROJECT must be set")
	}
	var err error
	Client, err = firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to connect to Firestore: %v", err)
	}
}

// GetUserByID fetches one user document by its ID
func GetUserByID(userID string) (*models.User, error) {
	ctx := context.Background()
	doc, err := Client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return nil, err
	}
	var u models.User
	if err := doc.DataTo(&u); err != nil {
		return nil, err
	}
	return &u, nil
}
