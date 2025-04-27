package requesters

import (
	"fmt"
	"io"
	"net/http"

	"weather_checker/types"
)

func GetWeather(city string) []byte {
	query := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=  &q=%s&days=5&aqi=no&alerts=no", city)

	resp, err := http.Get(query)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return types.WeatherUnmarshal([]byte(body))
}
