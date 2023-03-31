package employee

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog"
	"log"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/model"
)

type Repository interface {
	FindAllEmployees() (model.StorageEmployees, error)
	FindAllEmployeesLimit(limit string) (model.StorageEmployees, error)
	FindEmployee(id string) (model.StorageEmployee, error)
	AddEmployee(employee model.StorageEmployee) error
	RemoveEmployee(id string) error
}

type Handler struct {
	database Repository
	logger   zerolog.Logger
}

func NewHandler(db Repository, logger zerolog.Logger) Handler {
	return Handler{
		database: db,
		logger:   logger,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		employeeNumber := mux.Vars(r)["id"]
		if employeeNumber != "" {
			h.Get(employeeNumber, w)
			return
		}

		h.GetAll(w, r)
		return
	case http.MethodPost:
		h.Add(w, r)
		return
	case http.MethodDelete:
		employeeNumber := mux.Vars(r)["employee_number"]
		if employeeNumber == "" {
			h.logger.Error().Msg("no employee number provided")
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}
		h.Remove(w, employeeNumber)
		return
	}
}

func (h Handler) Get(employeeNumber string, w http.ResponseWriter) {
	employee, err := h.database.FindEmployee(employeeNumber)
	if err != nil {
		if errors.Is(err, model.NotFoundError) {
			h.logger.Error().Err(err).Msg("failed to find employee")
			http.Error(w, http.StatusText(404), http.StatusNotFound)
			return
		}

		h.logger.Error().Err(err).Msg(err.Error())
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(employee.ToHandlerEmployee())
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}

func (h Handler) Add(w http.ResponseWriter, r *http.Request) {
	employee := model.HandlerEmployee{}
	err := json.NewDecoder(r.Body).Decode(&employee)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to add employee")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	storageEmployee := employee.ToStorageEmployee()
	err = h.database.AddEmployee(storageEmployee)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h Handler) Remove(w http.ResponseWriter, employeeNumber string) {
	err := h.database.RemoveEmployee(employeeNumber)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to remove employee")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	h.logger.Info().Msg("removed employee")
	w.WriteHeader(http.StatusAccepted)

}

func (h Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		res, err := h.database.FindAllEmployeesLimit(limit)
		if err != nil {
			h.logger.Error().Err(err).Msg("failed to get all employees")
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(res.ToHandlerEmployees())
		if err != nil {
			h.logger.Error().Err(err).Msg("failed to json marshal employees")
			http.Error(w, http.StatusText(404), http.StatusNotFound)
		}
		return
	}
	res, err := h.database.FindAllEmployees()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get all employees")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	handlerEmployees := res.ToHandlerEmployees()
	err = json.NewEncoder(w).Encode(handlerEmployees)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to json marshal employees")
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	}
}
