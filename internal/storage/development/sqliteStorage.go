package development

import (
	"database/sql"
	"fmt"
	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/types"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type sqliteService struct {
	con *sql.DB
}

func New(config *config.Config) (*sqliteService, error) {
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
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `employees` (`id` INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE ,`first_name` VARCHAR(64), `last_name` VARCHAR(64), `salary` INTEGER, `birthday` TEXT, employee_number INTEGER UNIQUE);")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `companycars` (`id` INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE ,`manufacturer` VARCHAR(64), `type` VARCHAR(64), `number_plate` TEXT, employee_number INTEGER, FOREIGN KEY (employee_number) REFERENCES employees(employee_number));")
	err = insertSampleData(db)
	if err != nil {
		return nil, err
	}
	return &sqliteService{
		con: db,
	}, nil
}

func (sqlite *sqliteService) FindAllEmployees() ([]*types.Employee, error) {
	employees := make([]*types.Employee, 0)
	rows, err := sqlite.con.Query("SELECT * FROM employees")
	defer rows.Close()
	if err != nil {
		log.Fatalf("Could not get from sqliteService %v", err)
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

func (sqlite *sqliteService) FindAllEmployeesLimit(limit string) ([]*types.Employee, error) {
	employees := make([]*types.Employee, 0)
	rows, err := sqlite.con.Query("SELECT * FROM employees LIMIT $1", limit)
	defer rows.Close()
	if err != nil {
		log.Fatalf("Could not get from sqliteService %v", err)
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

func (sqlite *sqliteService) Find(id string) (*types.Employee, error) {
	result := &types.Employee{}
	row := sqlite.con.QueryRow("SELECT * FROM employees WHERE id = $1", id)
	switch err := row.Scan(&result.ID, &result.FirstName, &result.LastName, &result.Salary, &result.Birthday, &result.EmployeeNumber); err {
	case sql.ErrNoRows:
		return result, err
	}
	return result, nil
}

func (sqlite *sqliteService) Add(e *types.Employee) error {
	_, err := sqlite.con.Exec("INSERT INTO employees (first_name, last_name, salary, birthday, employee_number) VALUES ($1,$2,$3,$4, $5)", e.FirstName, e.LastName, e.Salary, e.Birthday, e.EmployeeNumber)
	if err != nil {
		return err
	}
	return nil
}


func (sqlite *sqliteService) Remove(id string) error {
	_, err := sqlite.con.Exec("DELETE FROM employees WHERE id = $1",id)
	if err != nil {
		return err
	}
	return nil
}

func (sqlite *sqliteService) GetCars(id string) ([]*types.EmployeeCars, error) {
	rows, err := sqlite.con.Query("SELECT employees.id, employees.first_name, employees.last_name, companycars.number_plate, companycars.type FROM employees JOIN companycars ON employees.employee_number=companycars.employee_number WHERE employees.id = $1", id)
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


func insertSampleData(db *sql.DB) error {
	employees := []*types.Employee{
		{
			FirstName: "Joe",
			LastName:  "Biehl",
			Salary:    50000,
			Birthday: types.JsonDate{
				Time: time.Now(),
			},
			EmployeeNumber: 1,
		},
		{
			FirstName: "Jan",
			LastName:  "Heinrich",
			Salary:    500000,
			Birthday: types.JsonDate{
				Time: time.Now().AddDate(0,0,1),
			},
			EmployeeNumber: 2,
		},{
			FirstName: "Rusalka",
			LastName:  "Ertel",
			Salary:    250000,
			Birthday: types.JsonDate{
				Time: time.Now().AddDate(0,0,2),
			},
			EmployeeNumber: 3,
		},{
			FirstName: "Tauseef",
			LastName:  "Al-Noor",
			Salary:    70000,
			Birthday: types.JsonDate{
				Time: time.Now().AddDate(0,0,3),
			},
			EmployeeNumber: 4,
		},{
			FirstName: "Lotte",
			LastName:  "Kwandt",
			Salary:    5000000,
			Birthday: types.JsonDate{
				Time: time.Now().AddDate(0,0,4),
			},
			EmployeeNumber: 5,
		},
	}

	cars := []*types.Car{
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
		stmt := "INSERT OR IGNORE INTO employees (first_name, last_name, salary, birthday, employee_number) VALUES (?, ?, ?, ?, ?)"
		_, err := db.Exec(stmt, item.FirstName, item.LastName, item.Salary, item.Birthday, item.EmployeeNumber)
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