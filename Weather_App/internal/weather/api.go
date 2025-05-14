package weather

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"weatherapp/models"

	"github.com/joho/godotenv"
)

var (
	provider     WeatherProvider
	outputWriter io.Writer
)

// InitProvider sets the active WeatherProvider.
func InitProvider(p WeatherProvider) {
	provider = p
}

// ShowWeather uses the configured provider to display either a one-day detailed view or a multi-day forecast
func ShowWeather(user models.User) {
	loc := user.Preferences.Location
	unit := strings.ToLower(user.Preferences.Unit)
	verbosity := strings.ToLower(user.Preferences.Verbosity)
	forecast := strings.ToLower(user.Preferences.Forecast)

	if forecast == "day" {
		data, err := provider.Current(loc)
		if err != nil {
			fmt.Fprintf(getWriter(), "Error: %v\n", err)
			return
		}
		renderDetailed(loc, data, verbosity, unit)
	} else {
		days := 7
		if forecast == "month" {
			days = 30
		}
		dataSlice, err := provider.Forecast(loc, days)
		if err != nil {
			fmt.Fprintf(getWriter(), "Error: %v\n", err)
			return
		}
		renderForecast(loc, dataSlice, verbosity, unit, forecast)
	}
}

// ShowOtherLocations prompts and then shows current weather for one city.
func ShowOtherLocations(reader *bufio.Reader) {
	fmt.Print("Enter location: ")
	loc, _ := reader.ReadString('\n')
	loc = strings.TrimSpace(loc)
	if loc == "" {
		fmt.Fprintln(getWriter(), "No location entered.")
		return
	}
	data, err := provider.Current(loc)
	if err != nil {
		fmt.Fprintf(getWriter(), "Error: %v\n", err)
		return
	}
	fmt.Fprintf(getWriter(), "Location: %s | %s | %.0f°C\n",
		strings.Title(loc), data.Description, data.Temperature)
}

// to print detailed view
func renderDetailed(loc string, d *WeatherData, verbosity, unit string) {
	unitLabel := "°C"
	temp := d.Temperature
	feels := d.FeelsLike
	if unit == "fahrenheit" {
		temp = temp*9/5 + 32
		feels = feels*9/5 + 32
		unitLabel = "°F"
	}
	out := getWriter()
	fmt.Fprintf(out, "\n Weather for %s\n", strings.Title(loc))
	fmt.Fprintln(out, "------------------------")
	fmt.Fprintf(out, "Description : %s\n", d.Description)
	fmt.Fprintf(out, "Temperature : %.0f %s\n", temp, unitLabel)
	if verbosity == "verbose" {
		fmt.Fprintf(out, "Feels Like  : %.0f %s\n", feels, unitLabel)
		fmt.Fprintf(out, "Humidity    : %.0f%%\n", d.Humidity)
		fmt.Fprintf(out, "Wind        : %.0f km/h (%s)\n", d.WindSpeed, d.WindDir)
	}
}

// to print multi-day forecast
func renderForecast(loc string, data []WeatherData, verbosity, unit, label string) {
	unitLabel := "°C"
	if unit == "fahrenheit" {
		unitLabel = "°F"
	}
	out := getWriter()
	fmt.Fprintf(out, "\n Forecast for %s (%s)\n", strings.Title(loc), label)
	fmt.Fprintln(out, "----------------------------")
	for i, d := range data {
		temp := d.Temperature
		if unit == "fahrenheit" {
			temp = temp*9/5 + 32
		}
		fmt.Fprintf(out, "Day %d: %s – %.0f%s\n", i+1, d.Description, temp, unitLabel)
		if verbosity == "verbose" {
			feels := d.FeelsLike
			if unit == "fahrenheit" {
				feels = feels*9/5 + 32
			}
			fmt.Fprintf(out, "  Feels like : %.0f%s\n", feels, unitLabel)
			fmt.Fprintf(out, "  Humidity   : %.0f%%\n", d.Humidity)
			fmt.Fprintf(out, "  Wind       : %.0f km/h\n", d.WindSpeed)
		}
	}
}

func init() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
}

func getWriter() io.Writer {
	if outputWriter != nil {
		return outputWriter
	}
	return io.Writer(os.Stdout)
}
