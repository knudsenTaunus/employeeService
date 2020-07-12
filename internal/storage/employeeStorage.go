package storage

import (
	"database/sql"
	"flag"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/storage/development"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/storage/production"
	"log"

	"github.com/janPhil/mySQLHTTPRestGolang/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type EmployeeStorage struct {
	con *sql.DB
}

func New() (*EmployeeStorage, error) {
	env := flag.String("environment", "development", "environment to run the app in")
	flag.Parse()
	if *env == "development" {
		db, err := development.NewSQLiteDatabase()
		if err != nil {
			return nil, err
		}
		return &EmployeeStorage{
			con: db,
		}, nil
	}
	db, err := production.NewSQLDatabase()
	if err != nil {
		return nil, err
	}
	return &EmployeeStorage{
		con: db,
	}, nil
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
		err := rows.Scan(&tmp.ID, &tmp.FirstName, &tmp.LastName, &tmp.Salary, &tmp.Birthday, &tmp.EmployeeNumber)
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

func (ed *EmployeeStorage) FindAllEmployeesLimit(limit string) ([]*types.Employee, error) {
	employees := make([]*types.Employee, 0)
	rows, err := ed.con.Query("SELECT * FROM employees LIMIT $1", limit)
	defer rows.Close()
	if err != nil {
		log.Fatalf("Could not get from storage %v", err)
	}
	for rows.Next() {
		tmp := new(types.Employee)
		err := rows.Scan(&tmp.ID, &tmp.FirstName, &tmp.LastName, &tmp.Salary, &tmp.Birthday)
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
	switch err := row.Scan(&result.ID, &result.FirstName, &result.LastName, &result.Salary, &result.Birthday); err {
	case sql.ErrNoRows:
		return result, err
	}
	return result, nil
}

func (ed *EmployeeStorage) Add(e *types.Employee) error {
	_, err := ed.con.Exec("INSERT INTO employees (first_name, last_name, salary, birthday, employee_number) VALUES ($1,$2,$3,$4, $5)", e.FirstName, e.LastName, e.Salary, e.Birthday, e.EmployeeNumber)
	if err != nil {
		return err
	}
	return nil
}


func (ed *EmployeeStorage) Remove(id string) error {
	_, err := ed.con.Exec("DELETE FROM employees WHERE id = $1",id)
	if err != nil {
		return err
	}
	return nil
}


