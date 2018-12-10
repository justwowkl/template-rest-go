package main // import "main"

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"

	"local/api"
	"local/util"
	"local/custommw"
)

// https://github.com/go-playground/validator
type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {

	util.Init()
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &customValidator{validator: validator.New()}

	// healthcheck API
	// need restricted access
	e.GET("/health", api.Health)

	// public API
	ePub := e.Group("/pub")
	ePub.GET("/test", api.PubTest)
	ePub.GET("/testData", api.PubTestData)
	ePub.POST("/signin", api.PubSignin)

	// user API
	// need auth with JWT
	eUser := e.Group("/user")
	eUser.Use(middleware.JWT([]byte(util.KeyJWT)))
	eUser.Use(custommw.JwtLoader)
	eUser.GET("/me", api.UserMe)

	// admin API
	// need auth with JWT & DB query
	// eAdmin := e.Group("/admin")

	e.Logger.Fatal(e.Start(":3000"))
}
