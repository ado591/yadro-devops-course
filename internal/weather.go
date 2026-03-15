package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DateLayout     = "2006-01-02"
	weatherAPIBase = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline"
)

var defaultHTTPClient = &http.Client{Timeout: 10 * time.Second}

type apiHour struct {
	Temp float64 `json:"temp"`
}

type apiDay struct {
	Datetime string    `json:"datetime"`
	Hours    []apiHour `json:"hours"`
	Temp     float64   `json:"temp"`
}

type apiResponse struct {
	Days []apiDay `json:"days"`
}

type WeatherClient struct {
	APIKey     string
	HTTPClient *http.Client
}

func (c *WeatherClient) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return defaultHTTPClient
}

func (c *WeatherClient) get(path, include string) (*apiResponse, error) {
	params := url.Values{
		"key":         {c.APIKey},
		"unitGroup":   {"metric"},
		"include":     {include},
		"contentType": {"json"},
	}
	apiURL := fmt.Sprintf("%s/%s?%s", weatherAPIBase, path, params.Encode())

	resp, err := c.httpClient().Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, body)
	}

	var r apiResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}
	return &r, nil
}

func collectTemps(days []apiDay) []float64 {
	var temps []float64
	for _, d := range days {
		if len(d.Hours) > 0 {
			for _, h := range d.Hours {
				temps = append(temps, h.Temp)
			}
		} else {
			temps = append(temps, d.Temp)
		}
	}
	return temps
}

func (c *WeatherClient) CurrentTemp(location string) (float64, error) {
	path := url.PathEscape(location) + "/today"
	r, err := c.get(path, "hours")
	if err != nil {
		return 0, err
	}
	if len(r.Days) == 0 {
		return 0, fmt.Errorf("no data for %s", location)
	}
	hours := r.Days[0].Hours
	if len(hours) == 0 {
		return r.Days[0].Temp, nil
	}
	return hours[len(hours)-1].Temp, nil
}

func (c *WeatherClient) HistoryTemps(location string, from, to time.Time) ([]float64, error) {
	path := fmt.Sprintf("%s/%s/%s",
		url.PathEscape(location),
		from.Format(DateLayout),
		to.Format(DateLayout),
	)
	r, err := c.get(path, "hours")
	if err != nil {
		return nil, err
	}
	return collectTemps(r.Days), nil
}

// Вообще прогноз погоды ничем от истории не отличается, но для расширяемости вынесла в отдельный метод
func (c *WeatherClient) ForecastTemps(location string, from, to time.Time) ([]float64, error) {
	return c.HistoryTemps(location, from, to)
}
