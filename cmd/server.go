package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"weather-simple-api/internal/collector"
)

func main() {
    http.HandleFunc("/weather", weatherHandler)
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
    }
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
    latStr := r.URL.Query().Get("lat")
    lonStr := r.URL.Query().Get("lon")
    dateStart := r.URL.Query().Get("date")
    dateEnd := r.URL.Query().Get("date")
    
    dates := make([]time.Time, 2)
    if dateStart == "" && dateEnd == "" { // Default to today and tomorrow if no date is provided
        dates[0] = time.Now()
        dates[1] = dates[1].AddDate(0, 0, 1)
    }  else if dateStart != "" {
        dates[0], _ = time.Parse("02/01/2006", dateStart)
    } else if dateEnd != "" {
        dates[1], _ = time.Parse("02/01/2006", dateEnd)
    } else {
        dates[0], _ = time.Parse("02/01/2006", dateStart)
        dates[1], _ = time.Parse("02/01/2006", dateEnd)
    }

    lat, err := strconv.ParseFloat(latStr, 32)
    if err != nil {
        http.Error(w, "Invalid latitude value", http.StatusBadRequest)
        return
    }

    lon, err := strconv.ParseFloat(lonStr, 32)
    if err != nil {
        http.Error(w, "Invalid longitude value", http.StatusBadRequest)
        return
    }

	data, err := collector.FetchWeatherForecast(lat, lon, dates)
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