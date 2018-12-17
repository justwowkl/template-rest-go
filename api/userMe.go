package api

import (
	"net/http"

	"github.com/labstack/echo"
)

// UserMe health check
func UserMe(c echo.Context) error {

	println("start me")
	type responseScheme struct {
		ID int `json:"id" validate:"required"`
	}

	responseJSON := &responseScheme{
		ID: c.Get("id").(int),
	}
	println("hi!")
	println("id : ", responseJSON.ID)
	if err := c.Validate(responseJSON); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, responseJSON)
}
