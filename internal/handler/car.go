package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/model"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
)

type CarDatabase interface {
	GetCars(id string) ([]model.Car, error)
	AddCar(car model.StorageCar) error
}

type Car struct {
	db     CarDatabase
	logger zerolog.Logger
}

func NewCar(db CarDatabase, logger zerolog.Logger) Car {
	return Car{
		db:     db,
		logger: logger,
	}
}

func (h Car) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		employeeNumber := mux.Vars(r)["id"]
		if employeeNumber != "" {
			h.Get(employeeNumber, w)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	case http.MethodPost:
		employeeNumber := mux.Vars(r)["id"]
		if employeeNumber != "" {
			h.Add(employeeNumber, w, r)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h Car) Get(id string, w http.ResponseWriter) {
	cars, err := h.db.GetCars(id)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get cars")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cars)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}

func (h Car) Add(id string, w http.ResponseWriter, r *http.Request) {
	car := model.StorageCar{}
	err := json.NewDecoder(r.Body).Decode(&car)

	employeeNumber, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to add car")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	car.EmployeeNumber = employeeNumber
	err = h.db.AddCar(car)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to add car")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}
