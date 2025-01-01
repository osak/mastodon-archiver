package db

import (
    "fmt"
    "github.com/jmoiron/sqlx"
    _ "github.com/go-sql-driver/mysql"
)

func Connect(host string, port int, user string, password string, dbname string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbname)
	return sqlx.Connect("mysql", connStr)
}
