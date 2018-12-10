module local/custommw

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/labstack/echo v3.3.5+incompatible
	local/util v0.0.0
)

replace local/util => ../util
