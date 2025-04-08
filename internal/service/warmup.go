package service

import (
	"context"
	"log"
	"time"
)

func (s *WeatherService) StartCacheWarmup(ctx context.Context, interval time.Duration, cities []string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.warmupCacheForCities(ctx, cities)
		}

	}

}

func (s *WeatherService) warmupCacheForCities(ctx context.Context, cities []string) {

	for _, city := range cities {

		select {
		case <-ctx.Done():
			return
		default:
			data, err := s.GetWeather(ctx, city)
			if err != nil {
				log.Printf("Failed to warmup cache for %s: %v", city, err)
				continue
			}
			log.Printf("Cache warmed up for %s: %d bytes", city, len(data))
		}

	}

}
