package api

import (
	"net/http"
	"sync"

	"local/util"

	"github.com/labstack/echo/v4"
)

// set for local test
// sudo docker run --name postgres -e POSTGRES_PASSWORD=1029 -p 5432:5432 -d postgres:10-alpine
// sudo docker run --name redis -p 6379:6379 -d redis:alpine

var _healthIsFullChecked = false
var _healthRWMutex = new(sync.RWMutex)
var _healthSuccessHandler func()
var _healthFailHandler func()

// HealthSetSuccessHandler HealthSetSuccessHandler
func HealthSetSuccessHandler(handler func()) {
	_healthSuccessHandler = handler
}

// HealthSetFailHandler HealthSetFailHandler
func HealthSetFailHandler(handler func()) {
	_healthFailHandler = handler
}

// Health health check
func Health(c echo.Context) error {
	_healthRWMutex.RLock()
	isFullCheckDone := _healthIsFullChecked
	_healthRWMutex.RUnlock()
	if isFullCheckDone {
		// simple test
		return c.String(http.StatusOK, "")
	}

	// do full test
	if util.Health() {
		// update full test flag
		_healthRWMutex.Lock()
		_healthIsFullChecked = true
		_healthRWMutex.Unlock()

		// return succeed
		if _healthSuccessHandler != nil {
			_healthSuccessHandler()
		}
		return c.String(http.StatusOK, "")
	}

	if _healthFailHandler != nil {
		_healthFailHandler()
	}
	return c.String(http.StatusInternalServerError, "")
}
