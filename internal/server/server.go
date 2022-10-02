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

// EmployeeServer is a struct which contains all dependencies for this microservice
type EmployeeServer struct {
	employeeHandler http.Handler
	carsHandler     http.Handler
	router          *mux.Router
}

// New returns an instance of the employeeServer with all dependencies.
// The served routes can be found in the routes.go file
func New(eh http.Handler, ch http.Handler, r *mux.Router) EmployeeServer {
	return EmployeeServer{
		employeeHandler: eh,
		carsHandler:     ch,
		router:          r,
	}
}

func (s EmployeeServer) SetRoutes() {
	s.router.Handle("/employee", s.employeeHandler).Methods(http.MethodPost, http.MethodDelete)
	s.router.Handle("/employees", s.employeeHandler).Methods(http.MethodGet)
	s.router.Handle("/employees/{id}", s.employeeHandler).Methods(http.MethodGet)

	s.router.Handle("/employee/{id}/cars", authorization.ValidateMiddleware(s.carsHandler))
}

// StartServer creates a http.Server and a channel where it waits for SIGINT or SIGTERM.
func (s EmployeeServer) StartServer(address string) {
	srv := &http.Server{
		Addr:    address,
		Handler: s.router,
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
	log.Print("Server Exited Properly")
}
