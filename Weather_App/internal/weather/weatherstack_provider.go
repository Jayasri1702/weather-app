package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type WeatherstackProvider struct {
	apiKey string
}

// NewWeatherstackProvider creates a WeatherstackProvider.
func NewWeatherstackProvider() *WeatherstackProvider {
	return &WeatherstackProvider{apiKey: os.Getenv("WEATHERSTACK_API_KEY")}
}

// Current fetches current weather from Weatherstack.
func (w *WeatherstackProvider) Current(location string) (*WeatherData, error) {
	url := fmt.Sprintf("http://api.weatherstack.com/current?access_key=%s&query=%s", w.apiKey, location)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r struct {
		Current struct {
			Temperature  float64  `json:"temperature"`
			FeelsLike    float64  `json:"feelslike"`
			Humidity     float64  `json:"humidity"`
			WindSpeed    float64  `json:"wind_speed"`
			WindDir      string   `json:"wind_dir"`
			Descriptions []string `json:"weather_descriptions"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	cd := r.Current
	return &WeatherData{
		Description: cd.Descriptions[0],
		Temperature: cd.Temperature,
		FeelsLike:   cd.FeelsLike,
		Humidity:    cd.Humidity,
		WindSpeed:   cd.WindSpeed,
		WindDir:     cd.WindDir,
	}, nil
}

// Forecast simulates a multi-day forecast (Weatherstack free tier lack real forecast)
func (w *WeatherstackProvider) Forecast(location string, days int) ([]WeatherData, error) {
	var out []WeatherData
	for i := 1; i <= days; i++ {
		out = append(out, WeatherData{
			Description: "Partly Cloudy",
			Temperature: 20 + float64(i%5),
			FeelsLike:   20 + float64(i%3),
			Humidity:    70,
			WindSpeed:   10,
			WindDir:     "NW",
		})
	}
	return out, nil
}
