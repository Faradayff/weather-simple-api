package apis

import (
	"time"
	"weather-simple-api/internal/models"
)

type OpenMeteo struct{}

const OpenMeteoURL = ""

func (o OpenMeteo) Fetch(lat, lon float64, date time.Time) (models.DailyForecast, error) {

	return models.DailyForecast{}, nil
}

func (o OpenMeteo) GetClientName() string {
	return "OpenMeteo"
}
