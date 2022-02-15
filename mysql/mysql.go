package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseCredentials struct {
	Username string
	Password string
	Database string
	Host     string
}

func Connect(credentials DatabaseCredentials) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", credentials.Username, credentials.Password, credentials.Host, credentials.Database))

	return db, err
}

func Run(dbPool *sql.DB, query string) (*sql.Rows, error) {
	rows, err := dbPool.Query(query)

	return rows, err
}
