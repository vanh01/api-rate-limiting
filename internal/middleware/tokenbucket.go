package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TokenBucket struct {
	Capacity int
	FillRate int
	NoTokens int
	LastTime time.Time
	Lock     sync.Mutex
}

func TokenBucketRateLimit(capacity, fillRate, token int) echo.MiddlewareFunc {
	clients := make(map[string]*TokenBucket)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Request().Header.Get("userId")
			userID, err := uuid.Parse(id)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			key := fmt.Sprintf("token_rate_limit:%s", userID)
			if _, found := clients[key]; !found {
				clients[key] = &TokenBucket{
					Capacity: capacity,
					FillRate: fillRate,
					NoTokens: capacity,
					LastTime: time.Now(),
				}

				return next(c)
			}

			client := clients[key]
			now := time.Now()
			timePassed := now.Sub(client.LastTime).Seconds()
			client.NoTokens = min(capacity, client.NoTokens+int(timePassed)*client.FillRate)

			if client.NoTokens < token {
				return c.String(http.StatusForbidden, "You have reached your access limit. Action not allowed.")
			}

			client.LastTime = now
			client.NoTokens -= token

			return next(c)
		}
	}
}
