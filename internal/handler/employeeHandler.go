package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/types"
	"net/http"
)

type EmployeeHandler struct {
	db *sql.DB
}

func New(db *sql.DB) *EmployeeHandler {
	h := &EmployeeHandler{
		db: db,
	}
	return h
}

func (h *EmployeeHandler) GetIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode("Welcome")
		if err != nil {
			http.Error(w, http.StatusText(404),404)
		}
	}
}

func (h *EmployeeHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := types.AllEmployees(h.db)
		if err != nil {
			http.Error(w, http.StatusText(500),500)
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
				http.Error(w, http.StatusText(404),404)
			}
	}
}