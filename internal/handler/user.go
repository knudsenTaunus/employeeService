package handler

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/rs/zerolog"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/model"
)

type UserDatabase interface {
	GetAll() ([]model.User, error)
	GetPaginatedAndFiltered(page, pageSize int, filter string) ([]model.User, error)
	Get(id string) (model.User, error)
	Create(User model.User) (model.User, error)
	Delete(id string) error
	Update(User model.User) (model.User, error)
}

type User struct {
	database UserDatabase
	userChan chan model.User
	logger   zerolog.Logger
}

func NewUser(db UserDatabase, userChan chan model.User, logger zerolog.Logger) User {
	return User{
		database: db,
		userChan: userChan,
		logger:   logger,
	}
}

func (h User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := mux.Vars(r)["id"]
		if id != "" {
			h.Get(id, w)
			return
		}

		h.GetAll(w, r)
		return
	case http.MethodPost:
		h.Create(w, r)
		return
	case http.MethodPatch:
		h.Update(w, r)
		return
	case http.MethodDelete:
		id := mux.Vars(r)["id"]
		if id == "" {
			h.logger.Error().Msg("request is missing id")
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}
		h.Delete(w, id)
		return
	}
}

func (h User) Get(id string, w http.ResponseWriter) {
	user, err := h.database.Get(id)
	if err != nil {
		if errors.Is(err, model.NotFoundError) {
			h.logger.Error().Err(err).Send()
			http.Error(w, http.StatusText(404), http.StatusNotFound)
			return
		}

		h.logger.Error().Err(err).Msg(err.Error())
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	return
}

func (h User) GetAll(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	if len(queryParams) == 0 {
		users, err := h.database.GetAll()
		if err != nil {
			h.logger.Error().Err(err).Send()
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			h.logger.Error().Err(err).Msg("failed to json marshal Users")
			http.Error(w, http.StatusText(404), http.StatusNotFound)
		}
		return
	}

	page, err := strconv.Atoi(queryParams.Get("page"))
	pageSize, err := strconv.Atoi(queryParams.Get("pageSize"))
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get query params")
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	filter := queryParams.Get("filter")
	users, err := h.database.GetPaginatedAndFiltered(page, pageSize, filter)
	if err != nil {
		h.logger.Error().Err(err).Send()
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to json marshal Users")
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	}
	return

}

func (h User) Create(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		h.logger.Error().Err(err).Send()
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	result, err := h.database.Create(user)
	if err != nil {
		if errors.Is(err, model.DuplicateNickError) || errors.Is(err, model.DuplicateMailError) {
			h.logger.Error().Err(err).Send()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.logger.Error().Err(err).Send()
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to json marshal Users")
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	}

	return
}

func (h User) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user := model.User{
		ID: id,
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.logger.Error().Err(err).Msgf("failed to unmarshal user %s", id)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}
	updatedUser, updateErr := h.database.Update(user)
	if updateErr != nil {
		if errors.Is(updateErr, model.DuplicateNickError) || errors.Is(updateErr, model.DuplicateMailError) {
			h.logger.Error().Err(updateErr).Send()
			http.Error(w, updateErr.Error(), http.StatusBadRequest)
			return
		}
		h.logger.Error().Err(updateErr).Send()
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	go func() {
		h.userChan <- updatedUser
	}()

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedUser)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	return
}

func (h User) Delete(w http.ResponseWriter, id string) {
	err := h.database.Delete(id)
	if err != nil {
		h.logger.Error().Err(err).Send()
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	h.logger.Info().Msg("removed User")
	w.WriteHeader(http.StatusAccepted)

}
