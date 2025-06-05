package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"weather-simple-api/internal/models"
)

type WeatherAPI struct {
	APIKey string
}

const WeatherAPIURL = "https://api.weatherapi.com/v1/forecast.json"

func (w WeatherAPI) Fetch(lat, lon, date string) (models.DailyForecast, error) {
	var forecast models.DailyForecast
	var respData weatherAPIResponse

	url := fmt.Sprintf("%s?key=%s&q=%s,%s&date=%s&day=maxtemp_c", WeatherAPIURL, w.APIKey, lat, lon, date)
	fmt.Println(url)
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

	if len(respData.Forecast.Forecastday) > 0 {
		forecast.Date = respData.Forecast.Forecastday[0].Date
		forecast.MaxTemp = respData.Forecast.Forecastday[0].Day.MaxtempC
		forecast.MinTemp = respData.Forecast.Forecastday[0].Day.MintempC
		forecast.UVIndex = respData.Forecast.Forecastday[0].Day.Uv
	}
	return forecast, nil
}

func (w WeatherAPI) GetClientName() string {
	return "WeatherAPI"
}

type weatherAPIResponse struct {
	Forecast struct {
		Forecastday []struct {
			Date string `json:"date"`
			Day  struct {
				MaxtempC float64 `json:"maxtemp_c"`
				MintempC float64 `json:"mintemp_c"`
				Uv       float64 `json:"uv"`
			} `json:"day"`
			// Puedes agregar "astro" y "hour" si los necesitas
		} `json:"forecastday"`
	} `json:"forecast"`
}
