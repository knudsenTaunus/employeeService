package cars

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/model"
	"net/http"
)

type Repository interface {
	GetCars(id string) ([]model.EmployeeCars, error)
}

type Handler struct {
	Database Repository
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		employeeNumber := mux.Vars(r)["id"]
		if employeeNumber != "" {
			h.Get(employeeNumber, w)
		}
	}
}

func NewHandler(db Repository) Handler {
	return Handler{
		Database: db,
	}
}

func (h Handler) Get(id string, w http.ResponseWriter) {
	cars, err := h.Database.GetCars(id)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cars)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}
