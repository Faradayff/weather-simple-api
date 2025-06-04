package apis

import (
	"time"
	"weather-simple-api/internal/models"
)

// WeatherClient define la interfaz que deben implementar todas las APIs de clima.
type WeatherClient interface {
	Fetch(lat, lon float64, date time.Time) (models.DailyForecast, error)
	GetClientName() string
}
