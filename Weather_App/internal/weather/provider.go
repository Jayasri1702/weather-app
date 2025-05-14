package weather

// WeatherData holds common weather fields.
type WeatherData struct {
	Description string
	Temperature float64
	FeelsLike   float64
	Humidity    float64
	WindSpeed   float64
	WindDir     string
}

// WeatherProvider defines the interface for any weather source.
type WeatherProvider interface {
	Current(location string) (*WeatherData, error)
	Forecast(location string, days int) ([]WeatherData, error)
}
