package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/config"
	"net/http"
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
		Router:      mux.NewRouter(),
	}

	user.SetRoutes()

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler: user.Router,
	}
}
