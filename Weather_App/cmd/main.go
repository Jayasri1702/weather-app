package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"weatherapp/internal/auth"
	"weatherapp/internal/config"
	"weatherapp/internal/storage"
	"weatherapp/internal/user"
	"weatherapp/internal/weather"
)

func init() {

	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("No .env found; relying on environment variables")
		}
	}
}

func main() {
	// Load app config
	cfg, err := config.Load("../config/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize the chosen weather provider
	switch strings.ToLower(cfg.WeatherProvider) {
	case "accuweather":
		weather.InitProvider(weather.NewAccuWeatherProvider())
	default:
		weather.InitProvider(weather.NewWeatherstackProvider())
	}

	// Initialize Firestore
	storage.InitFirestore()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n=== Weather CLI App ===")
		fmt.Println("1. Register")
		fmt.Println("2. Login")
		fmt.Println("3. Exit")
		fmt.Print("Enter choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			auth.Register(reader)

		case "2":
			userID := auth.Login(reader)
			if userID != "" {
				user.EnsurePreferences(reader, userID)
				dashboard(reader, userID)
			}

		case "3":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid option")
		}
	}
}

func dashboard(reader *bufio.Reader, userID string) {
	for {
		fmt.Println("\n=== Dashboard ===")
		fmt.Println("1. View My Weather")
		fmt.Println("2. Change Preferences")
		fmt.Println("3. View Other Locations")
		fmt.Println("4. List Users")
		fmt.Println("5. Logout")
		fmt.Print("Enter choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			u, err := storage.GetUserByID(userID)
			if err != nil {
				fmt.Println("Error fetching user:", err)
				continue
			}
			weather.ShowWeather(*u)

		case "2":
			user.ChangePreferences(reader, userID)

		case "3":
			weather.ShowOtherLocations(reader)

		case "4":
			user.ListUsers()

		case "5":
			return

		default:
			fmt.Println("Invalid choice")
		}
	}
}
