package custommw

import (
	"local/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

// AuthUser verify jwt and load user data
func AuthUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// get jwt Token
		jwtToken := c.Request().Header.Get("tkn")
		if jwtToken == "" {
			return c.String(http.StatusUnauthorized, "")
		}

		// verify & get Data
		data, err := util.JwtVerify(jwtToken, c.RealIP())
		if err != nil {
			return c.String(http.StatusUnauthorized, "")
		}
		// set datas to Context
		c.Set("id", data.ID)

		return next(c)

	}
}
