package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"weather_checker/conf_reader"
	"weather_checker/handlers"
	"weather_checker/redis_adapter"
)

type configFlag struct {
	ConfigPath string
}

func (f *configFlag) Set(s string) error {
	if s == "" {
		return fmt.Errorf("empty config path")
	}
	f.ConfigPath = s

	return nil
}

func (f *configFlag) String() string {
	return "path to config"
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conf_path := configFlag{""}
	flag.CommandLine.Var(&conf_path, "config", "Configuration file")
	flag.Parse()

	var config conf_reader.Conf
	config.GetConf(conf_path.ConfigPath)

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
