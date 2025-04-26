package main

import (
	"log"
	"net/http"
	"time"

	"weather_checker/conf_reader"
	"weather_checker/handlers"

	"github.com/alexflint/go-arg"
)

func main() {
	var args struct {
		ConfPath string `arg:"--config"`
	}
	arg.MustParse(&args)

	var config conf_reader.Conf
	config.GetConf(args.ConfPath)

	handlers.DialTimeout = time.Duration(config.DialTimeout) * time.Second
	handlers.Timeout = time.Duration(config.Timeout) * time.Second
	handlers.MaxRetries = int(config.MaxRetries)

	http.HandleFunc("/", handlers.Handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
