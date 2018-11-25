package util

import (
	"github.com/labstack/echo"
)

//ParseRequest try parse request data
func ParseRequest(c echo.Context, data interface{}) (interface{}, error) {
	if err := c.Bind(data); err != nil {
		return nil, err
	}
	if err := c.Validate(data); err != nil {
		return nil, err
	}
	return data, nil
}
