module main

require (
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c // indirect
	local/api v0.0.0
	local/custommw v0.0.0
	local/util v0.0.0
	local/worker v0.0.0
)

replace local/api => ./api

replace local/util => ./util

replace local/custommw => ./custommw

replace local/worker => ./worker
