package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

func ConnectToRedis(url string) (*redis.Client, error) {
	fmt.Println(url)
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Printf("Failed to connect to redis cache:%s\n", err.Error())
		return nil, err
	}
	return redis.NewClient(opts), nil
}

type BaseCache struct {
	*redis.Client
}

func NewBaseCache(client *redis.Client) *BaseCache {
	return &BaseCache{Client: client}
}

func (b *BaseCache) SetObject(k string, v interface{}, exp int64) error {
	bytes, err := msgpack.Marshal(v)
	if err != nil {
		panic(err)
	}
	res := b.Set(context.Background(), k, bytes, time.Duration(exp)*time.Second)
	return res.Err()
}

func (b *BaseCache) GetObject(k string, v interface{}) error {
	s, err := b.Get(context.Background(), k).Bytes()
	if err != nil {
		return err
	}

	return msgpack.Unmarshal(s, v)
}
