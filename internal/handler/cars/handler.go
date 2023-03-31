package cars

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/model"
	"github.com/rs/zerolog"
	"net/http"
)

type Repository interface {
	GetCars(id string) ([]model.EmployeeCars, error)
}

type Handler struct {
	Database Repository
	logger   zerolog.Logger
}

func NewHandler(db Repository, logger zerolog.Logger) Handler {
	return Handler{
		Database: db,
		logger:   logger,
	}
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

func (h Handler) Get(id string, w http.ResponseWriter) {
	cars, err := h.Database.GetCars(id)
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
