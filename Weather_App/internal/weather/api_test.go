package weather

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
	"weatherapp/models"

	"github.com/stretchr/testify/assert"
)

// fakeProvider implements WeatherProvider for testing.
type fakeProvider struct {
	currentData  *WeatherData
	forecastData []WeatherData
}

func (f *fakeProvider) Current(location string) (*WeatherData, error) {
	return f.currentData, nil
}

func (f *fakeProvider) Forecast(location string, days int) ([]WeatherData, error) {
	return f.forecastData, nil
}

// TestShowWeather_Day tests the "day" forecast path of ShowWeather.
func TestShowWeather_Day(t *testing.T) {
	// Setup fake provider with sample current data
	f := &fakeProvider{
		currentData: &WeatherData{
			Description: "Sunny",
			Temperature: 20,
			FeelsLike:   18,
			Humidity:    65,
			WindSpeed:   12,
			WindDir:     "NE",
		},
	}
	InitProvider(f)

	// Capture output
	var outBuf bytes.Buffer
	outputWriter = &outBuf

	// Prepare user preferences for a single-day detailed view
	user := models.User{
		Preferences: models.Preferences{
			Location:  "london",
			Unit:      "celsius",
			Verbosity: "verbose",
			Forecast:  "day",
		},
	}

	ShowWeather(user)
	out := outBuf.String()

	assert.Contains(t, out, "Weather for London")
	assert.Contains(t, out, "Description : Sunny")
	assert.Contains(t, out, "Temperature : 20 °C")
	assert.Contains(t, out, "Feels Like  : 18 °C")
	assert.Contains(t, out, "Humidity    : 65%")
	assert.Contains(t, out, "Wind        : 12 km/h (NE)")
}

// TestShowWeather_Week tests the "week" forecast path of ShowWeather.
func TestShowWeather_Week(t *testing.T) {
	// Setup fake provider with sample forecast data for two days
	f := &fakeProvider{
		forecastData: []WeatherData{
			{Description: "Partly Cloudy", Temperature: 15},
			{Description: "Rainy", Temperature: 12},
		},
	}
	InitProvider(f)

	var outBuf bytes.Buffer
	outputWriter = &outBuf

	// Prepare user preferences for a weekly forecast
	user := models.User{
		Preferences: models.Preferences{
			Location:  "london",
			Unit:      "celsius",
			Verbosity: "brief",
			Forecast:  "week",
		},
	}

	ShowWeather(user)
	out := outBuf.String()

	assert.Contains(t, out, "Forecast for London (week)")
	assert.Contains(t, out, "Day 1: Partly Cloudy – 15°C")
	assert.Contains(t, out, "Day 2: Rainy – 12°C")
}

// TestShowOtherLocations tests the interactive ShowOtherLocations function.
func TestShowOtherLocations(t *testing.T) {
	// Setup fake provider with sample current data
	f := &fakeProvider{
		currentData: &WeatherData{Description: "Cloudy", Temperature: 18},
	}
	InitProvider(f)

	var outBuf bytes.Buffer
	outputWriter = &outBuf

	// Simulate user entering "paris" as the location
	reader := bufio.NewReader(strings.NewReader("paris\n"))
	ShowOtherLocations(reader)
	out := outBuf.String()

	assert.Contains(t, out, "Location: Paris | Cloudy | 18°C")
}
