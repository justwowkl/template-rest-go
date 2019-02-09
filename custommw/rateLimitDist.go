package custommw

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/bluele/gcache"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

// RateDistRuleRaw rule for distruribted rate limit - for json parse
type RateDistRuleRaw struct {
	RedisEndpoint string `json:"redisEndpoint" validate:"required"`
	HeaderKey     string `json:"headerKey" validate:"required"`
	LookupSec     uint   `json:"lookupSec" validate:"required,gt=0"`
	Limit         int64  `json:"limit" validate:"required,gt=0"`
	IsRequiredKey bool   `json:"isRequired" validate:"required"`
	BlockSec      uint   `json:"blockTimeSec" validate:"required,gt=0"`
}

// RedisEndpoint rule for distruribted rate limit - internal
type rateDistRule struct {
	headerKey     string
	lookupTime    time.Duration
	limit         int64
	redisEndpoint string
	redisCli      *redis.Client
	isRequiredKey bool
	blockTime     time.Duration
}

// required global config - limit
var redisTimeoutMil = 100
var sizeLocalcache = 20

var _rateDistLimiterRules []rateDistRule
var _rateDistLimiterRuleRWMutex *sync.RWMutex
var _rateDistConfigUpdatedDate time.Time
var _rateDistcacheCli gcache.Cache

// InitRateDist init rate limiter
func InitRateDist() {
	_rateDistLimiterRuleRWMutex = new(sync.RWMutex)
	_rateDistConfigUpdatedDate = time.Unix(0, 0) // need check UTC required
	_rateDistcacheCli = gcache.New(sizeLocalcache).
		ARC().
		Build()
	rateDistUpdateConfig()
	// rateDistUpdateConfigWorker()
}

// RateLimitDist limit request by request IP
func RateLimitDist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		_rateDistLimiterRuleRWMutex.RLock()
		limitRules := _rateDistLimiterRules
		_rateDistLimiterRuleRWMutex.RUnlock()

		if len(limitRules) < 1 {
			return next(c)
		}

		ch := make(chan bool, len(limitRules))
		for _, _limitRule := range limitRules {
			go func(limitRule rateDistRule) {
				// check for each rules with redis
				// https://redis.io/commands/INCR

				// < psudo >
				// current = LLEN(ip)
				// IF current > 10 THEN
				// 	ERROR "too many requests per second"
				// ELSE
				// 	IF EXISTS(ip) == FALSE
				// 		MULTI
				// 			RPUSH(ip,ip)
				// 			EXPIRE(ip,1)
				// 		EXEC
				// 	ELSE
				// 		RPUSHX(ip,ip)
				// 	END
				// 	PERFORM_API_CALL()
				// END

				// get header value & check exists
				key := c.Request().Header.Get(limitRule.headerKey)
				if key == "" {
					ch <- !limitRule.isRequiredKey
					return
				}

				// check local blacklist cache
				keyCache := fmt.Sprintf("%s %s", limitRule.headerKey, key)
				_, err := _rateDistcacheCli.Get(keyCache)
				if err == nil {
					ch <- false
					return
				}

				// get from redis & check
				rateCur, err := limitRule.redisCli.LLen(key).Result()
				if err != nil {
					ch <- true
					return
				}

				// if exceeds, add blacklist and set fail
				if rateCur > limitRule.limit {
					_rateDistcacheCli.SetWithExpire(keyCache, true, limitRule.blockTime)
					ch <- false
					return
				}

				// check need create key in redis
				isKeyExists, err := limitRule.redisCli.Exists(key).Result()
				if err != nil {
					ch <- true
					return
				}
				if isKeyExists == 0 {
					// https://godoc.org/github.com/go-redis/redis#ex-Client-TxPipeline
					pipe := limitRule.redisCli.TxPipeline()
					pipe.RPush(key, " ")
					pipe.Expire(key, limitRule.lookupTime)
					_, err := pipe.Exec()
					if err != nil {
						fmt.Println(err) // just mark..
					}
				} else {
					err := limitRule.redisCli.RPushX(key, " ").Err()
					if err != nil {
						fmt.Println(err) // just mark..
					}
				}
				ch <- true

			}(_limitRule)
		}

		for i := 0; i < len(limitRules); i++ {
			if !<-ch {
				return c.String(http.StatusTooManyRequests, "")
			}
		}
		return next(c)
	}
}

// // Create a new rate limiter and add it to the visitors map, using the
// // IP address as the key.
// func rateLocalAddIP(ip string) *rate.Limiter {
// 	limiter := rate.NewLimiter(2, 2)
// 	_rateLocalLimiterIPsRWMutex.Lock()
// 	_rateLocalLimiterIPs[ip] = &typeLimitIP{
// 		limiter:  limiter,
// 		lastSeen: time.Now(),
// 	}
// 	println("new limiter : ", _rateLocalLimiterIPs[ip])
// 	_rateLocalLimiterIPsRWMutex.Unlock()
// 	return limiter
// }

