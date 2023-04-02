package rest

import "net/http"

func (s User) SetRoutes() {
	s.Router.Handle("/user", s.userHandler).Methods(http.MethodPost)
	s.Router.Handle("/users", s.userHandler).Methods(http.MethodGet)
	s.Router.Handle("/users/{id}", s.userHandler).Methods(http.MethodGet, http.MethodPatch, http.MethodDelete)
}
