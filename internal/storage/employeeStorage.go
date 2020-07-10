package storage

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/janPhil/mySQLHTTPRestGolang/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const (
	user     = "root"
	password = "test"
	address  = "127.0.0.1"
	port     = "3306"
	table    = "employees"
)

// NewDB builds the connection to a storage and returns a handle
// if the flag environment is set to development a local sqlite is created, sample
// data is inserted and the connection to it is returned
// If the flag is set to anything else the connection to an existing storage is
// established according to the const values.



type EmployeeStorage struct {
	con *sql.DB
}

func New() (*EmployeeStorage, error) {
	env := flag.String("environment", "development", "environment to run the app in")
	flag.Parse()
	if *env == "development" {
		db, err := newSQLiteDatabase()
		if err != nil {
			return nil, err
		}
		return db, nil
	}
	db, err := newSQLDatabase()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (ed *EmployeeStorage) FindAllEmployees() ([]*types.Employee, error) {
	employees := make([]*types.Employee, 0)
	rows, err := ed.con.Query("SELECT * FROM employees")
	defer rows.Close()
	if err != nil {
		log.Fatalf("Could not get from storage %v", err)
	}
	for rows.Next() {
		tmp := new(types.Employee)
		err := rows.Scan(&tmp.ID, &tmp.FirstName, &tmp.LastName, &tmp.Salary)
		employees = append(employees, tmp)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (ed *EmployeeStorage) Find(id string) (*types.Employee, error) {
	result := &types.Employee{}
	row := ed.con.QueryRow("SELECT * FROM employees WHERE id = $1", id)
	switch err := row.Scan(&result.ID, &result.FirstName, &result.LastName, &result.Salary); err {
	case sql.ErrNoRows:
		return result, err
	}
	return result, nil

}


func newSQLiteDatabase() (*EmployeeStorage, error) {
	db, err := sql.Open("sqlite3", "../development.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `employees` (`id` INTEGER PRIMARY KEY AUTOINCREMENT,`first_name` VARCHAR(64), `last_name` VARCHAR(64), `salary` INTEGER);")
	if err != nil {
		return nil, err
	}
	err = insertSampleData(db)
	if err != nil {
		return nil, err
	}
	return &EmployeeStorage{
		con: db,
	}, nil
}



func newSQLDatabase() (*EmployeeStorage, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, address, port, table)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	return &EmployeeStorage{
		con: db,
	}, nil
}

func insertSampleData(db *sql.DB) error {
	joe := &types.Employee{
		ID: 1,
		FirstName: "Joe",
		LastName:  "Sample",
		Salary:    50000,
	}

	stmt := "INSERT OR IGNORE INTO employees (id, first_name, last_name, salary) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(stmt, joe.ID, joe.FirstName, joe.LastName, joe.Salary)
	if err != nil {
		return err
	}
	return nil
}