// // Retrieve and return the rate limiter for the current visitor if it
// // already exists. Otherwise call the addVisitor function to add a
// // new entry to the map.
// func rateLocalGetIP(ip string) *rate.Limiter {
// 	_rateLocalLimiterIPsRWMutex.RLock()
// 	limiterIP, exists := _rateLocalLimiterIPs[ip]
// 	_rateLocalLimiterIPsRWMutex.RUnlock()
// 	if !exists {
// 		return rateLocalAddIP(ip)
// 	}
// 	println("old limiter : ", limiterIP)
// 	return limiterIP.limiter
// }

// Every minute check the map for visitors that haven't been seen for
// more than minutes and delete the entries.
func rateDistUpdateConfigWorker() {
	go func() {
		for {
			time.Sleep(time.Minute)
			rateDistUpdateConfig()
		}
	}()
}

func rateDistUpdateConfig() {
	// 1. watch limit config object - check changed date
	// local temp
	_rateDistLimiterRuleRWMutex.RLock()
	isAlreadyUpdate := len(_rateDistLimiterRules) > 0
	_rateDistLimiterRuleRWMutex.RUnlock()
	if isAlreadyUpdate {
		return
	}

	// 2. if updated, download & update date
	// local temp
	// Open our jsonFile
	jsonFile, err := os.Open("custommw/rateSample.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// 3. parse json & validate
	// local temp
	jsonByte, _ := ioutil.ReadAll(jsonFile)
	var rulesRaw []RateDistRuleRaw
	json.Unmarshal(jsonByte, &rulesRaw)

	_rateDistLimiterRuleRWMutex.RLock()
	rulesOld := _rateDistLimiterRules
	_rateDistLimiterRuleRWMutex.RUnlock()
	rulesNew := make([]rateDistRule, 0, len(rulesRaw))

	// create new rules object
	jsonValidate := validator.New()
	for _, ruleRaw := range rulesRaw {
		// validate
		err := jsonValidate.Struct(ruleRaw)
		if err != nil {
			// json error..!
			println(err.Error())
			return
		}

		// check new header key
		var reuseRedisCli *redis.Client
		for _, ruleOld := range rulesOld {
			if ruleRaw.HeaderKey == ruleOld.headerKey &&
				ruleRaw.RedisEndpoint == ruleOld.redisEndpoint &&
				ruleRaw.Limit == ruleOld.limit &&
				time.Second*time.Duration(ruleRaw.LookupSec) == ruleOld.lookupTime {
				reuseRedisCli = ruleOld.redisCli
				break
			}
		}
		ruleNew := rateDistRule{
			headerKey:     ruleRaw.HeaderKey,
			lookupTime:    time.Second * time.Duration(ruleRaw.LookupSec),
			limit:         ruleRaw.Limit,
			redisEndpoint: ruleRaw.RedisEndpoint,
			isRequiredKey: ruleRaw.IsRequiredKey,
			blockTime:     time.Second * time.Duration(ruleRaw.BlockSec),
		}
		if reuseRedisCli == nil {
			ruleNew.redisCli = redis.NewClient(&redis.Options{
				Addr:         ruleRaw.RedisEndpoint,
				ReadTimeout:  time.Millisecond * time.Duration(redisTimeoutMil),
				WriteTimeout: time.Millisecond * time.Duration(redisTimeoutMil),
			})
		} else {
			ruleNew.redisCli = reuseRedisCli
		}

		rulesNew = append(rulesNew, ruleNew)
	}

	// find redis client to release
	redisCliUnused := make([]*redis.Client, 0, len(rulesRaw))
	for _, ruleOld := range rulesOld {
		isUnused := true
		for _, ruleNew := range rulesNew {
			if ruleNew.headerKey == ruleOld.headerKey &&
				ruleNew.redisEndpoint == ruleOld.redisEndpoint &&
				ruleNew.limit == ruleOld.limit &&
				ruleNew.lookupTime == ruleOld.lookupTime {
				isUnused = false
				break
			}
		}
		if isUnused {
			redisCliUnused = append(redisCliUnused, ruleOld.redisCli)
		}
	}

	// 4. apply to config
	_rateDistLimiterRuleRWMutex.Lock()
	_rateDistLimiterRules = rulesNew
	_rateDistLimiterRuleRWMutex.Unlock()

	// 5. update updatedtime
	// actually, need set to object's last modefied time
	_rateDistConfigUpdatedDate = time.Now()

	// 6. release unused redis client resource
	for _, redisCli := range redisCliUnused {
		redisCli.Close()
	}

	// _rateLocalLimiterIPsRWMutex.RLock()
	// for ip, limiterIP := range _rateLocalLimiterIPs {
	// 	if time.Now().Sub(limiterIP.lastSeen) > 5*time.Minute {
	// 		_rateLocalLimiterIPsRWMutex.RUnlock()
	// 		_rateLocalLimiterIPsRWMutex.Lock()
	// 		delete(_rateLocalLimiterIPs, ip)
	// 		_rateLocalLimiterIPsRWMutex.Unlock()
	// 		_rateLocalLimiterIPsRWMutex.RLock()
	// 	}
	// }
	// _rateLocalLimiterIPsRWMutex.RUnlock()
}
