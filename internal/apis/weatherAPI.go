package apis

import (
	"time"
	"weather-simple-api/internal/models"
)

type WeatherAPI struct {
	APIKey string
}

const WeatherAPIURL = "https://api.weatherapi.com/v1/forecast.json" //?key={API_KEY}&q=44.34,10.99&date=2024-10-15&day=maxtemp_c

func (w WeatherAPI) Fetch(lat, lon float64, date time.Time) (models.DailyForecast, error) {

	return models.DailyForecast{}, nil
}

func (w WeatherAPI) GetClientName() string {
	return "WeatherAPI"
}
