package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"weather_checker/redis_adapter"

	"github.com/redis/go-redis/v9"
)

var DialTimeout time.Duration = 10 * time.Second
var Timeout time.Duration = 5 * time.Second
var MaxRetries = 5

func Handler(w http.ResponseWriter, r *http.Request) {
	cfg := redis_adapter.Config{
		Addr:        "localhost:6379",
		Password:    "",
		User:        "",
		DB:          0,
		MaxRetries:  MaxRetries,
		DialTimeout: DialTimeout,
		Timeout:     Timeout,
	}

	db, err := redis_adapter.NewClient(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	val, err := db.Get(context.Background(), "key").Result()
	if err == redis.Nil {
		fmt.Println("value not found")
	} else if err != nil {
		fmt.Printf("failed to get value, error: %v\n", err)
	}

	fmt.Fprintf(w, "%q\n", val)

	resp, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key={}&q=Moscow&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "%q\n", body)
}
