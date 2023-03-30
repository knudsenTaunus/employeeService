package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

type SQLLite struct {
	conn *sql.DB
}

func NewSQLite(config *config.Config) (*SQLLite, error) {
	db, err := sql.Open("sqlite3", config.Sqlitedatabase.Path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS `employees`")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS `companycars`")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `employees` (`id` INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE ,`first_name` VARCHAR(64), `last_name` VARCHAR(64), `salary` INTEGER, `birthday` datetime, employee_number INTEGER UNIQUE, entry_date datetime);")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `companycars` (`id` INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE ,`manufacturer` VARCHAR(64), `type` VARCHAR(64), `number_plate` TEXT, employee_number INTEGER, FOREIGN KEY (employee_number) REFERENCES employees(employee_number));")
	err = insertSampleData(db)
	if err != nil {
		return nil, err
	}
	return &SQLLite{
		conn: db,
	}, nil
}

func (sqlite *SQLLite) FindAllEmployees() (model.StorageEmployees, error) {
	employees := make(model.StorageEmployees, 0)
	rows, err := sqlite.conn.Query("SELECT * FROM employees")
	defer rows.Close()
	if err != nil {
		log.Fatalf("Could not get from sqliteService %v", err)
	}
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

func (sqlite *SQLLite) FindAllEmployeesLimit(limit string) (model.StorageEmployees, error) {
	l, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	employees := make(model.StorageEmployees, l)
	rows, err := sqlite.conn.Query("SELECT * FROM employees LIMIT $1", limit)
	defer rows.Close()
	if err != nil {
		log.Fatalf("Could not get from sqliteService %v", err)
	}
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

func (sqlite *SQLLite) FindEmployee(id string) (model.StorageEmployee, error) {
	result := model.StorageEmployee{}
	row := sqlite.conn.QueryRow("SELECT * FROM employees WHERE employee_number = $1", id)
	switch err := row.Scan(&result.ID, &result.FirstName, &result.LastName, &result.Salary, &result.Birthday, &result.EmployeeNumber, &result.EntryDate); err {
	case sql.ErrNoRows:
		return result, err
	}
	return result, nil
}

func (sqlite *SQLLite) AddEmployee(e model.StorageEmployee) error {
	_, err := sqlite.conn.Exec("INSERT INTO employees (first_name, last_name, salary, birthday, employee_number, entry_date) VALUES ($1,$2,$3,$4,$5,$6)", e.FirstName, e.LastName, e.Salary, e.Birthday, e.EmployeeNumber, e.EntryDate)
	if err != nil {
		return err
	}
	return nil
}

func (sqlite *SQLLite) RemoveEmployee(employeeNumber string) error {
	_, err := sqlite.conn.Exec("DELETE FROM employees WHERE employee_number = $1", employeeNumber)
	if err != nil {
		return err
	}
	return nil
}

func (sqlite *SQLLite) GetCars(id string) ([]model.EmployeeCars, error) {
	rows, err := sqlite.conn.Query("SELECT employees.id, employees.first_name, employees.last_name, companycars.number_plate, companycars.type FROM employees JOIN companycars ON employees.employee_number=companycars.employee_number WHERE employees.id = $1", id)
	if err != nil {
		return nil, err
	}
	var cars []model.EmployeeCars
	for rows.Next() {
		tmp := model.EmployeeCars{}
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

func insertSampleData(db *sql.DB) error {
	today, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		return err
	}
	employees := []*model.StorageEmployee{
		{
			FirstName:      "Joe",
			LastName:       "Biehl",
			Salary:         50000,
			Birthday:       time.Date(1986, 1, 18, 0, 0, 0, 0, time.UTC),
			EmployeeNumber: 1,
			EntryDate:      today,
		},
		{
			FirstName:      "Jan",
			LastName:       "Heinrich",
			Salary:         500000,
			Birthday:       time.Date(1984, 2, 24, 0, 0, 0, 0, time.UTC),
			EmployeeNumber: 2,
			EntryDate:      today,
		}, {
			FirstName:      "Rusalka",
			LastName:       "Ertel",
			Salary:         250000,
			Birthday:       time.Date(1988, 3, 10, 0, 0, 0, 0, time.UTC),
			EmployeeNumber: 3,
			EntryDate:      today,
		}, {
			FirstName:      "Tauseef",
			LastName:       "Al-Noor",
			Salary:         70000,
			Birthday:       time.Date(1987, 11, 2, 0, 0, 0, 0, time.UTC),
			EmployeeNumber: 4,
			EntryDate:      today,
		}, {
			FirstName:      "Lotte",
			LastName:       "Kwandt",
			Salary:         5000000,
			Birthday:       time.Date(1988, 7, 27, 0, 0, 0, 0, time.UTC),
			EmployeeNumber: 5,
			EntryDate:      today,
		},
	}

	cars := []*model.Car{
		{
			Manufacturer:   "Mercedes",
			Type:           "SL600",
			NumberPlate:    "B-Jo 666",
			EmployeeNumber: 1,
		},
		{
			Manufacturer:   "BMW",
			Type:           "525",
			NumberPlate:    "DD-CK 007",
			EmployeeNumber: 5,
		},
	}

	for _, item := range employees {
		stmt := "INSERT OR IGNORE INTO employees (first_name, last_name, salary, birthday, employee_number, entry_date) VALUES (?, ?, ?, ?, ?, ?)"
		_, err := db.Exec(stmt, item.FirstName, item.LastName, item.Salary, item.Birthday, item.EmployeeNumber, item.EntryDate)
		if err != nil {
			fmt.Print(err)
		}
	}

	for _, item := range cars {
		stmt := "INSERT OR IGNORE INTO companycars (manufacturer, type, number_plate, employee_number) VALUES (?, ?, ?, ?)"
		_, err := db.Exec(stmt, item.Manufacturer, item.Type, item.NumberPlate, item.EmployeeNumber)
		if err != nil {
			fmt.Print(err)
		}
	}
	return nil
}
