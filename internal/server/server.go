package server

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/handler/authorization"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type handler interface {
	GetAll() http.HandlerFunc
	Get() http.HandlerFunc
	Add() http.HandlerFunc
	Remove() http.HandlerFunc
	GetCars() http.HandlerFunc
}

// employeeServer is a struct which contains all dependencies for this microservice
type employeeServer struct {
	employeeHandler handler
	router          *mux.Router
}

// New returns an instance of the employeeServer with all dependencies.
// The served routes can be found in the routes.go file
func New(h handler, r *mux.Router) *employeeServer {
	return &employeeServer{
		employeeHandler: h,
		router: r,
	}
}

func (s * employeeServer) SetRoutes() *employeeServer {
	getRouter := s.router.Methods(http.MethodGet).Subrouter()
	getRouter.Handle("/employees", s.employeeHandler.GetAll())
	getRouter.Handle("/employees", s.employeeHandler.GetAll()).Queries("limit", "{limit}")
	getRouter.Handle("/employee/{employee_number}", authorization.ValidateMiddleware(s.employeeHandler.Get()))
	getRouter.Handle("/employee/{id}/cars", authorization.ValidateMiddleware(s.employeeHandler.GetCars()))

	postRouter := s.router.Methods(http.MethodPost).Subrouter()
	postRouter.Handle("/employee", s.employeeHandler.Add())

	deleteRouter := s.router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.Handle("/employees/{id}", s.employeeHandler.Remove())
	return s
}

// StartServer creates an http.Server and a channel where it waits for SIGINT or SIGTERM.
func (s *employeeServer) StartServer(address string) {
	srv := &http.Server{
		Addr:              address,
		Handler:           s.router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: %s\n", err)
		}
	}()

	log.Printf("Server started")
	<- done
	log.Printf("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}