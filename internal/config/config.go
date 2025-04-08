package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	RedisAddr     string
	WeatherAPIKey string
	CacheTTL      time.Duration
	ServerPort    string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func parseDuration(envVar string, defaultVal time.Duration) time.Duration {
	if envVar == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(envVar)
	if err != nil {
		log.Printf("Time parsing error: %s", err)
		return defaultVal
	}
	return d
}
func Load() *Config {
	return &Config{
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		WeatherAPIKey: os.Getenv("WEATHERAPI_KEY"),
		CacheTTL:      parseDuration(os.Getenv("CACHE_TTL"), 30*time.Second),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
	}
}
