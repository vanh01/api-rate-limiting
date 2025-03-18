package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type LeakyBucket struct {
	Capacity int
	LeakRate int
	Bucket   []time.Time
	LastLeak time.Time
	Lock     sync.Mutex
}

func LeakyBucketRateLimit(capacity, leakRate int) echo.MiddlewareFunc {
	clients := make(map[string]*LeakyBucket)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Request().Header.Get("userId")
			userID, err := uuid.Parse(id)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			key := fmt.Sprintf("leaky_rate_limit:%s", userID)
			if _, found := clients[key]; !found {
				clients[key] = &LeakyBucket{
					Capacity: capacity,
					LeakRate: leakRate,
					Bucket:   []time.Time{},
					LastLeak: time.Now(),
				}

				return next(c)
			}

			client := clients[key]
			client.Lock.Lock()
			defer client.Lock.Unlock()
			now := time.Now()
			timePassed := now.Sub(client.LastLeak).Seconds()
			leaked := int(timePassed) * leakRate

			if leaked > 0 {
				for range min(leaked, len(client.Bucket)) {
					client.Bucket = client.Bucket[1:]
				}

				client.LastLeak = now
			}

			if len(client.Bucket) > client.Capacity {
				return c.String(http.StatusForbidden, "You have reached your access limit. Action not allowed.")
			}

			client.Bucket = append(client.Bucket, now)
			return next(c)
		}
	}
}
