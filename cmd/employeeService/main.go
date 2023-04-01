package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	userpb "github.com/knudsenTaunus/employeeService/gen/go/proto/user/v1"
	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/db"
	"github.com/knudsenTaunus/employeeService/internal/handler"
	"github.com/knudsenTaunus/employeeService/internal/model"
	"github.com/knudsenTaunus/employeeService/internal/router"
	"github.com/knudsenTaunus/employeeService/internal/server/protobuf"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"log"
	"net"
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

	userChan := make(chan model.Employee)

	employeeHandler := handler.NewEmployee(database, userChan, logger)
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
	logger.Info().Msg("HTTP Server started")

	userRPCServer := protobuf.NewServer(userChan, logger)
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, userRPCServer)

	lis, err := net.Listen("tcp", "localhost:9879")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal().Err(err).Msg("failed to start grpc server")
		}
	}()
	logger.Info().Msg("gRPC Server started")
	go userRPCServer.Invoke()

	<-done
	logger.Info().Msg("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server Shutdown Failed:%+v")
	}
	logger.Print("Server Exited gracefully")
}
