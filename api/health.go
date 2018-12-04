package api

import (
	"net/http"
	"sync"

	"local/util"

	"github.com/labstack/echo"
)

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
	if util.CacheHealth() &&
		util.CacheHealth() &&
		util.CacheHealth() {

		// update full test flag
		_rwMutexHealth.Lock()
		_isFullCheckDone = true
		_rwMutexHealth.Unlock()

		// return succeed
		return c.String(http.StatusOK, "")
	}

	return c.String(http.StatusInternalServerError, "")
}
