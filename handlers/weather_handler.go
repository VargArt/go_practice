package handlers

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"weather_checker/redis_adapter"
	"weather_checker/requesters"
	"weather_checker/types"

	"github.com/redis/go-redis/v9"
)

func GetMarshaller(contentType string) func(v interface{}) ([]byte, error) {
	if contentType == "application/json" {
		return json.Marshal
	} else if contentType == "application/xml" {
		return func(v interface{}) ([]byte, error) {
			return xml.MarshalIndent(v, "", " ")
		}
	}

	return json.Marshal
}

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Print("start handler")
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		w.Write([]byte("405 - method not allowed\n"))
		return
	}

	requested_content_type := r.Header.Get("Content-Type")
	city := r.URL.Query().Get("city")
	if city == "" {
		log.Printf("invalid parameter")
		w.WriteHeader(400)
		w.Write([]byte("400 - invalid parameter\n"))
		return
	}

	marshal_func := GetMarshaller(requested_content_type)

	if redis_adapter.WeatherCache == nil {
		log.Print("invalid cache")
		panic("invalid cache")
	}

	db_data, err := redis_adapter.WeatherCache.Get(context.Background(), city)
	if err == redis.Nil {
		log.Print("empty cache, use api request")
		api_weather := requesters.GetWeather(city)

		data, err := marshal_func(api_weather)
		if err != nil {
			log.Print("invalid marshaling")
			w.WriteHeader(400)
			w.Write([]byte("400 - invalid answer\n"))
			return
		}

		w.Header().Set("Content-Type", requested_content_type)
		w.Write(data)

		if err := redis_adapter.WeatherCache.Set(context.Background(), city, data); err != nil {
			log.Printf("failed to set data, error: %s", err.Error())
			w.WriteHeader(400)
			w.Write([]byte("400 - invalid answer\n"))
			return
		}

		return

	} else if err != nil {
		fmt.Printf("failed to get value, error: %v\n", err)
	}

	cached_weather := types.WeatherUnmarshal([]byte(db_data))

	data, err := marshal_func(cached_weather)
	if err != nil {
		log.Print(err.Error())
		panic(err)
	}

	w.Header().Set("Content-Type", requested_content_type)
	w.Write(data)
}
