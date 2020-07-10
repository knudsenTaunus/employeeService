package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/types"
	"net/http"
)

type Database interface {
	FindAllEmployees() ([]*types.Employee, error)
	Find(id string) (*types.Employee, error)
	Add(employee *types.Employee) error
}

type EmployeeHandler struct {
	Database Database
}

func New(db Database) *EmployeeHandler {
	return &EmployeeHandler{
		Database: db,
	}
}

func (h *EmployeeHandler) GetIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode("Welcome")
		if err != nil {
			http.Error(w, http.StatusText(404),404)
		}
	}
}

func (h *EmployeeHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		idString := params["id"]
		employee, err := h.Database.Find(idString)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(500),500)
		}
		if err != nil {
			http.Error(w, http.StatusText(404),404)
		}
		err = json.NewEncoder(w).Encode(employee)
		if err != nil {
			http.Error(w, http.StatusText(404),404)
		}
	}
}

func (h *EmployeeHandler) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employee := &types.Employee{}
		err := employee.FromJSON(r.Body)
		err = h.Database.Add(employee)
		if err != nil {
			http.Error(w, http.StatusText(501),501)
		}
	}
}

func (h *EmployeeHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := h.Database.FindAllEmployees()
		if err != nil {
			http.Error(w, http.StatusText(500),500)
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
				http.Error(w, http.StatusText(404),404)
			}
	}
}