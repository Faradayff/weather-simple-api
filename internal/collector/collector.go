package collector

import (
	"context"
	"fmt"
	"strconv"
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

var taskQueue = make(chan ForecastTask)

func FetchWeatherForecastWorker(ctx context.Context, tm *TaskManager, lat, lon string) (map[string]map[string]models.DailyForecast, error) { // Send the task to the workers
	// Create channels for results and errors
	resultChan := make(chan ForecastResult)
	errChan := make(chan error)
	availableAPIs := ctx.Value("availableAPIs").([]apis.WeatherClient)

	// Send tasks to the workers
	for _, api := range availableAPIs {
		for i := range 5 {
			select {
			case <-ctx.Done(): // Continue working until the context is canceled
				return nil, ctx.Err()
			case tm.taskQueue <- ForecastTask{Api: api, Lat: lat, Lon: lon, Day: i, Result: resultChan, Err: errChan}:
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
