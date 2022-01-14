package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/handler/cars"
	"github.com/knudsenTaunus/employeeService/internal/handler/employee"
	"github.com/knudsenTaunus/employeeService/internal/server"
	"github.com/knudsenTaunus/employeeService/internal/store"
	"log"
)

func main() {

	employeeConfig, err := config.NewConfig("./config.yml")
	if err != nil {
		log.Fatalf("failed to configure service: %s", err)
	}

	var repository *store.Repository

	switch employeeConfig.Environment {
	case "development":
		var err error
		db, err := store.NewSQLite(employeeConfig)
		if err != nil {
			log.Fatal("failed to create development database")
		}
		repository = store.NewRepository(db)
	case "production":
		var err error
		db, err := store.NewMySQL(employeeConfig)
		if err != nil {
			log.Fatal("failed to create development database")
		}
		repository.DB = store.NewRepository(db)
	}

	router := mux.NewRouter()
	employeeHandler := employee.NewHandler(repository)
	carsHandler := cars.NewHandler(repository)
	address := fmt.Sprintf("%s:%s", employeeConfig.Server.Host, employeeConfig.Server.Port)
	s := server.New(employeeHandler, carsHandler, router).SetRoutes()
	s.StartServer(address)
}
