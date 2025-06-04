package collector

import (
	"fmt"
	"sync"
	"time"
	"weather-simple-api/internal/apis"
	"weather-simple-api/internal/models"
)

// Add here the available APIs that you want to use for fetching weather data.
// You can add more APIs by implementing the WeatherAPIClient interface in the apis package.
var availableAPIs = []apis.WeatherClient{
	apis.OpenMeteo{},
	apis.WeatherAPI{APIKey: "",},
}

func FetchWeatherForecast(lat, lon float64, dates[] time.Time) ([]models.WeatherForecast, error) {
	var wg sync.WaitGroup
	var mu = make([]sync.Mutex, len(availableAPIs))

	forecasts := make([]models.WeatherForecast, len(availableAPIs))

	for i, api := range availableAPIs {
		forecasts[i].ApiName = api.GetClientName()
		forecasts[i].ForecastList = make(map[time.Time]models.DailyForecast)
		for _, date := range dates {
			wg.Add(1)
			go func() {
				defer wg.Done()
				data, err := api.Fetch(lat, lon, date)
				if err != nil {
					fmt.Println("Error fetching from", api.GetClientName(), ":", err)
					return
				}
				mu[i].Lock()
				forecasts[i].ForecastList[date] = data
				mu[i].Unlock()
			}()
		}
	}

	wg.Wait()

	var total int
	for i := range forecasts {
		total += len(forecasts[i].ForecastList)
	}
	if total == 0 {
		return []models.WeatherForecast{}, fmt.Errorf("no data fetched from any API")
	}
	return forecasts, nil
}
