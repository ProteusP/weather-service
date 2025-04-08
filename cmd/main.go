package main

import (
	"context"
	"fmt"
	"time"
	"weather-service/internal/cache"
	"weather-service/internal/config"
	"weather-service/internal/handler"
	"weather-service/internal/middleware"
	"weather-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.Load()

	redisCache, err := cache.NewRedisClient(cfg.RedisAddr, cfg.CacheTTL)
	if err != nil {
		panic(fmt.Sprintf("Failed to init Redis: %v", err))
	}

	//Можно добавить limit в config
	rateLimiter := middleware.RateLimiter(redisCache, 2, time.Hour)

	weatherService := service.NewWeatherService(redisCache, cfg.WeatherAPIKey)
	weatherHandler := handler.NewWeatherHandler(weatherService)
	forecastHandler := handler.NewForecastHandler(weatherService)

	r := gin.Default()
	r.Use(rateLimiter)

	r.GET("/weather/:city", weatherHandler.GetCurrentWeather)
	r.GET("/forecast/:city/:days", forecastHandler.GetForecast)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go weatherService.StartCacheWarmup(
		ctx,
		cfg.CacheWarmupInterval,
		cfg.PopularCities,
	)

	r.Run(":" + cfg.ServerPort)

}
