package production

import (
	"database/sql"
	"fmt"
)

const (
	user     = "root"
	password = "test"
	address  = "127.0.0.1"
	port     = "3306"
	table    = "employees"
)

func NewSQLDatabase() (*sql.DB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, address, port, table)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
