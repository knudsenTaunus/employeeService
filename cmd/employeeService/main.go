package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/db"
	"github.com/knudsenTaunus/employeeService/internal/handler/cars"
	"github.com/knudsenTaunus/employeeService/internal/handler/employee"
	"github.com/knudsenTaunus/employeeService/internal/repository"
	"github.com/knudsenTaunus/employeeService/internal/router"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	employeeRepository employee.Repository
	carRepository      cars.Repository
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func main() {

	logger := zerolog.New(os.Stdout)
	serviceConfig, configErr := config.NewConfig("./config.yml")
	if configErr != nil {
		logger.Fatal().Err(configErr).Msg("failed to configure service")
	}

	switch serviceConfig.Environment {
	case "development":
		database, err := db.NewSQLite(serviceConfig)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to create development database")
			return
		}
		employeeRepository = repository.NewEmployee(database)
		carRepository = repository.NewCar(database)
	case "production":
		database, err := db.NewMySQL(serviceConfig)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to create development database")
			return
		}
		employeeRepository = repository.NewEmployee(database)
		carRepository = repository.NewCar(database)
	}

	employeeHandler := employee.NewHandler(employeeRepository, logger)
	carsHandler := cars.NewHandler(carRepository, logger)
	address := fmt.Sprintf("%s:%s", serviceConfig.Server.Host, serviceConfig.Server.Port)
	employeeRouter := router.New(employeeHandler, carsHandler)

	employeeRouter.SetRoutes()

	srv := &http.Server{
		Addr:    address,
		Handler: employeeRouter.Router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("listen: %s\n")
		}
	}()

	logger.Printf("Server started")
	<-done
	logger.Printf("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server Shutdown Failed:%+v")
	}
	logger.Print("Server Exited gracefully")
}