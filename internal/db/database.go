package db

import (
	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/model"
)

type Database interface {
	FindAllEmployees() ([]model.Employee, error)
	FindAllEmployeesLimit(limit string) ([]model.Employee, error)
	FindEmployee(id string) (model.Employee, error)
	AddEmployee(employee model.Employee) error
	RemoveEmployee(id string) error
	GetCars(id string) ([]model.Car, error)
	AddCar(car model.StorageCar) error
	UpdateEmployee(employee model.Employee) error
}

func New(cfg *config.Config) (Database, error) {
	switch cfg.Environment {
	case "development":
		database, err := NewSQLite(cfg)
		if err != nil {
			return nil, err
		}

		return database, nil
	case "production":
		database, err := NewMySQL(cfg)
		if err != nil {
			return nil, err
		}

		return database, nil
	default:
		database, err := NewSQLite(cfg)
		if err != nil {
			return nil, err
		}

		return database, nil

	}
}
