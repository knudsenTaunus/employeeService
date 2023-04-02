package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/config"
)

// User is a struct which contains all dependencies for this microservice
type User struct {
	userHandler http.Handler
	Router      *mux.Router
}

// NewHTTP returns an instance of the userServer with all dependencies.
// The served routes can be found in the routes.go file
func NewHTTP(eh http.Handler, config *config.Config) *http.Server {
	user := User{
		userHandler: eh,
		Router:      MuxRouter(eh),
	}

	user.SetRoutes()

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler: user.Router,
	}
}

func MuxRouter(handler http.Handler) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/user", handler).Methods(http.MethodPost)
	r.Handle("/users", handler).Methods(http.MethodGet)
	r.Handle("/users/{id}", handler).Methods(http.MethodGet, http.MethodPatch, http.MethodDelete)
	return r
}
