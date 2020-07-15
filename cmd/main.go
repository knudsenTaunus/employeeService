package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/handler/employee"
	"github.com/knudsenTaunus/employeeService/internal/server"
	"github.com/knudsenTaunus/employeeService/internal/storage"
	"log"
)

func main() {
	db, err := storage.New()
	if err != nil {
		log.Fatalf("Failed to create storage %s", err)
	}
	router := mux.NewRouter()
	employeeHandler := employee.New(db)
	s := server.New(employeeHandler, router)
	s.StartServer()
}
