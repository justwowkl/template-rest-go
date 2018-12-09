package util

import (
	"fmt"
	"sync"
)

// CacheGet try to get from cache
func CacheGet(key string) (interface{}, error) {
	val, err := gcacheGet(key)
	if err == nil {
		fmt.Println("[cache] hit!! gcache : ", key, " / ", val)
		return val, nil
	}
	val, err = redisGet(key)
	if err == nil {
		fmt.Println("[cache] hit!! redis : ", key, " / ", val)
		gcacheSet(key, val)
		return val, nil
	}
	return nil, err
}

// CacheSet try to set cache
func CacheSet(key string, val interface{}) {
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		defer wait.Done()
		redisSet(key, val)
		fmt.Println("[cache] set redis")
	}()
	gcacheSet(key, val)
	fmt.Println("[cache] set gcache")
	wait.Wait()
}
