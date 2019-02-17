module local/custommw

require (
	github.com/bluele/gcache v0.0.0-20190203144525-2016d595ccb0
	github.com/go-redis/redis v6.15.1+incompatible
	github.com/labstack/echo/v4 v4.0.0
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c
	gopkg.in/go-playground/validator.v9 v9.27.0
	local/util v0.0.0
)

replace local/util => ../util
