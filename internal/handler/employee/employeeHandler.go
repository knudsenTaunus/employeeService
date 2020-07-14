package employee

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/types"
	"net/http"
)

type database interface {
	FindAllEmployees() ([]*types.Employee, error)
	FindAllEmployeesLimit(limit string) ([]*types.Employee, error)
	Find(id string) (*types.Employee, error)
	Add(employee *types.Employee) error
	Remove(id string) error
	GetCars(id string) ([]*types.EmployeeCars, error)
}

type handler struct {
	Database database
}

func New(db database) *handler {
	return &handler{
		Database: db,
	}
}

func (h *handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		employee, err := h.Database.Find(id)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(500), 500)
		}
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
		}
		err = json.NewEncoder(w).Encode(employee)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
		}
	}
}

func (h *handler) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employee := &types.Employee{}
		err := employee.FromJSON(r.Body)
		err = h.Database.Add(employee)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func (h *handler) Remove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		err := h.Database.Remove(id)
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
			}
			err = json.NewEncoder(w).Encode(res)
			if err != nil {
				http.Error(w, http.StatusText(404), 404)
			}
			return
		}
		res, err := h.Database.FindAllEmployees()
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
		err = json.NewEncoder(w).Encode(res)
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
		err = json.NewEncoder(w).Encode(cars)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
	}
}
