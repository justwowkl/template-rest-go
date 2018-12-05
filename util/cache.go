package util

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var _redisClient *redis.Client

// cacheInit init cache
func cacheInit() {
	_redisClient = redis.NewClient(&redis.Options{
		// https://godoc.org/github.com/go-redis/redis#Options
		Addr:         "localhost:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB
		PoolSize:     20, // default is 10, per CPU
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	})
	fmt.Println("cache init done")
}

// CacheHealth healthcheck
func CacheHealth() bool {
	_, err := _redisClient.Ping().Result()
	if err != nil {
		fmt.Println("cache not respond")
		return false
	}
	return true
}

// CacheGet get from cache
func CacheGet(key string) {
	val, err := _redisClient.Get(key).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		panic(err) // no need panic
	} else {
		fmt.Println("key", val)
	}
}

// CacheSet set cache
func CacheSet(key string, value interface{}) {
	err := _redisClient.Set(key, value, 0).Err()
	if err != nil {
		panic(err) // no need panic
	}
}
