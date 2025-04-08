package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"weather-service/internal/cache"
)

type WeatherService struct {
	cache  *cache.RedisCache
	apiKey string
}

func NewWeatherService(cache *cache.RedisCache, apiKey string) *WeatherService {
	return &WeatherService{
		cache:  cache,
		apiKey: apiKey,
	}
}

func (s *WeatherService) GetWeather(ctx context.Context, city string) (string, error) {

	cached, err := s.cache.Get(ctx, "weather:"+city)
	if err == nil {
		return cached, nil
	}

	url := fmt.Sprintf(
		"https://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=ru",
		s.apiKey, city,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Weather API request failed: %v", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Weather API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read response body: %v", err)
	}

	if err := s.cache.Set(ctx, "weather:"+city, string(body)); err != nil {
		fmt.Printf("Warning: failed to cache data: %v\n", err)
	}

	return string(body), err
}

func (s *WeatherService) GetForecast(ctx context.Context, city string, days int) (string, error) {
	strDays := string(days)
	cached, err := s.cache.Get(ctx, "forecast:"+city+":"+strDays)

	if err == nil {
		return cached, nil
	}

	url := fmt.Sprintf(
		"https://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=%d&lang=ru",
		s.apiKey,
		city,
		days)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Weather API request failed: %v", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Weather API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read response body: %v", err)
	}

	if err := s.cache.Set(ctx, "forecast:"+city+":"+strDays, string(body)); err != nil {
		fmt.Printf("Warning: failed to cache data: %v\n", err)
	}

	return string(body), nil
}
