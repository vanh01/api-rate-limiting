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

type LeakyBucket struct {
	Capacity int
	LeakRate int
	Bucket   []time.Time
	LastLeak time.Time
}

func LeakyBucketRateLimit(capacity, leakRate int, bcache *cache.BaseCache) echo.MiddlewareFunc {
	clients := make(map[string]*sync.Mutex)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Request().Header.Get("userId")
			userID, err := uuid.Parse(id)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			key := fmt.Sprintf("leaky_rate_limit:%s", userID)
			if _, exist := clients[key]; !exist {
				clients[key] = &sync.Mutex{}
			}

			clients[key].Lock()
			defer clients[key].Unlock()

			var client *LeakyBucket
			err = bcache.GetObject(key, &client)
			if err == redis.Nil {
				client := &LeakyBucket{
					Capacity: capacity,
					LeakRate: leakRate,
					Bucket:   []time.Time{},
					LastLeak: time.Now(),
				}
				bcache.SetObject(key, client, int64(capacity/leakRate))

				return next(c)
			}

			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

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
			bcache.SetObject(key, client, int64(capacity/leakRate))

			return next(c)
		}
	}
}
