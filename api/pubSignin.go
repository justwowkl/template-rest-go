package api

import (
	"local/util"
	"net/http"

	"github.com/labstack/echo"
)

// PubSignin hey
func PubSignin(c echo.Context) error {

	type requestScheme struct {
		Provider  string `json:"provider" validate:"required,oneof=fb gg"`
		AuthToken string `json:"authToken" validate:"required,alphanum"`
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
	data := util.JwtData{ID: 88, IP: "0.0.0.0"} //c.RealIP()
	token, err := util.JwtCreate(data)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// -----------------
	responseJSON := &responseScheme{
		Token: token,
	}
	if err := c.Validate(responseJSON); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, responseJSON)
}
