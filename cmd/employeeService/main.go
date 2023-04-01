package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/db"
	"github.com/knudsenTaunus/employeeService/internal/handler"
	"github.com/knudsenTaunus/employeeService/internal/router"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	database, err := db.New(serviceConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create development database")
		return
	}

	employeeHandler := handler.NewEmployee(database, logger)
	carsHandler := handler.NewCar(database, logger)

	employeeRouter := router.New(employeeHandler, carsHandler)
	employeeRouter.SetRoutes()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", serviceConfig.Server.Host, serviceConfig.Server.Port),
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
