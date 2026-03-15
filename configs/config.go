package configs

import "os"

const (
	ServiceName    = "weather"
	defaultVersion = "1.0.0"
	defaultAuthor  = "a.obraztsova"
	defaultPort    = "8000"
)

type Config struct {
	Version string
	Author  string
	APIKey  string
	Port    string
}

// обертка с дефолтными значениями
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func LoadConfig() Config {
	return Config{
		Version: getEnv("VERSION", defaultVersion),
		Author:  getEnv("AUTHOR", defaultAuthor),
		APIKey:  os.Getenv("API_KEY"),
		Port:    getEnv("PORT", defaultPort),
	}
}
