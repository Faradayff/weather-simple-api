package collector

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
	"weather-simple-api/internal/apis"
	"weather-simple-api/internal/models"
)

type ForecastTask struct {
	Lat, Lon string
	Result   chan map[string]map[string]models.DailyForecast
	Err      chan error
}

var (
	taskQueue = make(chan ForecastTask)
	once      sync.Once
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
)

const nWorkers int = 3 // Number of concurrent workers

// Add here the available APIs that you want to use for fetching weather data.
// You can add more APIs by implementing the WeatherAPIClient interface in the apis package.
var availableAPIs = []apis.WeatherClient{
	apis.OpenMeteo{},
	apis.WeatherAPI{APIKey: os.Getenv("WEATHER_API_KEY")},
}

func StartWorker() {
	once.Do(func() {
		ctx, cancel = context.WithCancel(context.Background()) // Start context with cancel function
		for range nWorkers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done(): // Checks if the context is canceled, signaling the worker to stop.
						return
					case task := <-taskQueue: // Checks if there is a task in the taskQueue channel.
						result, err := fetchWeatherForecast(task.Lat, task.Lon) // Processes the task by fetching the weather forecast.
						if err != nil {
							task.Err <- err // Sends the error back through the task's error channel.
						} else {
							task.Result <- result // Sends the result back through the task's result channel.
						}
					}
				}
			}()
		}
	})
}

func StopWorkers() { // End function to stop the workers
	cancel()
	wg.Wait()
}

func FetchWeatherForecastWorker(lat, lon string) (map[string]map[string]models.DailyForecast, error) { // Send the task to thje workers
	resultChan := make(chan map[string]map[string]models.DailyForecast)
	errChan := make(chan error)
	taskQueue <- ForecastTask{Lat: lat, Lon: lon, Result: resultChan, Err: errChan}
	select {
	case res := <-resultChan:
		return res, nil
	case err := <-errChan:
		return nil, err
	}
}

func fetchWeatherForecast(lat, lon string) (map[string]map[string]models.DailyForecast, error) { // Get the weather forecast (worker's job)
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
			go func(i, j int, api apis.WeatherClient, date string) {
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
			}(i, j, api, date)
		}
	}

	wg.Wait()
	return forecasts, nil
}
