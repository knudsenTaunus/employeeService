package production

import (
	"database/sql"
	"fmt"
	"github.com/knudsenTaunus/employeeService/internal/types"
)

const (
	user     = "root"
	password = "test"
	address  = "127.0.0.1"
	port     = "3306"
	table    = "employees"
)

type mySQLervice struct {
	con *sql.DB
}

func New() (*mySQLervice, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, address, port, table)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	return &mySQLervice{con: db}, nil
}

func (s *mySQLervice) FindAllEmployees() ([]*types.Employee, error) {
	return nil, nil
}
func (s *mySQLervice) FindAllEmployeesLimit(limit string) ([]*types.Employee, error) {
	return nil, nil
}
func (s *mySQLervice) Find(id string) (*types.Employee, error) {
	return  nil, nil
}
func (s *mySQLervice) Add(employee *types.Employee) error {
	return nil
}
func (s *mySQLervice) Remove(id string) error {
	return nil
}
func (s *mySQLervice) GetCars(id string) ([]*types.EmployeeCars, error) {
	return nil, nil
}
