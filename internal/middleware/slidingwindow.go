package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"

	"github.com/vanh01/api-rate-limiting/pkg/cache"
)

func SlidingWindowRateLimit(limit, windowSec int64, bcache *cache.BaseCache) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Request().Header.Get("userId")
			userID, err := uuid.Parse(id)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			key := fmt.Sprintf("sliding_rate_limit:%s", userID)
			now := float64(time.Now().UnixNano()) / 1e9 // Lấy timestamp hiện tại theo giây

			bcache.ZRemRangeByScore(context.Background(), key, "0", fmt.Sprintf("%f", now-float64(windowSec)))

			count, err := bcache.ZCard(context.Background(), key).Result()
			if err != nil {
				log.Println("Redis error:", err)
				return c.String(http.StatusInternalServerError, err.Error())
			}

			if count >= limit {
				return c.String(http.StatusForbidden, "You have reached your access limit. Action not allowed.")
			}

			bcache.ZAdd(context.Background(), key, redis.Z{
				Score:  now,
				Member: now,
			})

			bcache.Expire(context.Background(), key, time.Duration(windowSec)*time.Second)

			return next(c)
		}
	}
}
