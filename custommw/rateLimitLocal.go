package custommw

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/time/rate"
)

// refer https://www.alexedwards.net/blog/how-to-rate-limit-http-requests

type typeLimitIP struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// required global config - limit
var _rateLocalLimiterIPs map[string]*typeLimitIP
var _rateLocalLimiterIPsRWMutex *sync.RWMutex

// InitRateLocal init rate limiter
func InitRateLocal() {
	_rateLocalLimiterIPs = make(map[string]*typeLimitIP)
	_rateLocalLimiterIPsRWMutex = new(sync.RWMutex)
	rateLocalCleanupIPsWorker()
}

// RateLimitLocal limit request by request IP
func RateLimitLocal(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()
		if strings.HasPrefix(ip, "10.") {
			return next(c)
		}
		limiter := rateLocalGetIP(c.RealIP())
		if limiter.Allow() == false {
			return c.String(http.StatusTooManyRequests, "")
		}
		return next(c)
	}
}

// Create a new rate limiter and add it to the visitors map, using the
// IP address as the key.
func rateLocalAddIP(ip string) *rate.Limiter {
	limiter := rate.NewLimiter(2, 2)
	_rateLocalLimiterIPsRWMutex.Lock()
	_rateLocalLimiterIPs[ip] = &typeLimitIP{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	_rateLocalLimiterIPsRWMutex.Unlock()
	println("new limiter : ", _rateLocalLimiterIPs[ip])
	return limiter
}

// Retrieve and return the rate limiter for the current visitor if it
// already exists. Otherwise call the addVisitor function to add a
// new entry to the map.
func rateLocalGetIP(ip string) *rate.Limiter {
	_rateLocalLimiterIPsRWMutex.RLock()
	limiterIP, exists := _rateLocalLimiterIPs[ip]
	_rateLocalLimiterIPsRWMutex.RUnlock()
	if !exists {
		return rateLocalAddIP(ip)
	}
	println("old limiter : ", limiterIP)
	return limiterIP.limiter
}

// Every minute check the map for visitors that haven't been seen for
// more than 5 minutes and delete the entries.
func rateLocalCleanupIPsWorker() {
	go func() {
		for {
			time.Sleep(time.Minute)
			_rateLocalLimiterIPsRWMutex.RLock()
			for ip, limiterIP := range _rateLocalLimiterIPs {
				if time.Now().Sub(limiterIP.lastSeen) > 5*time.Minute {
					_rateLocalLimiterIPsRWMutex.RUnlock()
					_rateLocalLimiterIPsRWMutex.Lock()
					delete(_rateLocalLimiterIPs, ip)
					_rateLocalLimiterIPsRWMutex.Unlock()
					_rateLocalLimiterIPsRWMutex.RLock()
				}
			}
			_rateLocalLimiterIPsRWMutex.RUnlock()
		}
	}()
}
