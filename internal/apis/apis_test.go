package apis

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestOpenMeteoFetch_CorrectJSON(t *testing.T) {
	mockJSON := `{
		"latitude": 44.35,
		"longitude": 10.983,
		"daily": {
			"time": ["2025-06-06"],
			"temperature_2m_max": [25.5],
			"temperature_2m_min": [15.2],
			"uv_index_max": [7.3]
		}
	}`

	mockClient := &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(mockJSON)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	om := OpenMeteo{
		Client: mockClient,
	}

	forecast, err := om.Fetch("44.35", "10.983", "2025-06-06")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if forecast.Date != "2025-06-06" {
		t.Errorf("expected date '2025-06-06', got %s", forecast.Date)
	}
	if forecast.MaxTemp != 25.5 {
		t.Errorf("expected MaxTemp 25.5, got %v", forecast.MaxTemp)
	}
	if forecast.MinTemp != 15.2 {
		t.Errorf("expected MinTemp 15.2, got %v", forecast.MinTemp)
	}
	if forecast.UVIndex != 7.3 {
		t.Errorf("expected UVIndex 7.3, got %v", forecast.UVIndex)
	}
}

func TestWeatherAPIFetch_CorrectJSON(t *testing.T) {
	mockJSON := `{
		"forecast": {
			"forecastday": [
				{
					"date": "2025-06-06",
					"day": {
						"maxtemp_c": 25.5,
						"mintemp_c": 15.2,
						"uv": 7.3
					}
				}
			]
		}
	}`

	mockClient := &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(mockJSON)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	w := WeatherAPI{
		Client: mockClient,
	}

	forecast, err := w.Fetch("44.35", "10.983", "2025-06-06")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if forecast.Date != "2025-06-06" {
		t.Errorf("expected date '2025-06-06', got %s", forecast.Date)
	}
	if forecast.MaxTemp != 25.5 {
		t.Errorf("expected MaxTemp 25.5, got %v", forecast.MaxTemp)
	}
	if forecast.MinTemp != 15.2 {
		t.Errorf("expected MinTemp 15.2, got %v", forecast.MinTemp)
	}
	if forecast.UVIndex != 7.3 {
		t.Errorf("expected UVIndex 7.3, got %v", forecast.UVIndex)
	}
}
