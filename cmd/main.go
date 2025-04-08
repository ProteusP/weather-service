package main

import (
	"fmt"
	"weather-service/internal/cache"
	"weather-service/internal/config"
	"weather-service/internal/handler"
	"weather-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.Load()

	redisCache, err := cache.NewRedisClient(cfg.RedisAddr, cfg.CacheTTL)
	if err != nil {
		panic(fmt.Sprintf("Failed to init Redis: %v", err))
	}

	weatherService := service.NewWeatherService(redisCache, cfg.WeatherAPIKey)

	weatherHandler := handler.NewWeatherHandler(weatherService)
	forecastHandler := handler.NewForecastHandler(weatherService)

	r := gin.Default()
	r.GET("/weather/:city", weatherHandler.GetCurrentWeather)
	r.GET("/forecast/:city/:days", forecastHandler.GetForecast)

	r.Run(":" + cfg.ServerPort)

}
