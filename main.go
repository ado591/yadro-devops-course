package main

import (
	"fmt"
	"net/http"
	"os"

	"weather/configs"
	"weather/internal"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	app := internal.CreateApp(cfg, &internal.WeatherClient{APIKey: cfg.APIKey})

	mux := http.NewServeMux()
	mux.HandleFunc("/info", app.HandleInfo)
	mux.HandleFunc("/info/weather", app.HandleWeather)

	addr := fmt.Sprintf(":%s", cfg.Port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		os.Exit(1)
	}
}
