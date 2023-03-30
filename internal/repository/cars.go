package repository

import "github.com/knudsenTaunus/employeeService/internal/model"

type CarDatabase interface {
	GetCars(id string) ([]model.EmployeeCars, error)
}

type CarRepository struct {
	DB CarDatabase
}

func NewCar(db CarDatabase) *CarRepository {
	return &CarRepository{DB: db}
}

func (r CarRepository) GetCars(id string) ([]model.EmployeeCars, error) {
	return r.DB.GetCars(id)
}
