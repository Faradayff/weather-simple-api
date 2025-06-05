package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"weather-simple-api/internal/models"
)

type OpenMeteo struct{}

const OpenMeteoURL = "https://api.open-meteo.com/v1/forecast"

func (o OpenMeteo) Fetch(lat, lon, date string) (models.DailyForecast, error) {
	var forecast models.DailyForecast
	var respData openMeteoResponse

	startDate := date
	endDate := date
	url := fmt.Sprintf("%s?latitude=%s&longitude=%s&start_date=%s&end_date=%s&daily=temperature_2m_max,temperature_2m_min,uv_index_max", OpenMeteoURL, lat, lon, startDate, endDate)

	resp, err := http.Get(url)
	if err != nil {
		return forecast, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return forecast, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return forecast, err
	}

	if len(respData.Daily.Time) > 0 {
		forecast.Date = respData.Daily.Time[0]
		forecast.MaxTemp = respData.Daily.Temperature2mMax[0]
		forecast.MinTemp = respData.Daily.Temperature2mMin[0]
		forecast.UVIndex = respData.Daily.UVIndexMax[0]
	}
	return forecast, nil
}

func (o OpenMeteo) GetClientName() string {
	return "OpenMeteo"
}

type openMeteoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Daily     struct {
		Time             []string  `json:"time"`
		Temperature2mMax []float64 `json:"temperature_2m_max"`
		Temperature2mMin []float64 `json:"temperature_2m_min"`
		UVIndexMax       []float64 `json:"uv_index_max"`
	} `json:"daily"`
}
