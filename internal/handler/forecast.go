package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"weather-service/internal/service"

	"github.com/gin-gonic/gin"
)

type ForecastHandler struct {
	service *service.WeatherService
}

func NewForecastHandler(service *service.WeatherService) *ForecastHandler {
	return &ForecastHandler{
		service: service,
	}
}

func (h *ForecastHandler) GetForecast(c *gin.Context) {
	city := c.Param("city")
	days, err := strconv.Atoi(c.Param("days"))
	if err != nil {
		fmt.Printf("Unable to convert days to Int: %s", c.Param("days"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Days param should be `Int`"})
	}

	data, err := h.service.GetForecast(c, city, days)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var jsonData map[string]any

	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to parse forecast data"})
		return
	}

	c.JSON(http.StatusOK, jsonData)
}
