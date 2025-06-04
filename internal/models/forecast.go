package models

type DailyForecast struct {
    Date         string  `json:"date"`
    MaxTempC    float64 `json:"max_temp_c"`
    MinTempC    float64 `json:"min_temp_c"`
	//Humidity    int     `json:"humidity"`
    //Weather     string  `json:"weather"`
}

type WeatherForecast struct {
    WeatherAPI1 map[string]DailyForecast `json:"weatherAPI1"`
    WeatherAPI2 map[string]DailyForecast `json:"weatherAPI2"`
}