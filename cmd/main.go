package main

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/handler/employee"
	"github.com/knudsenTaunus/employeeService/internal/server"
	"github.com/knudsenTaunus/employeeService/internal/storage"
)

var (
	environment string
	port string
)

func main() {
	configuration()
	flag.Parse()
	db := storage.New(&environment)
	router := mux.NewRouter()
	employeeHandler := employee.New(db)
	s := server.New(employeeHandler, router)
	s.StartServer(port)
}

func configuration() {
	flag.StringVar(&environment,"environment", "development", "environment to run the app in")
	flag.StringVar(&port, "port",":8080", "the port which the server runs on")
}
