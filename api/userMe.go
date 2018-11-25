package api

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// UserMe health check
func UserMe(c echo.Context) error {

	type responseScheme struct {
		Name    string `validate:"required"`
		IsAdmin bool   `validate:"required"`
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	responseJSON := &responseScheme{
		Name:    claims["name"].(string),
		IsAdmin: claims["admin"].(bool),
	}
	if err := c.Validate(responseJSON); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, responseJSON)
}
