package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"

	"github.com/vanh01/api-rate-limiting/pkg/cache"
)

type TokenBucket struct {
	Capacity int
	FillRate int
	NoTokens int
	LastTime time.Time
}

func TokenBucketRateLimit(capacity, fillRate, token int, bcache *cache.BaseCache) echo.MiddlewareFunc {
	clients := make(map[string]*sync.Mutex)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Request().Header.Get("userId")
			userID, err := uuid.Parse(id)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			key := fmt.Sprintf("token_rate_limit:%s", userID)
			if _, exist := clients[key]; !exist {
				clients[key] = &sync.Mutex{}
			}

			clients[key].Lock()
			defer clients[key].Unlock()

			var client *TokenBucket
			err = bcache.GetObject(key, &client)
			if err == redis.Nil {
				client := &TokenBucket{
					Capacity: capacity,
					FillRate: fillRate,
					NoTokens: capacity,
					LastTime: time.Now(),
				}
				bcache.SetObject(key, client, int64(capacity/fillRate))

				return next(c)
			}

			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			now := time.Now()
			timePassed := now.Sub(client.LastTime).Seconds()
			client.NoTokens = min(capacity, client.NoTokens+int(timePassed)*client.FillRate)

			if client.NoTokens < token {
				return c.String(http.StatusForbidden, "You have reached your access limit. Action not allowed.")
			}

			client.LastTime = now
			client.NoTokens -= token
			bcache.SetObject(key, client, int64(capacity/fillRate))

			return next(c)
		}
	}
}
