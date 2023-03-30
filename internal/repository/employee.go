package repository

import "github.com/knudsenTaunus/employeeService/internal/model"

type EmployeeDatabase interface {
	FindAllEmployees() (model.StorageEmployees, error)
	FindAllEmployeesLimit(limit string) (model.StorageEmployees, error)
	FindEmployee(id string) (model.StorageEmployee, error)
	AddEmployee(employee model.StorageEmployee) error
	RemoveEmployee(id string) error
}

type EmployeeRepository struct {
	DB EmployeeDatabase
}

func NewEmployee(db EmployeeDatabase) *EmployeeRepository {
	return &EmployeeRepository{DB: db}
}

func (r EmployeeRepository) FindAllEmployees() (model.StorageEmployees, error) {
	return r.DB.FindAllEmployees()
}

func (r EmployeeRepository) FindAllEmployeesLimit(limit string) (model.StorageEmployees, error) {
	return r.DB.FindAllEmployeesLimit(limit)
}

func (r EmployeeRepository) FindEmployee(id string) (model.StorageEmployee, error) {
	return r.DB.FindEmployee(id)
}

func (r EmployeeRepository) AddEmployee(employee model.StorageEmployee) error {
	return r.DB.AddEmployee(employee)
}

func (r EmployeeRepository) RemoveEmployee(id string) error {
	return r.DB.RemoveEmployee(id)
}
