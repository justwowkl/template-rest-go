package util

import (
	"fmt"

	"github.com/bluele/gcache"
)

var _gcache gcache.Cache

// cacheInit init cache
func gcacheInit() {
	_gcache = gcache.New(20). // num of cache object
					ARC().
					Build()
}

func gcacheHealth() bool {
	_gcache.Set(true, true)
	val, err := _gcache.Get(true)
	if err == gcache.KeyNotFoundError {
		fmt.Println("ok, not found")
		return false
	} else if err != nil {
		fmt.Println("oh, real error")
		return false
	} else if val != true {
		fmt.Println("wrong value")
		return false
	}
	fmt.Println("Get:", val)
	_gcache.Remove(true)
	return true
}

func gcacheGet(key string) (interface{}, error) {
	val, err := _gcache.Get(key)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func gcacheSet(key string, val interface{}) {
	_gcache.Set(key, val)
}
