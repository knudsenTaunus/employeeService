package db

import (
	"database/sql"
	"fmt"
	migrate "github.com/rubenv/sql-migrate"
	"log"
	"strconv"

	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/model"
)

type MySQL struct {
	conn *sql.DB
}

func NewMySQL(config *config.Config) (*MySQL, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.Mysqldatabase.User, config.Mysqldatabase.Password, config.Mysqldatabase.Host, config.Mysqldatabase.Port, "employees")

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		return nil, err
	}

	log.Default().Printf("Applied %d migration files", n)

	return &MySQL{conn: db}, nil
}

func (mysql *MySQL) FindAllEmployees() (model.StorageEmployees, error) {
	employees := make([]model.StorageEmployee, 0)
	stmt, err := mysql.conn.Prepare("SELECT * FROM employees")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tmp := model.StorageEmployee{}
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
func (mysql *MySQL) FindAllEmployeesLimit(limit string) (model.StorageEmployees, error) {
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

	employees := make([]model.StorageEmployee, 0, l)
	for rows.Next() {
		tmp := model.StorageEmployee{}
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

func (mysql *MySQL) FindEmployee(id string) (model.StorageEmployee, error) {
	result := model.StorageEmployee{}
	stmt, err := mysql.conn.Prepare("SELECT * FROM employees WHERE employee_number = ?")
	if err != nil {
		return model.StorageEmployee{}, err
	}

	err = stmt.QueryRow(id).Scan(&result.ID, &result.FirstName, &result.LastName, &result.Salary, &result.Birthday, &result.EmployeeNumber, &result.EntryDate)
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
func (mysql *MySQL) AddEmployee(e model.StorageEmployee) error {
	stmt, err := mysql.conn.Prepare("INSERT INTO employees (first_name, last_name, salary, birthday, employee_number, entry_date) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(e.FirstName, e.LastName, e.Salary, e.Birthday, e.EmployeeNumber, e.EntryDate)
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
func (mysql *MySQL) GetCars(id string) ([]model.EmployeeCars, error) {
	stmt, err := mysql.conn.Prepare("SELECT employees.id, employees.first_name, employees.last_name, companycars.number_plate, companycars.type FROM employees JOIN companycars ON employees.employee_number=companycars.employee_number WHERE employees.employee_number = ?")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	
	cars := make([]model.EmployeeCars, 0)
	for rows.Next() {
		tmp := new(model.EmployeeCars)
		err := rows.Scan(&tmp.ID, &tmp.FirstName, &tmp.LastName, &tmp.NumberPlate, &tmp.Type)
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
