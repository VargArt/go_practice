package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"weather_checker/conf_reader"
	"weather_checker/handlers"
	"weather_checker/redis_adapter"

	"github.com/alexflint/go-arg"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var args struct {
		ConfPath string `arg:"--config"`
	}
	arg.MustParse(&args)

	var config conf_reader.Conf
	config.GetConf(args.ConfPath)

	redis_adapter.Cfg.DialTimeout = time.Duration(config.DialTimeout) * time.Second
	redis_adapter.Cfg.Timeout = time.Duration(config.Timeout) * time.Second
	redis_adapter.Cfg.MaxRetries = int(config.MaxRetries)
	redis_adapter.Cfg.Ttl = time.Duration(config.CacheDuration) * time.Second

	var err error = nil
	redis_adapter.WeatherCache, err = redis_adapter.NewLRUCache(context.Background())
	if err != nil {
		errorLog.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", handlers.Handler)

	infoLog.Print("start server")
	srv := &http.Server{
		Addr:     "localhost:8000",
		ErrorLog: errorLog,
		Handler:  mux,
	}

	errorLog.Fatal(srv.ListenAndServe())
}
