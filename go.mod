module main

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/labstack/echo v3.3.5+incompatible
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.2.2 // indirect
	golang.org/x/sys v0.0.0-20181122145206-62eef0e2fa9b // indirect
	gopkg.in/go-playground/validator.v9 v9.23.0
	local/api v0.0.0
	local/custommw v0.0.0
	local/util v0.0.0
)

replace local/api => ./api

replace local/util => ./util

replace local/custommw => ./middleware
