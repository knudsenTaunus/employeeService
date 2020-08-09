package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/handler/employee"
	"github.com/knudsenTaunus/employeeService/internal/server"
	"github.com/knudsenTaunus/employeeService/internal/storage/development"
	"github.com/knudsenTaunus/employeeService/internal/storage/production"
	"log"
)

func main() {

	employeeConfig, err := config.NewConfig("./config.yml")
	if err != nil {
		log.Fatalf("failed to configure service: %s", err)
	}
	
	var db employee.Storage
	
	switch employeeConfig.Environment {
	case "development":
		var err error
		db, err = development.New(employeeConfig)
		if err != nil {
			log.Fatal("failed to create development database")
		}
	case "production":
		var err error
		db, err = production.New(employeeConfig)
		if err != nil {
			log.Fatal("failed to create development database")
		}
	}

	router := mux.NewRouter()
	employeeHandler := employee.New(db)
	address := fmt.Sprintf("%s:%s", employeeConfig.Server.Host, employeeConfig.Server.Port)
	s := server.New(employeeHandler, router).SetRoutes()
	s.StartServer(address)
}