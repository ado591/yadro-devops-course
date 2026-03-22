package internal

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"weather/configs"
	"weather/internal/models"
)

type Weatherer interface {
	HistoryTemps(city string, from, to time.Time) ([]float64, error)
	ForecastTemps(city string, from, to time.Time) ([]float64, error)
}

type App struct {
	config configs.Config
	client Weatherer
}

func CreateApp(cfg configs.Config, client Weatherer) *App {
	return &App{config: cfg, client: client}
}

func writeJSON(responseWriter http.ResponseWriter, status int, data any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(status)
	_ = json.NewEncoder(responseWriter).Encode(data)
}

func (app *App) HandleInfo(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeJSON(responseWriter, http.StatusBadRequest, models.ErrorResponse{Error: "/info only supports GET"})
		return
	}
	writeJSON(responseWriter, http.StatusOK, models.InfoResponse{
		Version: app.config.Version,
		Service: configs.ServiceName,
		Author:  app.config.Author,
	})
}

func (a *App) HandleWeather(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "/weather only supports GET"})
		return
	}
	if a.config.APIKey == "" {
		writeJSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: "API_KEY not configured"})
		return
	}

	city := r.URL.Query().Get("city")
	if city == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "city parameter is required"})
		return
	}

	dateFromStr := r.URL.Query().Get("date_from")
	dateToStr := r.URL.Query().Get("date_to")

	today := time.Now().UTC().Truncate(24 * time.Hour)
	var temps []float64

	if dateFromStr == "" && dateToStr == "" {
		yesterday := today.AddDate(0, 0, -1)
		t, err := a.client.HistoryTemps(city, yesterday, today)
		if err != nil {
			log.Printf("HistoryTemps error: %v", err)
			writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		temps = t
	} else {
		dateFrom, dateTo, err := ParseDateRange(dateFromStr, dateToStr, today)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		if dateTo.After(today) {
			t, err := a.client.ForecastTemps(city, dateFrom, dateTo)
			if err != nil {
				log.Printf("ForecastTemps error: %v", err)
				writeJSON(w, http.StatusBadGateway, models.ErrorResponse{Error: err.Error()})
				return
			}
			temps = append(temps, t...)
		} else if !dateFrom.After(today) {
			t, err := a.client.HistoryTemps(city, dateFrom, dateTo)
			if err != nil {
				log.Printf("HistoryTemps error: %v", err)
				writeJSON(w, http.StatusBadGateway, models.ErrorResponse{Error: err.Error()})
				return
			}
			temps = append(temps, t...)
		}
	}

	writeJSON(w, http.StatusOK, models.WeatherResponse{
		Service: configs.ServiceName,
		Data:    models.WeatherData{TemperatureC: ComputeStats(temps)},
	})
}

func ParseDateRange(dateFromStr, dateToStr string, today time.Time) (time.Time, time.Time, error) {
	var dateFrom, dateTo time.Time
	var err error

	if dateFromStr != "" {
		dateFrom, err = time.Parse(DateLayout, dateFromStr)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("invalid format, expected YYYY-MM-DD")
		}
	}

	if dateToStr != "" {
		dateTo, err = time.Parse(DateLayout, dateToStr)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("invalid format, expected YYYY-MM-DD")
		}
	}

	switch {
	case dateFromStr == "" && dateToStr != "":
		dateFrom = dateTo
	case dateFromStr != "" && dateToStr == "":
		dateTo = dateFrom
	case dateFromStr == "" && dateToStr == "":
		dateFrom = today.AddDate(0, 0, -1)
		dateTo = today
	}

	if dateTo.Before(dateFrom) {
		return time.Time{}, time.Time{}, errors.New("date_to must be >= date_from")
	}

	return dateFrom, dateTo, nil
}
