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
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

// RateDistRuleRaw rule for distruribted rate limit - for json parse
type RateDistRuleRaw struct {
	RedisEndpoint string `json:"redisEndpoint" validate:"required"`
	HeaderKey     string `json:"headerKey" validate:"required"`
	LookupSec     uint   `json:"lookupSec" validate:"required,gt=0"`
	Limit         int64  `json:"limit" validate:"required,gt=0"`
	IsRequiredKey bool   `json:"isRequiredKey"`
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
	// rateDistUpdateConfig()
	// rateDistUpdateConfigWorker()

	rateDistUpdateConfig("custommw/rateSample.json")
	rateDistUpdateConfig("custommw/rateSample.1.json")
	rateDistUpdateConfig("custommw/rateSample.json")
	// runtime.GC()
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
				// << FOR TEST >>
				// if key == "" {
				// // flase is default value in go
				// 	ch <- !limitRule.isRequiredKey
				// 	return
				// }
				// << FOR TEST >>

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
					// intead blockTime, can use TTL redis commnad result for acutally remain TTL at redis (maybe good to add +1 second)
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
					pipe.RPush(key, "_")
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

// Every minute check the map for visitors that haven't been seen for
// more than minutes and delete the entries.
func rateDistUpdateConfigWorker() {
	go func() {
		for {
			time.Sleep(time.Minute)
			rateDistUpdateConfig("custommw/rateSample.json")
		}
	}()
}

func rateDistUpdateConfig(filepath string) {
	// 1. watch limit config object - check changed date
	// local temp
	// << FOR TEST >>
	// _rateDistLimiterRuleRWMutex.RLock()
	// isAlreadyUpdate := len(_rateDistLimiterRules) > 0
	// _rateDistLimiterRuleRWMutex.RUnlock()
	// if isAlreadyUpdate {
	// 	return
	// }
	// << FOR TEST >>

	// 2. if updated, download & update date
	// local temp
	// Open our jsonFile
	jsonFile, err := os.Open(filepath)
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

	// map of current redis clients
	redisEndpointMap := map[string][]*redis.Client{}
	for _, ruleOld := range rulesOld {
		redisEndpointMap[ruleOld.redisEndpoint] = append(redisEndpointMap[ruleOld.redisEndpoint], ruleOld.redisCli)
	}

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

		ruleNew := rateDistRule{
			headerKey:     ruleRaw.HeaderKey,
			lookupTime:    time.Second * time.Duration(ruleRaw.LookupSec),
			limit:         ruleRaw.Limit,
			redisEndpoint: ruleRaw.RedisEndpoint,
			isRequiredKey: ruleRaw.IsRequiredKey,
			blockTime:     time.Second * time.Duration(ruleRaw.BlockSec),
		}

		// check possiable redis client
		redisCliReuse, exists := redisEndpointMap[ruleRaw.RedisEndpoint]
		if exists && len(redisCliReuse) > 0 {
			// reuse redis client
			ruleNew.redisCli = redisCliReuse[len(redisCliReuse)-1]
			redisEndpointMap[ruleRaw.RedisEndpoint] = redisEndpointMap[ruleRaw.RedisEndpoint][:len(redisCliReuse)-1]
			fmt.Println("hah! reuse")
		} else {
			// need new redis client
			ruleNew.redisCli = redis.NewClient(&redis.Options{
				Addr:         ruleRaw.RedisEndpoint,
				ReadTimeout:  time.Millisecond * time.Duration(redisTimeoutMil),
				WriteTimeout: time.Millisecond * time.Duration(redisTimeoutMil),
			})
			fmt.Println("hah! new client")
		}
		fmt.Println(ruleNew)

		rulesNew = append(rulesNew, ruleNew)
	}

	// 4. apply to config
	_rateDistLimiterRuleRWMutex.Lock()
	_rateDistLimiterRules = rulesNew
	_rateDistLimiterRuleRWMutex.Unlock()

	// 5. update updatedtime
	// actually, need set to object's last modefied time
	_rateDistConfigUpdatedDate = time.Now()

	// 6. release unused redis client resource
	fmt.Println("redisEndpointMap before", redisEndpointMap)
	for endpoint, redisClis := range redisEndpointMap {
		for _, redisCli := range redisClis {
			fmt.Println("redis client deleted for", endpoint)
			redisCli.Close()
		}
		redisEndpointMap[endpoint] = nil
	}
	fmt.Println("redisEndpointMap after", redisEndpointMap)

	fmt.Println("--------- config done ----------")

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
