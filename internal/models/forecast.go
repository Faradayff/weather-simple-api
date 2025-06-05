package models

type DailyForecast struct {
	Date    string  `json:"date"`
	MaxTemp float64 `json:"max_temp"`
	MinTemp float64 `json:"min_temp"`
	UVIndex float64 `json:"uv_index,omitempty"` // Optional field for UV index
}
