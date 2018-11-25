package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// Test test api
func Test(c echo.Context) error {

	// payload scheme
	type requestScheme struct {
		Name string `json:"name" validate:"required"`
		Msg  string `json:"msg" validate:"required"`
		Num  int    `json:"num" validate:"required"`
	}
	type responseScheme struct {
		Name string `validate:"required"`
		Msg  string `validate:"required"`
		IP   string `validate:"required"`
		Num  int    `validate:"required"`
	}

	// verify request payload
	requestJSON := new(requestScheme)
	if err := c.Bind(requestJSON); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(requestJSON); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// need error standard

	// verify response payload
	responseJSON := &responseScheme{
		Name: requestJSON.Name,
		Msg:  requestJSON.Msg,
		Num:  requestJSON.Num,
		IP:   c.RealIP(),
	}
	if err := c.Validate(responseJSON); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, responseJSON)
}
