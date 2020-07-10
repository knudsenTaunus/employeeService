package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/handler"
	"net/http"
)


// Server is a struct which contains all dependencies for this microservice
type Server struct {
	employeeHandler handler.HandlerInterface
	router *mux.Router
}

// NewServer returns an instance of the server with all dependencies.
// The served routes can be found in the routes.go file
func NewServer(h handler.HandlerInterface, r *mux.Router) *Server {
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
	s.router.Handle("/", s.employeeHandler.GetIndex())
	s.router.Handle("/all", s.employeeHandler.GetAll())
	err := http.ListenAndServe(":8080", s.router)
	if err != nil {

	}
	defer s.StopServer()
}

// StopServer stops the server and closes the connection to the database
func (s *Server) StopServer() {
	fmt.Println("Shutting down server...")
	fmt.Println("closing database...")
}
