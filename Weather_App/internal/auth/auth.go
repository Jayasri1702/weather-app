package auth

import (
	"bufio"
	"fmt"
	"strings"
	"weatherapp/internal/storage"
	"weatherapp/models"

	"golang.org/x/crypto/bcrypt"
)

// Register prompts for user details and saves a new User
func Register(reader *bufio.Reader) {
	fmt.Print("Enter UserID: ")
	userID, _ := reader.ReadString('\n')
	userID = strings.TrimSpace(userID)

	fmt.Print("Enter Name (Username): ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := models.User{
		UserID:      userID,
		Name:        name,
		Password:    string(hash),
		Preferences: models.Preferences{},
	}

	err := storage.SaveUser(user)
	if err != nil {
		fmt.Println("Error saving user:", err)
	} else {
		fmt.Println("User registered successfully!")
	}
}

// Login prompts for credentials, authenticates, and returns the UserID if successful
func Login(reader *bufio.Reader) string {
	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	users := storage.LoadUsers()
	for _, u := range users {
		if u.Name == username {
			err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
			if err == nil {
				fmt.Println("Login successful!")
				return u.UserID
			}
		}
	}

	fmt.Println("Invalid credentials")
	return ""
}
