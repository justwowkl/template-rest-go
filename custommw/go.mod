module local/custommw

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/labstack/echo v3.3.5+incompatible
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c
	local/util v0.0.0
)

replace local/util => ../util
