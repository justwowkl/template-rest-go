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
var _limiterIPs = make(map[string]*typeLimitIP)
var _limiterIPsRWMutex = new(sync.RWMutex)

// RateLimit limit request by request IP
func RateLimit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()
		if strings.HasPrefix(ip, "10.") {
			return next(c)
		}
		limiter := getIP(c.RealIP())
		if limiter.Allow() == false {
			return c.String(http.StatusTooManyRequests, "")
		}
		return next(c)
	}
}

// Run a background goroutine to remove old entries from the visitors map.
func init() {
	go cleanupIPs()
}

// Create a new rate limiter and add it to the visitors map, using the
// IP address as the key.
func addIP(ip string) *rate.Limiter {
	limiter := rate.NewLimiter(2, 2)
	_limiterIPsRWMutex.Lock()
	_limiterIPs[ip] = &typeLimitIP{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	println("new limiter : ", _limiterIPs[ip])
	_limiterIPsRWMutex.Unlock()
	return limiter
}

// Retrieve and return the rate limiter for the current visitor if it
// already exists. Otherwise call the addVisitor function to add a
// new entry to the map.
func getIP(ip string) *rate.Limiter {
	_limiterIPsRWMutex.RLock()
	limiterIP, exists := _limiterIPs[ip]
	_limiterIPsRWMutex.RUnlock()
	if !exists {
		return addIP(ip)
	}
	println("old limiter : ", limiterIP)
	return limiterIP.limiter
}

// Every minute check the map for visitors that haven't been seen for
// more than 5 minutes and delete the entries.
func cleanupIPs() {
	for {
		time.Sleep(time.Minute)
		_limiterIPsRWMutex.RLock()
		for ip, limiterIP := range _limiterIPs {
			if time.Now().Sub(limiterIP.lastSeen) > 5*time.Minute {
				_limiterIPsRWMutex.RUnlock()
				_limiterIPsRWMutex.Lock()
				delete(_limiterIPs, ip)
				_limiterIPsRWMutex.Unlock()
				_limiterIPsRWMutex.RLock()
			}
		}
		_limiterIPsRWMutex.RUnlock()
	}
}
