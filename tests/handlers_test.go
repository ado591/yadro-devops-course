package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"weather/configs"
	"weather/internal"
	"weather/internal/models"
)

func TestParseDateRange_BothEmpty(t *testing.T) {
	today := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	from, to, err := internal.ParseDateRange("", "", today)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !from.Equal(today.AddDate(0, 0, -1)) {
		t.Errorf("expected from=yesterday, got %v", from)
	}
	if !to.Equal(today) {
		t.Errorf("expected to=today, got %v", to)
	}
}

func TestParseDateRange_OnlyFrom(t *testing.T) {
	today := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	from, to, err := internal.ParseDateRange("2024-06-05", "", today)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !from.Equal(to) {
		t.Errorf("expected from==to when only date_from provided, got from=%v to=%v", from, to)
	}
}

func TestParseDateRange_OnlyTo(t *testing.T) {
	today := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	from, to, err := internal.ParseDateRange("", "2024-06-08", today)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !from.Equal(to) {
		t.Errorf("expected from==to when only date_to provided, got from=%v to=%v", from, to)
	}
}

func TestParseDateRange_ToBeforeFrom(t *testing.T) {
	today := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	_, _, err := internal.ParseDateRange("2024-06-10", "2024-06-01", today)
	if err == nil {
		t.Error("expected error when date_to < date_from")
	}
}

func TestHandleInfo_GET(t *testing.T) {
	cfg := configs.Config{Version: "1000.7.993", Author: "k.kaneki"}
	app := internal.CreateApp(cfg, nil)

	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	rr := httptest.NewRecorder()
	app.HandleInfo(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp models.InfoResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Version != "1000.7.993" {
		t.Errorf("expected version=1000.7.993, got %q", resp.Version)
	}
	if resp.Author != "k.kaneki" {
		t.Errorf("expected author=k.kaneki, got %q", resp.Author)
	}
	if resp.Service != configs.ServiceName {
		t.Errorf("expected service=%q, got %q", configs.ServiceName, resp.Service)
	}
}

func TestHandleInfo_POST(t *testing.T) {
	app := internal.CreateApp(configs.Config{}, nil)
	req := httptest.NewRequest(http.MethodPost, "/info", nil)
	rr := httptest.NewRecorder()
	app.HandleInfo(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

type mockWeatherer struct {
	historyTemps  []float64
	forecastTemps []float64
	err           error
}

func (m *mockWeatherer) HistoryTemps(_ string, _, _ time.Time) ([]float64, error) {
	return m.historyTemps, m.err
}

func (m *mockWeatherer) ForecastTemps(_ string, _, _ time.Time) ([]float64, error) {
	return m.forecastTemps, m.err
}

func TestHandleWeather_MissingCity(t *testing.T) {
	app := internal.CreateApp(configs.Config{APIKey: "key"}, &mockWeatherer{})
	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	rr := httptest.NewRecorder()
	app.HandleWeather(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestHandleWeather_NoAPIKey(t *testing.T) {
	app := internal.CreateApp(configs.Config{APIKey: ""}, &mockWeatherer{})
	req := httptest.NewRequest(http.MethodGet, "/weather?city=Moscow", nil)
	rr := httptest.NewRecorder()
	app.HandleWeather(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestHandleWeather_DefaultRange(t *testing.T) {
	mock := &mockWeatherer{historyTemps: []float64{10, 20, 30}}
	app := internal.CreateApp(configs.Config{APIKey: "key"}, mock)

	req := httptest.NewRequest(http.MethodGet, "/weather?city=Moscow", nil)
	rr := httptest.NewRecorder()
	app.HandleWeather(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp models.WeatherResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Data.TemperatureC.Min != 10 || resp.Data.TemperatureC.Max != 30 {
		t.Errorf("unexpected stats: %+v", resp.Data.TemperatureC)
	}
}

func TestHandleWeather_InvalidDateFormat(t *testing.T) {
	app := internal.CreateApp(configs.Config{APIKey: "key"}, &mockWeatherer{})
	req := httptest.NewRequest(http.MethodGet, "/weather?city=Moscow&date_from=bad-date", nil)
	rr := httptest.NewRecorder()
	app.HandleWeather(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestHandleWeather_POST(t *testing.T) {
	app := internal.CreateApp(configs.Config{APIKey: "key"}, &mockWeatherer{})
	req := httptest.NewRequest(http.MethodPost, "/weather?city=Moscow", nil)
	rr := httptest.NewRecorder()
	app.HandleWeather(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}
