package db

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/model"
)

type MySQL struct {
	conn *sql.DB
}

func NewMySQL(config *config.Config) (*MySQL, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.Mysqldatabase.User, config.Mysqldatabase.Password, config.Mysqldatabase.Host, config.Mysqldatabase.Port, "employees")
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	return &MySQL{conn: db}, nil
}

func (mysql *MySQL) FindAllEmployees() ([]model.Employee, error) {
	employees := make([]model.Employee, 0)
	stmt, err := mysql.conn.Prepare("SELECT first_name, last_name, salary, birthday, employee_number, entry_date FROM employees")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tmp := model.Employee{}
		err := rows.Scan(&tmp.FirstName, &tmp.LastName, &tmp.Salary, &tmp.Birthday.Time, &tmp.EmployeeNumber, &tmp.EntryDate.Time)
		if err != nil {
			return nil, err
		}
		employees = append(employees, tmp)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (mysql *MySQL) FindAllEmployeesLimit(limit string) ([]model.Employee, error) {
	stmt, err := mysql.conn.Prepare("SELECT * FROM employees LIMIT ?")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	l, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}

	employees := make([]model.Employee, 0, l)
	for rows.Next() {
		tmp := model.Employee{}
		err := rows.Scan(&tmp.EmployeeNumber, &tmp.FirstName, &tmp.LastName, &tmp.Salary, &tmp.Birthday.Time, &tmp.EmployeeNumber, &tmp.EntryDate.Time)
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

func (mysql *MySQL) FindEmployee(id string) (model.Employee, error) {
	result := model.Employee{}
	stmt, err := mysql.conn.Prepare("SELECT first_name, last_name, salary, birthday, employee_number, entry_date FROM employees WHERE employee_number = ?")
	if err != nil {
		return model.Employee{}, err
	}

	err = stmt.QueryRow(id).Scan(&result.FirstName, &result.LastName, &result.Salary, &result.Birthday.Time, &result.EmployeeNumber, &result.EntryDate.Time)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return result, model.NotFoundError
		default:
			return result, err
		}
	}
	return result, nil
}
func (mysql *MySQL) AddEmployee(e model.Employee) error {
	stmt, err := mysql.conn.Prepare("INSERT INTO employees (id, first_name, last_name, salary, birthday, employee_number, entry_date) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id.String(), e.FirstName, e.LastName, e.Salary, e.Birthday.Time, e.EmployeeNumber, e.EntryDate.Time)
	if err != nil {
		return err
	}

	return nil
}
func (mysql *MySQL) RemoveEmployee(id string) error {
	stmt, err := mysql.conn.Prepare("DELETE FROM employees WHERE employee_number = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
func (mysql *MySQL) GetCars(id string) ([]model.Car, error) {
	stmt, err := mysql.conn.Prepare("SELECT employees.first_name, employees.last_name, companycars.number_plate, companycars.type, employees.employee_number FROM employees JOIN companycars ON employees.employee_number=companycars.employee_number WHERE employees.employee_number = ?")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	cars := make([]model.Car, 0)
	for rows.Next() {
		tmp := new(model.Car)
		err := rows.Scan(&tmp.FirstName, &tmp.LastName, &tmp.NumberPlate, &tmp.Type, &tmp.EmployeeNumber)
		if err != nil {
			fmt.Println(err)
		}
		cars = append(cars, *tmp)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return cars, nil
}

func (mysql *MySQL) AddCar(car model.StorageCar) error {
	stmt, err := mysql.conn.Prepare("INSERT INTO companycars (id, manufacturer, type, number_plate, employee_number) VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id.String(), car.Manufacturer, car.Type, car.NumberPlate, car.EmployeeNumber)
	if err != nil {
		return err
	}

	return nil
}
