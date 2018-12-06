package util

import (
	"fmt"

	"github.com/jackc/pgx"
)

var _pgConfig *pgx.ConnConfig

// DBinit init db config
func DBinit() {
	_pgConfig = &pgx.ConnConfig{
		Host:     "localhost",
		User:     "postgres",
		Password: "1029",
		Database: "pgx_test",
	}
	fmt.Println("db init done")
}

// DBhealth healthcheck
func DBhealth() bool {
	conn, err := pgx.Connect(*_pgConfig)
	if err != nil || conn.PID() == 0 {
		return false
	}
	if _, present := conn.RuntimeParams["server_version"]; !present {
		return false
	}
	var currentDB string
	err = conn.QueryRow("select current_database()").Scan(&currentDB)
	if err != nil || currentDB != "pgx_test" {
		return false
	}
	var user string
	err = conn.QueryRow("select current_user").Scan(&user)
	if err != nil || user != "postgres" {
		return false
	}
	err = conn.Close()
	if err != nil {
		return false
	}
	fmt.Println("db health okay")
	return true
}
