package handler

import (
	"encoding/json"
	"net/http"
	"weather-service/internal/service"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	service *service.WeatherService
}

func NewWeatherHandler(service *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		service: service,
	}
}

func (h *WeatherHandler) GetCurrentWeather(c *gin.Context) {
	city := c.Param("city")

	data, err := h.service.GetWeather(c, city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var jsonData map[string]any

	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to parse weather data"})
		return
	}

	c.JSON(http.StatusOK, jsonData)
}
