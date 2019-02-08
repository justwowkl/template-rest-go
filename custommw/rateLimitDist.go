package custommw

import (
	"sync"
	"time"
)

// // refer https://www.alexedwards.net/blog/how-to-rate-limit-http-requests

type limitRule struct {
	redisEndpoint string
	headerKey     string
	lookupMin     uint
	limit         uint
}

// // required global config - limit
var _rateDistLimiterRules []limitRule
var _rateDistLimiterRuleRWMutex *sync.RWMutex
var _rateDistConfigUpdatedDate time.Time

// InitRateDist init rate limiter
func InitRateDist() {
	_rateDistLimiterRuleRWMutex = new(sync.RWMutex)
	_rateLocalLimiterIPsRWMutex.Lock()
	_rateDistConfigUpdatedDate = time.Unix(0, 0).UTC()
	_rateLocalLimiterIPsRWMutex.Unlock()
	rateDistUpdateConfig()
	rateDistUpdateConfigWorker()
}

// func Unix(sec int64, nsec int64) Time

// // RateLimitLocal limit request by request IP
// func RateLimitLocal(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		ip := c.RealIP()
// 		if strings.HasPrefix(ip, "10.") {
// 			return next(c)
// 		}
// 		limiter := rateLocalGetIP(c.RealIP())
// 		if limiter.Allow() == false {
// 			return c.String(http.StatusTooManyRequests, "")
// 		}
// 		return next(c)
// 	}
// }

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
	// watch limit config object - changed date
	// if updated, download & update date
	// parse & apply to config

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
