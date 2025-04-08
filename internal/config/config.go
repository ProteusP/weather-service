package config

import (
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	RedisAddr           string
	WeatherAPIKey       string
	CacheTTL            time.Duration
	ServerPort          string
	PopularCities       []string
	CacheWarmupInterval time.Duration
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
		RedisAddr:           os.Getenv("REDIS_ADDR"),
		WeatherAPIKey:       os.Getenv("WEATHERAPI_KEY"),
		CacheTTL:            parseDuration(os.Getenv("CACHE_TTL"), 30*time.Second),
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		PopularCities:       getEnvAsSlice("POPULAR_CITIES", []string{"Moscow", "London", "Paris", "Saint-Petersburg", "New-York"}, ","),
		CacheWarmupInterval: parseDuration(os.Getenv("CACHE_WARMUP_INTERVAL"), 30*time.Second),
	}
}

func getEnvAsSlice(key string, fallback []string, sep string) []string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	return strings.Split(val, sep)
}
