package main // import "main"

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"

	"local/api"
	"local/custommw"
	"local/util"
	"local/worker"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {

	util.Init()
	// custommw.InitRateLocal()
	custommw.InitRateDist()
	e := echo.New()

	e.Use(middleware.Logger())
	// e.Use(custommw.Logger)
	// e.Use(custommw.RateLimitLocal)
	e.Use(middleware.Recover())
	e.Validator = &customValidator{validator: validator.New()}

	// healthcheck API
	// need restricted access (only internal IP)
	api.HealthSetSuccessHandler(worker.Start)
	api.HealthSetFailHandler(worker.Stop)
	e.GET("/health", api.Health)

	// public API
	ePub := e.Group("/pub")
	ePub.GET("/test", api.PubTest)
	ePub.GET("/testData", api.PubTestData)
	ePub.POST("/signin", api.PubSignin)

	// user API
	// need auth with JWT
	eUser := e.Group("/user")
	// eUser.Use(middleware.JWT([]byte(util.KeyJWT)))
	eUser.Use(custommw.AuthUser)
	eUser.GET("/me", api.UserMe)

	// admin API
	// need auth with JWT & DB query
	// eAdmin := e.Group("/admin")

	e.Logger.Fatal(e.Start(":3000"))
}
