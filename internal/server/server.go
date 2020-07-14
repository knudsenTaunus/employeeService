package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/handler/authorization"
	"net/http"
)

type handler interface {
	GetAll() http.HandlerFunc
	Get() http.HandlerFunc
	Add() http.HandlerFunc
	Remove() http.HandlerFunc
	GetCars() http.HandlerFunc
}

// server is a struct which contains all dependencies for this microservice
type server struct {
	employeeHandler handler
	router          *mux.Router
}

// New returns an instance of the server with all dependencies.
// The served routes can be found in the routes.go file
func New(h handler, r *mux.Router) *server {
	s := &server{
		employeeHandler: h,
		router: r,
	}
	fmt.Println("server created")

	return s
}

// StartServer starts the new created server
func (s *server) StartServer() {
	fmt.Println("server started")
	getRouter := s.router.Methods(http.MethodGet).Subrouter()
	getRouter.Handle("/employees", s.employeeHandler.GetAll())
	getRouter.Handle("/employees", s.employeeHandler.GetAll()).Queries("limit", "{limit}")
	getRouter.Handle("/employee/{id}", authorization.ValidateMiddleware(s.employeeHandler.Get()))
	getRouter.Handle("/employees/{id}/cars", authorization.ValidateMiddleware(s.employeeHandler.GetCars()))

	postRouter := s.router.Methods(http.MethodPost).Subrouter()
	postRouter.Handle("/employee", s.employeeHandler.Add())

	deleteRouter := s.router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.Handle("/employees/{id}", s.employeeHandler.Remove())


	err := http.ListenAndServe(":8080", s.router)
	if err != nil {

	}
	defer s.StopServer()
}

// StopServer stops the server and closes the connection to the storage
func (s *server) StopServer() {
	fmt.Println("Shutting down server...")
	fmt.Println("closing storage...")
}
