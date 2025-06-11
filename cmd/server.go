package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"weather-simple-api/internal/collector"

	"github.com/joho/godotenv"
)

func init() { // Get environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("Warning: Could not load .env file: %v", err) // Could give a false error when running from Docker
	}
}

func main() { // Start the workers and the server
	collector.StartWorker()
	defer collector.StopWorkers() // The workers will end the tasks when closing the server

	http.HandleFunc("/weather", weatherHandler)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v", err)
	}
}

func weatherHandler(w http.ResponseWriter, r *http.Request) { // Handle the endpoint /weather
	// Get latitude and longitude query parameters
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	// Get the context
	ctx := r.Context()

	if lat == "" || lon == "" {
		http.Error(w, "Missing lat or lon query parameters", http.StatusBadRequest)
		return
	}

	// Fetch weather forecast using the collector package
	data, err := collector.FetchWeatherForecastWorker(ctx, lat, lon)
	if err != nil {
		if ctx.Err() == context.Canceled {
			http.Error(w, "Request canceled by the client", http.StatusRequestTimeout)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Set response headers and encode the data as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
