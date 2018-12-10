package custommw

import (
	"local/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// JwtLoader JwtLoader
func JwtLoader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		dataEncrypted := claims["data"].(string)

		// get Data
		data, err := util.JwtVerifyData(dataEncrypted)
		if err != nil {
			c.Error(err)
		}
		// set datas to Context
		c.Set("id", data.ID)
		// clear token in context
		c.Set("user", nil)

		return nil
	}
}
