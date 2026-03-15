package internal

type InfoResponse struct {
	Version string `json:"version"`
	Service string `json:"service"`
	Author  string `json:"author"`
}

type TemperatureStats struct {
	Average float64 `json:"average"`
	Median  float64 `json:"median"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
}

type WeatherData struct {
	TemperatureC TemperatureStats `json:"temperature_c"`
}

type WeatherResponse struct {
	Service string      `json:"service"`
	Data    WeatherData `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
