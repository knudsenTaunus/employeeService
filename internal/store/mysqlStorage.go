package store

import (
	"database/sql"
	"fmt"
	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/types"
	"log"
	"strconv"
)

const (
	user     = "root"
	password = "employeedatabase"
	host     = "127.0.0.1"
	port     = "3306"
	table    = "employees"
)

type mySQLService struct {
	con *sql.DB
}

func NewMySQL(config *config.Config) (*mySQLService, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.Mysqldatabase.User, config.Mysqldatabase.Password, config.Mysqldatabase.Host, config.Mysqldatabase.Port, "employees")
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	return &mySQLService{con: db}, nil
}

func (mysql *mySQLService) FindAllEmployees() (types.StorageEmployees, error) {
	employees := make([]*types.StorageEmployee, 0)
	rows, err := mysql.con.Query("SELECT * FROM employees")
	defer rows.Close()
	if err != nil {
		log.Fatalf("Could not get from sqliteService %v", err)
	}
	for rows.Next() {
		tmp := new(types.StorageEmployee)
		err := rows.Scan(&tmp.ID, &tmp.FirstName, &tmp.LastName, &tmp.Salary, &tmp.Birthday, &tmp.EmployeeNumber, &tmp.EntryDate)
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
func (mysql *mySQLService) FindAllEmployeesLimit(limit string) (types.StorageEmployees, error) {
	l, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	employees := make([]*types.StorageEmployee, l)
	rows, err := mysql.con.Query("SELECT * FROM employees LIMIT ?", limit)
	if err != nil {
		log.Fatalf("Could not get from sqliteService %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		tmp := new(types.StorageEmployee)
		err := rows.Scan(&tmp.ID, &tmp.FirstName, &tmp.LastName, &tmp.Salary, &tmp.Birthday, &tmp.EmployeeNumber, &tmp.EntryDate)
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
func (mysql *mySQLService) Find(id string) (*types.StorageEmployee, error) {
	result := &types.StorageEmployee{}
	row := mysql.con.QueryRow("SELECT * FROM employees WHERE employee_number = ?", id)
	switch err := row.Scan(&result.ID, &result.FirstName, &result.LastName, &result.Salary, &result.Birthday, &result.EmployeeNumber, &result.EntryDate); err {
	case sql.ErrNoRows:
		return result, err
	}
	return result, nil
}
func (mysql *mySQLService) Add(e *types.StorageEmployee) error {
	_, err := mysql.con.Exec("INSERT INTO employees (first_name, last_name, salary, birthday, employee_number, entry_date) VALUES (?,?,?,?,?,?)", e.FirstName, e.LastName, e.Salary, e.Birthday, e.EmployeeNumber, e.EntryDate)
	if err != nil {
		return err
	}
	return nil
}
func (mysql *mySQLService) Remove(id string) error {
	_, err := mysql.con.Exec("DELETE FROM employees WHERE employee_number = ?", id)
	if err != nil {
		return err
	}
	return nil
}
func (mysql *mySQLService) GetCars(id string) ([]*types.EmployeeCars, error) {
	rows, err := mysql.con.Query("SELECT employees.id, employees.first_name, employees.last_name, companycars.number_plate, companycars.type FROM employees JOIN companycars ON employees.employee_number=companycars.employee_number WHERE employees.employee_number = ?", id)
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
