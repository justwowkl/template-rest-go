package api

import (
	"net/http"
	"sync"

	"local/util"

	"github.com/labstack/echo"
)

// set for local test
// sudo docker run --name postgres -e POSTGRES_PASSWORD=1029 -p 5432:5432 -d postgres:10-alpine
// sudo docker run --name redis -p 6379:6379 -d redis:alpine

var _isFullCheckDone = false
var _rwMutexHealth = new(sync.RWMutex)

// Health health check
func Health(c echo.Context) error {
	_rwMutexHealth.RLock()
	isFullCheckDone := _isFullCheckDone
	_rwMutexHealth.RUnlock()
	if isFullCheckDone {
		// simple test
		return c.String(http.StatusOK, "")
	}

	// do full test
	if util.Health() {
		// update full test flag
		_rwMutexHealth.Lock()
		_isFullCheckDone = true
		_rwMutexHealth.Unlock()

		// return succeed
		return c.String(http.StatusOK, "")
	}

	return c.String(http.StatusInternalServerError, "")
}
