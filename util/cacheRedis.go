package util

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var _redis *redis.Client

// redisInit init redis
func redisInit() {
	_redis = redis.NewClient(&redis.Options{
		// https://godoc.org/github.com/go-redis/redis#Options
		Addr:         "localhost:6379",
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	})
	fmt.Println("redis init done")
}

// redisHealth healthcheck
func redisHealth() bool {
	_, err := _redis.Ping().Result()
	if err != nil {
		fmt.Println("redis not respond")
		return false
	}
	fmt.Println("redis health okay")
	return true
}

// redisGet get from redis
func redisGet(key string) (interface{}, error) {
	val, err := _redis.Get(key).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

// redisSet set redis
func redisSet(key string, val interface{}) {
	err := _redis.Set(key, val, 0).Err()
	if err != nil {
		fmt.Println("fail to write")
	}
}
