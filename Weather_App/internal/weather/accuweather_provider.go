package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// AccuWeatherProvider implements WeatherProvider using AccuWeather APIs
type AccuWeatherProvider struct {
	apiKey string
}

// NewAccuWeatherProvider reads ACCUWEATHER_API_KEY from the environment
func NewAccuWeatherProvider() *AccuWeatherProvider {
	return &AccuWeatherProvider{apiKey: os.Getenv("ACCUWEATHER_API_KEY")}
}

// lookupLocationKey finds the AccuWeather location key for a city
func (a *AccuWeatherProvider) lookupLocationKey(location string) (string, error) {
	url := fmt.Sprintf(
		"http://dataservice.accuweather.com/locations/v1/cities/search?apikey=%s&q=%s",
		a.apiKey, location,
	)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var locs []struct{ Key string }
	if err := json.NewDecoder(resp.Body).Decode(&locs); err != nil {
		return "", err
	}
	if len(locs) == 0 {
		return "", fmt.Errorf("location not found: %s", location)
	}
	return locs[0].Key, nil
}

// Current fetches the current conditions for a location.
func (a *AccuWeatherProvider) Current(location string) (*WeatherData, error) {
	key, err := a.lookupLocationKey(location)
	if err != nil {
		return nil, err
	}
	condURL := fmt.Sprintf(
		"http://dataservice.accuweather.com/currentconditions/v1/%s?apikey=%s&details=true",
		key, a.apiKey,
	)
	resp, err := http.Get(condURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cs []struct {
		WeatherText string `json:"WeatherText"`
		Temperature struct {
			Metric struct{ Value float64 } `json:"Metric"`
		} `json:"Temperature"`
		RealFeelTemperature struct {
			Minimum struct{ Value float64 } `json:"Minimum"`
			Maximum struct{ Value float64 } `json:"Maximum"`
		} `json:"RealFeelTemperature"`
		RelativeHumidity float64 `json:"RelativeHumidity"`
		Wind             struct {
			Speed     struct{ Value float64 }    `json:"Speed"`
			Direction struct{ Localized string } `json:"Direction"`
		} `json:"Wind"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&cs); err != nil {
		return nil, err
	}
	if len(cs) == 0 {
		return nil, fmt.Errorf("no current conditions for %s", location)
	}
	c := cs[0]
	return &WeatherData{
		Description: c.WeatherText,
		Temperature: c.Temperature.Metric.Value,
		FeelsLike:   c.RealFeelTemperature.Maximum.Value,
		Humidity:    c.RelativeHumidity,
		WindSpeed:   c.Wind.Speed.Value,
		WindDir:     c.Wind.Direction.Localized,
	}, nil
}

// Forecast retrieves up to 5-day forecasts, padded to the requested days.
func (a *AccuWeatherProvider) Forecast(location string, days int) ([]WeatherData, error) {
	key, err := a.lookupLocationKey(location)
	if err != nil {
		return nil, err
	}

	// AccuWeather supports daily forecasts up to 5 days
	requestDays := days
	if requestDays > 5 {
		requestDays = 5
	}
	url := fmt.Sprintf(
		"http://dataservice.accuweather.com/forecasts/v1/daily/%dday/%s?apikey=%s&metric=true&details=true",
		requestDays, key, a.apiKey,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r struct {
		DailyForecasts []struct {
			Temperature struct {
				Minimum struct{ Value float64 } `json:"Minimum"`
				Maximum struct{ Value float64 } `json:"Maximum"`
			} `json:"Temperature"`
			RealFeelTemperature struct {
				Minimum struct{ Value float64 } `json:"Minimum"`
				Maximum struct{ Value float64 } `json:"Maximum"`
			} `json:"RealFeelTemperature"`
			Day struct {
				IconPhrase               string  `json:"IconPhrase"`
				PrecipitationProbability float64 `json:"PrecipitationProbability"`
				Wind                     struct {
					Speed     struct{ Value float64 }    `json:"Speed"`
					Direction struct{ Localized string } `json:"Direction"`
				} `json:"Wind"`
			} `json:"Day"`
		} `json:"DailyForecasts"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	var out []WeatherData
	for i, fc := range r.DailyForecasts {
		if i >= requestDays {
			break
		}
		out = append(out, WeatherData{
			Description: fc.Day.IconPhrase,
			Temperature: fc.Temperature.Maximum.Value,
			FeelsLike:   fc.RealFeelTemperature.Maximum.Value,
			Humidity:    fc.Day.PrecipitationProbability,
			WindSpeed:   fc.Day.Wind.Speed.Value,
			WindDir:     fc.Day.Wind.Direction.Localized,
		})
	}
	// pad for days beyond those returned
	for i := requestDays; i < days; i++ {
		out = append(out, WeatherData{Description: "Forecast unavailable"})
	}
	return out, nil
}
