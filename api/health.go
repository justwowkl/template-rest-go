package api

import (
	"net/http"
	"sync"

	"github.com/labstack/echo"
)

var _isFullCheckDone = false
var _rwMutex = new(sync.RWMutex)

// Health health check
func Health(c echo.Context) error {
	_rwMutex.RLock()
	isFullCheckDone := _isFullCheckDone
	_rwMutex.RUnlock()
	if isFullCheckDone {
		// simple test
		return c.String(http.StatusOK, "")
	}

	// do full test
	// ...

	// if succeed full test
	_rwMutex.Lock()
	_isFullCheckDone = true
	_rwMutex.Unlock()
	return c.String(http.StatusOK, "")
}
