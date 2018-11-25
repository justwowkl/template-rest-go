package api

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// PubSignin hey
func PubSignin(c echo.Context) error {

	type requestScheme struct {
		Name string `json:"name" validate:"required"`
		Msg  string `json:"msg" validate:"required"`
		Num  int    `json:"num" validate:"required"`
	}

	type responseScheme struct {
		Token string `validate:"required"`
	}

	requestJSON := new(requestScheme)
	if err := c.Bind(requestJSON); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(requestJSON); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// need error standard

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = requestJSON.Name
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	responseJSON := &responseScheme{
		Token: t,
	}
	if err := c.Validate(responseJSON); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, responseJSON)
}
