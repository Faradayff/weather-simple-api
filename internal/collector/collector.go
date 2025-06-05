package collector

import (
	"fmt"
	"os"
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
	apis.WeatherAPI{APIKey: os.Getenv("WEATHER_API_KEY")},
}

func FetchWeatherForecast(lat, lon string) (map[string]map[string]models.DailyForecast, error) {
	var wg sync.WaitGroup
	var mu = make([]sync.Mutex, len(availableAPIs))

	forecasts := make(map[string]map[string]models.DailyForecast)
	dates := make([]string, 5)
	for i := range dates {
		if i == 0 {
			dates[0] = time.Now().Format("2006-01-02")
		} else {
			dates[i] = time.Now().AddDate(0, 0, i).Format("2006-01-02")
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
