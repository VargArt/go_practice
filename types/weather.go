package types

import "encoding/json"

type ForecastDay struct {
	Date string `json:"date"`
	Day  struct {
		MaxtempC  float32 `json:"maxtemp_c"`
		MintempC  float32 `json:"mintemp_c"`
		AvgtempC  float32 `json:"avgtemp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"day"`
}

type Weather struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
	Forecast struct {
		Forecastday []ForecastDay `json:"forecastday"`
	} `json:"forecast"`
}

func WeatherUnmarshal(val []byte) Weather {
	var weather Weather
	err := json.Unmarshal((val), &weather)
	if err != nil {
		panic(err)
	}

	return weather
}
