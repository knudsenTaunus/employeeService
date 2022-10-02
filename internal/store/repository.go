package store

import "github.com/knudsenTaunus/employeeService/internal/types"

type Database interface {
	GetCars(id string) ([]*types.EmployeeCars, error)
	FindAllEmployees() (types.StorageEmployees, error)
	FindAllEmployeesLimit(limit string) (types.StorageEmployees, error)
	Find(id string) (types.StorageEmployee, error)
	Add(employee types.StorageEmployee) error
	Remove(id string) error
}

type Repository struct {
	DB Database
}

func NewRepository(db Database) *Repository {
	return &Repository{DB: db}
}

func (r Repository) GetCars(id string) ([]*types.EmployeeCars, error) {
	return r.DB.GetCars(id)
}

func (r Repository) FindAllEmployees() (types.StorageEmployees, error) {
	return r.DB.FindAllEmployees()
}

func (r Repository) FindAllEmployeesLimit(limit string) (types.StorageEmployees, error) {
	return r.DB.FindAllEmployeesLimit(limit)
}

func (r Repository) Find(id string) (types.StorageEmployee, error) {
	return r.DB.Find(id)
}

func (r Repository) Add(employee types.StorageEmployee) error {
	return r.DB.Add(employee)
}

func (r Repository) Remove(id string) error {
	return r.DB.Remove(id)
}
