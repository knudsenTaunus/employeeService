package router

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

// Employee is a struct which contains all dependencies for this microservice
type Employee struct {
	employeeHandler http.Handler
	carsHandler     http.Handler
	Router          *mux.Router
}

// New returns an instance of the employeeServer with all dependencies.
// The served routes can be found in the routes.go file
func New(eh http.Handler, ch http.Handler) Employee {

	return Employee{
		employeeHandler: eh,
		carsHandler:     ch,
		Router:          mux.NewRouter(),
	}
}

func (s Employee) SetRoutes() {
	s.Router.Handle("/employee", s.employeeHandler).Methods(http.MethodPost, http.MethodDelete)
	s.Router.Handle("/employee/{employee_number}", s.employeeHandler).Methods(http.MethodDelete)
	s.Router.Handle("/employees", s.employeeHandler).Methods(http.MethodGet)
	s.Router.Handle("/employees/{id}", s.employeeHandler).Methods(http.MethodGet)

	s.Router.Handle("/employee/{id}/cars", authorization.ValidateMiddleware(s.carsHandler))
}

// StartServer creates a http.Server and a channel where it waits for SIGINT or SIGTERM.
func (s Employee) StartServer(address string) {
	srv := &http.Server{
		Addr:    address,
		Handler: s.Router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started")
	<-done
	log.Printf("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited gracefully")
}
