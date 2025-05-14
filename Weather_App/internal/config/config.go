package config

import (
	"encoding/json"
	"os"
)

// AppConfig holds which weather provider to use.
type AppConfig struct {
	WeatherProvider string `json:"weather_provider"`
}

// Load reads and parses the JSON file at path.
func Load(path string) (*AppConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg AppConfig
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
