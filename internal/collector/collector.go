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
	Api      apis.WeatherClient
	Lat, Lon string
	Day      int
	Result   chan ForecastResult
	Err      chan error
}

type ForecastResult struct {
	Api      string
	Date     string
	Forecast models.DailyForecast
}

var (
	taskQueue = make(chan ForecastTask)
	once      sync.Once
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
)

var nWorkers int = len(availableAPIs) * 5 * 2 // Number of concurrent workers

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
						result, err := fetchWeatherForecast(task.Api, task.Lat, task.Lon, task.Day) // Processes the task by fetching the weather forecast.
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

func FetchWeatherForecastWorker(ctx context.Context, lat, lon string) (map[string]map[string]models.DailyForecast, error) { // Send the task to the workers
	// Create channels for results and errors
	resultChan := make(chan ForecastResult)
	errChan := make(chan error)

	// Send tasks to the workers
	for _, api := range availableAPIs {
		for i := range 5 {
			select {
			case <-ctx.Done(): // Continue working until the context is canceled
				return nil, ctx.Err()
			case taskQueue <- ForecastTask{Api: api, Lat: lat, Lon: lon, Day: i, Result: resultChan, Err: errChan}:
			}
		}
	}

	forecast := make(map[string]map[string]models.DailyForecast)
	errorCount := 0
	errors := make([]error, 0)

	// Wait for results or errors
	for range len(availableAPIs) * 5 {
		select {
		case <-ctx.Done(): // Continue working until the context is canceled
			return nil, ctx.Err()
		case res := <-resultChan:
			if _, ok := forecast[res.Api]; !ok {
				forecast[res.Api] = make(map[string]models.DailyForecast)
			}
			forecast[res.Api][res.Date] = res.Forecast
		case err := <-errChan:
			errorCount++
			errors = append(errors, err)
			if errorCount > len(availableAPIs)*5/2 { // If the numbers of errors is greater than half of the total tasks, return an error
				return nil, fmt.Errorf("failed to fetch data from all APIs: %v", errors)
			}
		}
	}

	return forecast, nil
}

func fetchWeatherForecast(api apis.WeatherClient, lat, lon string, day int) (ForecastResult, error) { // Get the one day weather forecast (worker's job)
	date := time.Now().AddDate(0, 0, day).Format("2006-01-02")
	result, err := api.Fetch(lat, lon, date)
	if err != nil {
		return ForecastResult{}, fmt.Errorf("error fetching data from %s: %w", api.GetClientName(), err)
	}

	return ForecastResult{Api: api.GetClientName(), Date: "day" + strconv.Itoa(day+1), Forecast: result}, nil
}
