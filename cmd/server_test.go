package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Handler function to create a weather handler with a mock fetchWeather function
func createWeatherHandler(fetchWeather func(lat, lon string) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lat := r.URL.Query().Get("lat")
		lon := r.URL.Query().Get("lon")

		data, err := fetchWeather(lat, lon)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func TestWeatherHandler_MissingParams(t *testing.T) {
	req := httptest.NewRequest("GET", "/weather", nil)
	w := httptest.NewRecorder()
	weatherHandler(w, req, nil) // Pass nil or a mock TaskManager as appropriate
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestWeatherHandler_Success(t *testing.T) {
	mockFetchWeather := func(lat, lon string) (any, error) {
		return map[string]string{"weather": "sunny"}, nil
	}

	handler := createWeatherHandler(mockFetchWeather)

	req := httptest.NewRequest("GET", "/weather?lat=1&lon=1", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	if response["weather"] != "sunny" {
		t.Errorf("expected 'sunny', got '%s'", response["weather"])
	}
}

func TestWeatherHandler_InternalServerError(t *testing.T) {
	mockFetchWeather := func(lat, lon string) (any, error) {
		return nil, errors.New("simulated error")
	}

	handler := createWeatherHandler(mockFetchWeather)

	req := httptest.NewRequest("GET", "/weather?lat=1&lon=1", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
