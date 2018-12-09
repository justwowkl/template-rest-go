package util

import (
	"fmt"

	"github.com/jackc/pgx"
)

var _pgConfig *pgx.ConnConfig

// dbInit init db config
func pgInit() {
	_pgConfig = &pgx.ConnConfig{
		Host:     "localhost",
		User:     "postgres",
		Password: "1029",
		Database: "pgx_test",
	}
	fmt.Println("db init done")
}

// dbHealth healthcheck
func pgHealth() bool {
	conn, err := pgx.Connect(*_pgConfig)
	if err != nil {
		return false
	}
	defer conn.Close()
	if conn.PID() == 0 {
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
	fmt.Println("db health okay")
	return true
}

// PGquerySingle healthcheck
func PGquerySingle(query string) (interface{}, error) {
	conn, errConn := pgx.Connect(*_pgConfig)
	if errConn != nil {
		return nil, errConn
	}
	defer conn.Close()
	rows := (*pgx.Rows)(conn.QueryRow(query))
	defer rows.Close()
	errQuery := rows.Err()
	if errQuery != nil {
		return nil, errQuery
	}
	vals, errVal := rows.Values()
	if errVal != nil {
		return nil, errVal
	}
	return vals[0], nil
}

// PGqueryMultiple healthcheck
func PGqueryMultiple(query string) ([]interface{}, error) {
	conn, errConn := pgx.Connect(*_pgConfig)
	if errConn != nil {
		return nil, errConn
	}
	defer conn.Close()
	rows, errQuery := conn.Query(query)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()
	return rows.Values()
}

// PGexec healthcheck
func PGexec(query string) (int, error) {
	conn, errConn := pgx.Connect(*_pgConfig)
	if errConn != nil {
		return -1, errConn
	}
	defer conn.Close()
	result, errExec := conn.Exec(query)
	if errExec != nil {
		return -1, errExec
	}
	return int(result.RowsAffected()), nil
}
