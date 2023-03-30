package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/db"
	"github.com/knudsenTaunus/employeeService/internal/handler/cars"
	"github.com/knudsenTaunus/employeeService/internal/handler/employee"
	"github.com/knudsenTaunus/employeeService/internal/repository"
	"github.com/knudsenTaunus/employeeService/internal/server"
	"log"
)

var (
	employeeRepository employee.Repository
	carRepository      cars.Repository
)

func main() {

	serviceConfig, configErr := config.NewConfig("./config.yml")
	if configErr != nil {
		log.Fatalf("failed to configure service: %s", configErr)
	}

	switch serviceConfig.Environment {
	case "development":
		db, err := db.NewSQLite(serviceConfig)
		if err != nil {
			log.Fatal("failed to create development database")
		}
		employeeRepository = repository.NewEmployee(db)
		carRepository = repository.NewCar(db)
	case "production":
		db, err := db.NewMySQL(serviceConfig)
		if err != nil {
			log.Fatal("failed to create development database")
		}
		employeeRepository = repository.NewEmployee(db)
		carRepository = repository.NewCar(db)
	}

	router := mux.NewRouter()
	employeeHandler := employee.NewHandler(employeeRepository)
	carsHandler := cars.NewHandler(carRepository)
	address := fmt.Sprintf("%s:%s", serviceConfig.Server.Host, serviceConfig.Server.Port)
	s := server.New(employeeHandler, carsHandler, router)
	s.SetRoutes()
	s.StartServer(address)
}
