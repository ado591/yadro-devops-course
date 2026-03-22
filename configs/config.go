package configs

import (
	"fmt"
	"os"
	"regexp"
)

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

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

var digitsOnly = regexp.MustCompile(`^\d+$`)

func LoadConfig() (Config, error) {
	port := getEnv("PORT", defaultPort)
	if !digitsOnly.MatchString(port) {
		return Config{}, fmt.Errorf("invalid PORT value %q: must contain digits only", port)
	}
	return Config{
		Version: getEnv("VERSION", defaultVersion),
		Author:  getEnv("AUTHOR", defaultAuthor),
		APIKey:  os.Getenv("API_KEY"),
		Port:    port,
	}, nil
}
