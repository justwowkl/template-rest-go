package api

import (
	"net/http"

	"github.com/labstack/echo"
)

var isFullCheckDone = false

// Health health check
func Health(c echo.Context) error {
	if isFullCheckDone {
		return c.String(http.StatusOK, "")
	}
	// do full function test
	return c.String(http.StatusOK, "")
}
