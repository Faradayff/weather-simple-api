package apis

import (
	"weather-simple-api/internal/models"
)

// WeatherClient define la interfaz que deben implementar todas las APIs de clima.
type WeatherClient interface {
	Fetch(lat, lon, date string) (models.DailyForecast, error)
	GetClientName() string
}
