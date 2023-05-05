package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var Cache *redis.Client
var Cachechannel chan string

func SetupRedis() {
	Cache = redis.NewClient(&redis.Options{
		Addr: "redis:2000",
		DB:   0,
	})
}

func SetupCacheChannel() {
	Cachechannel = make(chan string)

	go func(ch chan string) {
		for {
			time.Sleep(5 * time.Second)
			value := <-ch
			Cache.Del(context.Background(), value)
			fmt.Printf("cache cleared, %v", value)
		}
	}(Cachechannel)
}

func ClearCache(keys ...string) {
	for _, key := range keys {
		Cachechannel <- key
	}
}
