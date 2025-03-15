package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/vanh01/api-rate-limiting/pkg/cache"
)

func FixedWindowRateLimit(limit, windowSec int64, bcache *cache.BaseCache) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Request().Header.Get("userId")
			userID, err := uuid.Parse(id)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			currentWindow := time.Now().Unix() / windowSec
			key := fmt.Sprintf("fixed_rate_limit:%s:%d", userID, currentWindow)

			count, err := bcache.Incr(context.Background(), key).Result()
			if err != nil {
				log.Println("Redis error:", err)
				return c.String(http.StatusInternalServerError, err.Error())
			}

			if count == 1 {
				bcache.Expire(context.Background(), key, time.Duration(windowSec)*time.Second)
			}

			if count > limit {
				return c.String(http.StatusForbidden, "You have reached your access limit. Action not allowed.")
			}

			return next(c)
		}
	}
}
