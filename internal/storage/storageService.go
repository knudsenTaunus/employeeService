package storage

import (
	"github.com/knudsenTaunus/employeeService/internal/storage/development"
	"github.com/knudsenTaunus/employeeService/internal/storage/production"
	"github.com/knudsenTaunus/employeeService/internal/types"
	"log"
)

type storage interface {
	FindAllEmployees() ([]*types.Employee, error)
	FindAllEmployeesLimit(limit string) ([]*types.Employee, error)
	Find(id string) (*types.Employee, error)
	Add(employee *types.Employee) error
	Remove(id string) error
	GetCars(id string) ([]*types.EmployeeCars, error)
}

func New(env string) storage {
	if env == "development" {
		var err error
		db, err := development.New()
		if err != nil {
			log.Fatal("failed to create development database")
		}
		return db
	}
	if env == "production" {
		var err error
		db, err := production.New()
		if err != nil {
			log.Fatal("failed to create production database")
		}
		return db
	}
	return nil
}
