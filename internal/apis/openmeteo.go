package apis

import (
	"time"
	"weather-simple-api/internal/models"
)

type OpenMeteo struct{}

const OpenMeteoURL = "https://api.open-meteo.com/v1/forecast" //?latitude=52.52&longitude=13.41&start_date=2024-10-15&end_date=2024-10-15&daily=temperature_2m_max

func (o OpenMeteo) Fetch(lat, lon float64, date time.Time) (models.DailyForecast, error) {

	return models.DailyForecast{}, nil
}

func (o OpenMeteo) GetClientName() string {
	return "OpenMeteo"
}
