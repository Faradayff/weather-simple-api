package collector

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"weather-simple-api/internal/apis"
	"weather-simple-api/internal/models"
)

// Add here the available APIs that you want to use for fetching weather data.
// You can add more APIs by implementing the WeatherAPIClient interface in the apis package.
var availableAPIs = []apis.WeatherClient{
	apis.OpenMeteo{},
	apis.WeatherAPI{APIKey: ""},
}

func FetchWeatherForecast(lat, lon float64) (map[string]map[string]models.DailyForecast, error) {
	var wg sync.WaitGroup
	var mu = make([]sync.Mutex, len(availableAPIs))

	forecasts := make(map[string]map[string]models.DailyForecast)
	dates := make([]time.Time, 5)
	for i := range dates {
		if i == 0 {
			dates[0] = time.Now()
		} else {
			dates[i] = dates[i-1].AddDate(0, 0, 1)
		}
	}

	for i, api := range availableAPIs {
		forecasts[api.GetClientName()] = make(map[string]models.DailyForecast)
		for j, date := range dates {
			wg.Add(1)
			go func() {
				defer wg.Done()
				data, err := api.Fetch(lat, lon, date)
				if err != nil {
					fmt.Println("Error fetching from", api.GetClientName(), ":", err)
					return
				}
				day := "day" + strconv.Itoa(j+1)

				mu[i].Lock()
				forecasts[api.GetClientName()][day] = data
				mu[i].Unlock()
			}()
		}
	}

	wg.Wait()
	return forecasts, nil
}
