package handler

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog"
	"io"
	"log"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/model"
)

type EmployeeDatabase interface {
	FindAllEmployees() ([]model.Employee, error)
	FindAllEmployeesLimit(limit string) ([]model.Employee, error)
	FindEmployee(id string) (model.Employee, error)
	AddEmployee(employee model.Employee) error
	RemoveEmployee(id string) error
	UpdateEmployee(employee model.Employee) error
}

type Employee struct {
	database EmployeeDatabase
	userChan chan model.Employee
	logger   zerolog.Logger
}

func NewEmployee(db EmployeeDatabase, userChan chan model.Employee, logger zerolog.Logger) Employee {
	return Employee{
		database: db,
		userChan: userChan,
		logger:   logger,
	}
}

func (h Employee) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	case http.MethodPatch:
		h.Update(w, r)
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

func (h Employee) Get(employeeNumber string, w http.ResponseWriter) {
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
	err = json.NewEncoder(w).Encode(employee)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}

func (h Employee) Add(w http.ResponseWriter, r *http.Request) {
	employee := model.Employee{}
	err := json.NewDecoder(r.Body).Decode(&employee)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to add employee")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	err = h.database.AddEmployee(employee)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}

func (h Employee) Update(w http.ResponseWriter, r *http.Request) {
	employeeNumber := mux.Vars(r)["id"]
	employee, findErr := h.database.FindEmployee(employeeNumber)
	if findErr != nil {
		h.logger.Error().Err(findErr).Msgf("failed to find employee %s", employeeNumber)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}

	body, bodyErr := io.ReadAll(r.Body)
	if bodyErr != nil {
		h.logger.Error().Err(bodyErr).Msg("failed to read body")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	patchErr := employee.Patch(body)
	if patchErr != nil {
		h.logger.Error().Err(patchErr).Msg("failed to patch employee")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

	updateErr := h.database.UpdateEmployee(employee)
	if updateErr != nil {
		h.logger.Error().Err(updateErr).Msg("failed to save updated employee")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

	go func() {
		h.userChan <- employee
	}()

	w.WriteHeader(http.StatusAccepted)
	return
}

func (h Employee) Remove(w http.ResponseWriter, employeeNumber string) {
	err := h.database.RemoveEmployee(employeeNumber)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to remove employee")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	h.logger.Info().Msg("removed employee")
	w.WriteHeader(http.StatusAccepted)

}

func (h Employee) GetAll(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		employees, err := h.database.FindAllEmployeesLimit(limit)
		if err != nil {
			h.logger.Error().Err(err).Msg("failed to get all employees")
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(employees)
		if err != nil {
			h.logger.Error().Err(err).Msg("failed to json marshal employees")
			http.Error(w, http.StatusText(404), http.StatusNotFound)
		}
		return
	}

	employees, err := h.database.FindAllEmployees()
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get all employees")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	EmployeeEmployees := employees
	err = json.NewEncoder(w).Encode(EmployeeEmployees)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to json marshal employees")
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	}
}
