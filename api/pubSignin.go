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
		Provider  string `json:"provider" validate:"required, oneof=fb gg"`
		AuthToken string `json:"authToken" validate:"required, Alphanumeric"`
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
	// -----------------
	// need error standard

	// async - ask to DB & get info
	// async - ask to Porvider

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	// claims["name"] = requestJSON.Name
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// -----------------
	responseJSON := &responseScheme{
		Token: t,
	}
	if err := c.Validate(responseJSON); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, responseJSON)
}
