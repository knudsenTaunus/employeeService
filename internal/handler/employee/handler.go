package employee

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/types"
	"net/http"
)

type Repository interface {
	FindAllEmployees() (types.StorageEmployees, error)
	FindAllEmployeesLimit(limit string) (types.StorageEmployees, error)
	Find(id string) (*types.StorageEmployee, error)
	Add(employee *types.StorageEmployee) error
	Remove(id string) error
}

type handler struct {
	Database Repository
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		employeeNumber := mux.Vars(r)["id"]
		if employeeNumber != "" {
			h.Get(employeeNumber, w)
		}

		h.GetAll(w, r)
	case http.MethodPost:
		h.Add(w, r)
		return
	case http.MethodDelete:
		h.Remove(w, r)
		return
	}
}

func NewHandler(db Repository) *handler {
	return &handler{
		Database: db,
	}
}

func (h *handler) Get(employeeNumber string, w http.ResponseWriter) {
	employee, err := h.Database.Find(employeeNumber)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(employee.ToHandlerEmployee())
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

}

func (h *handler) Add(w http.ResponseWriter, r *http.Request) {
	employee := &types.HandlerEmployee{}
	err := employee.FromJSON(r.Body)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	err = h.Database.Add(employee.ToStorageEmployee())
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *handler) Remove(w http.ResponseWriter, r *http.Request) {
	employeeNumber := mux.Vars(r)["employee_number"]
	err := h.Database.Remove(employeeNumber)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		res, err := h.Database.FindAllEmployeesLimit(limit)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(res.ToHandlerEmployees())
		if err != nil {
			http.Error(w, http.StatusText(404), http.StatusNotFound)
		}
		return
	}
	res, err := h.Database.FindAllEmployees()
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res.ToHandlerEmployees())
	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	}
}
