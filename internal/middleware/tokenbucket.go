package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TokenBucket struct {
	Capacity int
	FillRate int
	NoTokens int
	LastTime time.Time
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

			now := time.Now()
			timePassed := now.Sub(clients[key].LastTime).Seconds()
			clients[key].NoTokens = min(capacity, clients[key].NoTokens+int(timePassed)*clients[key].FillRate)

			if clients[key].NoTokens < token {
				return c.String(http.StatusForbidden, "You have reached your access limit. Action not allowed.")
			}

			clients[key].LastTime = now
			clients[key].NoTokens -= token

			return next(c)
		}
	}
}
