package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"weather-simple-api/internal/collector"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Could not load .env file: %v", err) // Could give a false error when running from Docker
	}
}

func main() {
	http.HandleFunc("/weather", weatherHandler)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v", err)
	}
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	if lat == "" || lon == "" {
		http.Error(w, "Missing lat or lon query parameters", http.StatusBadRequest)
		return
	}

	data, err := collector.FetchWeatherForecastWorker(lat, lon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
