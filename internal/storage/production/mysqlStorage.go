package production

import (
	"database/sql"
	"fmt"
	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/types"
	"log"
)

const (
	user     = "root"
	password = "employeedatabase"
	host  = "127.0.0.1"
	port     = "3306"
	table    = "employees"
)

type mySQLService struct {
	con *sql.DB
}

func New(config *config.Config) (*mySQLService, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.Mysqldatabase.User, config.Mysqldatabase.Password, config.Mysqldatabase.Host, config.Mysqldatabase.Port, "employees")
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	return &mySQLService{con: db}, nil
}

func (mysql *mySQLService) FindAllEmployees() ([]*types.Employee, error) {
	employees := make([]*types.Employee, 0)
	rows, err := mysql.con.Query("SELECT * FROM employees")
	defer rows.Close()
	if err != nil {
		log.Fatalf("Could not get from sqliteService %v", err)
	}
	for rows.Next() {
		tmp := new(types.Employee)
		err := rows.Scan(&tmp.ID, &tmp.FirstName, &tmp.LastName, &tmp.Salary, &tmp.Birthday.Time, &tmp.EmployeeNumber, &tmp.EntryDate.Time)
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
func (mysql *mySQLService) FindAllEmployeesLimit(limit string) ([]*types.Employee, error) {
	employees := make([]*types.Employee, 0)
	rows, err := mysql.con.Query("SELECT * FROM employees LIMIT ?", limit)
	defer rows.Close()
	if err != nil {
		log.Fatalf("Could not get from sqliteService %v", err)
	}
	for rows.Next() {
		tmp := new(types.Employee)
		err := rows.Scan(&tmp.ID, &tmp.FirstName, &tmp.LastName, &tmp.Salary, &tmp.Birthday.Time, &tmp.EmployeeNumber, &tmp.EntryDate.Time)
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
func (mysql *mySQLService) Find(id string) (*types.Employee, error) {
	result := &types.Employee{}
	row := mysql.con.QueryRow("SELECT * FROM employees WHERE id = ?", id)
	switch err := row.Scan(&result.ID, &result.FirstName, &result.LastName, &result.Salary, &result.Birthday, &result.EmployeeNumber, &result.EntryDate); err {
	case sql.ErrNoRows:
		return result, err
	}
	return result, nil
}
func (mysql *mySQLService) Add(e *types.Employee) error {
	_, err := mysql.con.Exec("INSERT INTO employees (first_name, last_name, salary, birthday, employee_number, entry_date) VALUES (?,?,?,?,?,?)", e.FirstName, e.LastName, e.Salary, e.Birthday.Time, e.EmployeeNumber, e.EntryDate.Time)
	if err != nil {
		return err
	}
	return nil
}
func (mysql *mySQLService) Remove(id string) error {
	_, err := mysql.con.Exec("DELETE FROM employees WHERE id = ?",id)
	if err != nil {
		return err
	}
	return nil
}
func (mysql *mySQLService) GetCars(id string) ([]*types.EmployeeCars, error) {
	rows, err := mysql.con.Query("SELECT employees.id, employees.first_name, employees.last_name, companycars.number_plate, companycars.type FROM employees JOIN companycars ON employees.employee_number=companycars.employee_number WHERE employees.id = ?", id)
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
