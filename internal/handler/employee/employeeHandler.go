package employee

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/types"
	"net/http"
)

type Storage interface {
	FindAllEmployees() (types.StorageEmployees, error)
	FindAllEmployeesLimit(limit string) (types.StorageEmployees, error)
	Find(id string) (*types.StorageEmployee, error)
	Add(employee *types.StorageEmployee) error
	Remove(id string) error
	GetCars(id string) ([]*types.EmployeeCars, error)
}

type handler struct {
	Database Storage
}

func New(db Storage) *handler {
	return &handler{
		Database: db,
	}
}

func (h *handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employeeNumber := mux.Vars(r)["employee_number"]
		employee, err := h.Database.Find(employeeNumber)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
		}
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(employee.ToHandlerEmployee())
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
		}
	}
}

func (h *handler) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employee := &types.HandlerEmployee{}
		err := employee.FromJSON(r.Body)
		w.Header().Add("Content-Type", "application/json")
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
		err = h.Database.Add(employee.ToStorageEmployee())
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *handler) Remove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employeeNumber := mux.Vars(r)["employee_number"]
		err := h.Database.Remove(employeeNumber)
		if err != nil {
			http.Error(w, http.StatusText(501), 501)
		}
	}
}

func (h *handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := r.URL.Query().Get("limit")
		if limit != "" {
			res, err := h.Database.FindAllEmployeesLimit(limit)
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
				return
			}
			w.Header().Add("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(res.ToHandlerEmployees())
			if err != nil {
				http.Error(w, http.StatusText(404), 404)
			}
			return
		}
		res, err := h.Database.FindAllEmployees()
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(res.ToHandlerEmployees())
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
		}
	}
}

func (h *handler) GetCars() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		cars, err := h.Database.GetCars(id)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(cars)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
	}
}
