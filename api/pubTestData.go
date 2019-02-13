package api

import (
	"local/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

// PubTestData test api
func PubTestData(c echo.Context) error {

	// verify request payload
	util.CacheSet("test", "okay")
	util.CacheGet("test")

	return c.String(http.StatusOK, "")
}
