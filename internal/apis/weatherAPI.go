package apis

import (
	"time"
	"weather-simple-api/internal/models"
)

type WeatherAPI struct {
    APIKey string
}

const WeatherAPIURL = ""

func (w WeatherAPI) Fetch(lat, lon float32, date time.Time) (models.DailyForecast, error) {
    
    return models.DailyForecast{}, nil
}

func (w WeatherAPI) GetClientName() string {
    return "WeatherAPI"
}