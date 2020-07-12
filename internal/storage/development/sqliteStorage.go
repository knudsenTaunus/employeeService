package development

import (
	"database/sql"
	"fmt"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/types"
)

func NewSQLiteDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "../development.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS `employees`")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `employees` (`id` INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE ,`first_name` VARCHAR(64), `last_name` VARCHAR(64), `salary` INTEGER, `birthday` TEXT, employee_number INTEGER UNIQUE);")
	if err != nil {
		return nil, err
	}
	err = insertSampleData(db)

	if err != nil {
		return nil, err
	}
	return db, nil
}

func insertSampleData(db *sql.DB) error {
	data := []*types.Employee{
		{
			FirstName: "Joe",
			LastName:  "Biehl",
			Salary:    50000,
			Birthday: "14.07.1988",
			EmployeeNumber: 1,
		},
		{
			FirstName: "Jan",
			LastName:  "Heinrich",
			Salary:    500000,
			Birthday: "24.02.1984",
			EmployeeNumber: 2,
		},{
			FirstName: "Rusalka",
			LastName:  "Ertel",
			Salary:    250000,
			Birthday: "06.03.1988",
			EmployeeNumber: 3,
		},{
			FirstName: "Tauseef",
			LastName:  "Al-Noor",
			Salary:    70000,
			Birthday: "30.08.1988",
			EmployeeNumber: 4,
		},{
			FirstName: "Lotte",
			LastName:  "Kwandt",
			Salary:    5000000,
			Birthday: "26.07.1988",
			EmployeeNumber: 5,
		},
	}

	for _, item := range data {
		stmt := "INSERT OR IGNORE INTO employees (first_name, last_name, salary, birthday, employee_number) VALUES (?, ?, ?, ?, ?)"
		_, err := db.Exec(stmt, item.FirstName, item.LastName, item.Salary, item.Birthday, item.EmployeeNumber)
		if err != nil {
			fmt.Print(err)
		}
	}
	return nil
}