package database

import (
	"database/sql"
	"fmt"
)

// InitDatabase initializes the database - Only to call if one needs a new database
func InitDatabase() {
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1:3306)/employees")
	if err != nil {
		fmt.Printf("Could not connect to database %v", err)
	}

	_, err = db.Query("CREATE Table IF NOT EXISTS employees(id int NOT NULL AUTO_INCREMENT, first_name varchar(50), last_name varchar(50), salary int, PRIMARY KEY(id));")
	if err != nil {
		fmt.Printf("Couldnt create table %v", err)
	}
	defer db.Close()
}

// InsertSampleData puts some employees into the database - Only to populate the new database with some testdata
func InsertSampleData() {

	type employee struct {
		firstName string
		lastName  string
		salary    int
	}

	joe := employee{firstName: "Joe", lastName: "Sample", salary: 50000}
	db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1:3306)/employees")
	if err != nil {
		fmt.Printf("Could not connect to database %v", err)
	}

	stmt := "INSERT IGNORE INTO employees (id, first_name, last_name, salary) VALUES (?, ?, ?, ?)"
	_, err = db.Query(stmt, 1, joe.firstName, joe.lastName, joe.salary)
	if err != nil {
		fmt.Printf("Could not insert sample data %v", err)
	}
}
