package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type HandlerInterface interface {
	GetIndex() http.HandlerFunc
	GetAll() http.HandlerFunc
	Get() http.HandlerFunc
	Add() http.HandlerFunc
	Remove() http.HandlerFunc
}

// Server is a struct which contains all dependencies for this microservice
type Server struct {
	employeeHandler HandlerInterface
	router *mux.Router
}

// NewServer returns an instance of the server with all dependencies.
// The served routes can be found in the routes.go file
func NewServer(h HandlerInterface, r *mux.Router) *Server {
	s := &Server{
		employeeHandler: h,
		router: r,
	}
	fmt.Println("Server created")

	return s
}

// StartServer starts the new created Server
func (s *Server) StartServer() {
	fmt.Println("Server started")

	getRouter := s.router.Methods(http.MethodGet).Subrouter()
	getRouter.Handle("/all", s.employeeHandler.GetAll())
	getRouter.Handle("/{id}", s.employeeHandler.Get())


	postRouter := s.router.Methods(http.MethodPost).Subrouter()
	postRouter.Handle("/employee", s.employeeHandler.Add())

	deleteRouter := s.router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.Handle("/{id}", s.employeeHandler.Remove())


	err := http.ListenAndServe(":8080", s.router)
	if err != nil {

	}
	defer s.StopServer()
}

// StopServer stops the server and closes the connection to the storage
func (s *Server) StopServer() {
	fmt.Println("Shutting down server...")
	fmt.Println("closing storage...")
}
