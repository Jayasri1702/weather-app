package user

import (
	"bufio"
	"fmt"
	"strings"
	"weatherapp/internal/storage"
	"weatherapp/models"
)

// ChangePreferences prompts for and updates a User’s Preferences
func ChangePreferences(reader *bufio.Reader, userID string) {
	users := storage.LoadUsers()
	for _, u := range users {
		if u.UserID == userID {
			promptPreferences(reader, &u)
			storage.UpdateUser(u)
			fmt.Println("Preferences updated")
			return
		}
	}
}

// EnsurePreferences checks if a User has Preferences and prompts if they are empty
func EnsurePreferences(reader *bufio.Reader, userID string) {
	users := storage.LoadUsers()
	for _, u := range users {
		if u.UserID == userID {
			if u.Preferences.Location == "" {
				fmt.Println("\nPlease set your weather preferences:")
				promptPreferences(reader, &u)
				storage.UpdateUser(u)
			}
			return
		}
	}
}

// promptPreferences prompts the user to input and set weather preferences via CLI.
func promptPreferences(reader *bufio.Reader, u *models.User) {
	fmt.Print("Enter your location: ")
	loc, _ := reader.ReadString('\n')
	u.Preferences.Location = strings.TrimSpace(loc)

	fmt.Print("Unit (celsius/fahrenheit): ")
	u.Preferences.Unit, _ = reader.ReadString('\n')
	u.Preferences.Unit = strings.TrimSpace(u.Preferences.Unit)

	fmt.Print("Verbosity (brief/verbose): ")
	u.Preferences.Verbosity, _ = reader.ReadString('\n')
	u.Preferences.Verbosity = strings.TrimSpace(u.Preferences.Verbosity)

	fmt.Print("Forecast (day/week/month): ")
	u.Preferences.Forecast, _ = reader.ReadString('\n')
	u.Preferences.Forecast = strings.TrimSpace(u.Preferences.Forecast)
}

// ListUsers prints all registered users to stdout
func ListUsers() {
	users := storage.LoadUsers()
	for _, u := range users {
		fmt.Printf("UserID: %s, Name: %s\n", u.UserID, u.Name)
	}
}
