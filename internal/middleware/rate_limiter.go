package middleware

import (
	"context"
	"fmt"

	"net/http"
	"time"
	"weather-service/internal/cache"

	"github.com/gin-gonic/gin"
)

func RateLimiter(redisCache *cache.RedisCache, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Чтобы не считать запрос к "/favicon.ico"
		if c.Request.Method == http.MethodOptions || c.Request.URL.Path == "/favicon.ico" {
			c.Next()
			return
		}

		ip := c.ClientIP()
		key := fmt.Sprintf("rate:%s:%s", ip, c.Request.URL.Path)

		count, err := redisCache.GetClient().Incr(context.Background(), key).Result()

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if count == 1 {
			redisCache.GetClient().Expire(context.Background(), key, window)
		}

		if count > int64(limit) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Too many requests. Try again later."})
			return
		}

		c.Next()
	}
}
