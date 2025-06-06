package collector

import (
	"testing"
	"weather-simple-api/internal/apis"
	"weather-simple-api/internal/models"
)

type mockAPI struct{}

func (m mockAPI) Fetch(lat, lon, date string) (models.DailyForecast, error) {
	return models.DailyForecast{
		Date:    date,
		MaxTemp: 20.0,
		MinTemp: 10.0,
		UVIndex: 5.0,
	}, nil
}

func (m mockAPI) GetClientName() string {
	return "mockAPI"
}

func TestFetchWeatherForecastWorker(t *testing.T) {
	availableAPIs = []apis.WeatherClient{
		mockAPI{},
	}

	result, err := FetchWeatherForecastWorker("1", "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	apiName := "mockAPI"
	forecasts, ok := result[apiName]
	if !ok {
		t.Fatalf("expected key %s in results", apiName)
	}

	if len(forecasts) != 5 {
		t.Fatalf("expected 5 days for %s, got %d", apiName, len(forecasts))
	}

	for dayKey, fc := range forecasts {
		if fc.MaxTemp != 20.0 {
			t.Errorf("expected MaxTemp 20.0 for %s, got %v", dayKey, fc.MaxTemp)
		}
		if fc.MinTemp != 10.0 {
			t.Errorf("expected MinTemp 10.0 for %s, got %v", dayKey, fc.MinTemp)
		}
		if fc.UVIndex != 5.0 {
			t.Errorf("expected UVIndex 5.0 for %s, got %v", dayKey, fc.UVIndex)
		}
	}
}
