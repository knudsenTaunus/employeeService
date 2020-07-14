package storage

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/storage/development"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/storage/production"
	"log"

	"github.com/janPhil/mySQLHTTPRestGolang/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type service struct {
	con *sql.DB
}

func New() (*service, error) {
	env := flag.String("environment", "development", "environment to run the app in")
	flag.Parse()
	if *env == "development" {
		db, err := development.New()
		if err != nil {
			return nil, err
		}
		return &service{
			con: db,
		}, nil
	}
	db, err := production.New()
	if err != nil {
		return nil, err
	}
	return &service{
		con: db,
	}, nil
}

func (ed *service) FindAllEmployees() ([]*types.Employee, error) {
	employees := make([]*types.Employee, 0)
	rows, err := ed.con.Query("SELECT * FROM employees")
	defer rows.Close()
	if err != nil {
		log.Fatalf("Could not get from service %v", err)
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

func (ed *service) FindAllEmployeesLimit(limit string) ([]*types.Employee, error) {
	employees := make([]*types.Employee, 0)
	rows, err := ed.con.Query("SELECT * FROM employees LIMIT $1", limit)
	defer rows.Close()
	if err != nil {
		log.Fatalf("Could not get from service %v", err)
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

func (ed *service) Find(id string) (*types.Employee, error) {
	result := &types.Employee{}
	row := ed.con.QueryRow("SELECT * FROM employees WHERE id = $1", id)
	switch err := row.Scan(&result.ID, &result.FirstName, &result.LastName, &result.Salary, &result.Birthday, &result.EmployeeNumber); err {
	case sql.ErrNoRows:
		return result, err
	}
	return result, nil
}

func (ed *service) Add(e *types.Employee) error {
	_, err := ed.con.Exec("INSERT INTO employees (first_name, last_name, salary, birthday, employee_number) VALUES ($1,$2,$3,$4, $5)", e.FirstName, e.LastName, e.Salary, e.Birthday, e.EmployeeNumber)
	if err != nil {
		return err
	}
	return nil
}


func (ed *service) Remove(id string) error {
	_, err := ed.con.Exec("DELETE FROM employees WHERE id = $1",id)
	if err != nil {
		return err
	}
	return nil
}

func (ed *service) GetCars(id string) ([]*types.EmployeeCars, error) {
	rows, err := ed.con.Query("SELECT employees.id, employees.first_name, employees.last_name, companycars.number_plate, companycars.type FROM employees JOIN companycars ON employees.employee_number=companycars.employee_number WHERE employees.id = $1", id)
	if err != nil {
		return nil, err
	}
	var cars []*types.EmployeeCars
	for rows.Next() {
		tmp := new(types.EmployeeCars)
		err := rows.Scan(&tmp.ID, &tmp.FirstName, &tmp.LastName, &tmp.NumberPlate, &tmp.Type)
		if err != nil {
			fmt.Println(err)
		}
		cars = append(cars, tmp)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return cars, nil
}