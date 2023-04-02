package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/knudsenTaunus/employeeService/internal/crypt"

	_ "github.com/go-sql-driver/mysql"
	userpb "github.com/knudsenTaunus/employeeService/generated/go/proto/user/v1"
	"github.com/knudsenTaunus/employeeService/internal/config"
	"github.com/knudsenTaunus/employeeService/internal/db"
	"github.com/knudsenTaunus/employeeService/internal/handler"
	"github.com/knudsenTaunus/employeeService/internal/model"
	"github.com/knudsenTaunus/employeeService/internal/server/proto"
	"github.com/knudsenTaunus/employeeService/internal/server/rest"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
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

	if serviceConfig.Environment == "development" {
		logger = logger.With().Caller().Logger()
	}

	cryptor := crypt.New(serviceConfig.Crypt.Secret)
	database, err := db.NewMySQL(serviceConfig, cryptor)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create development database")
		return
	}

	userChan := make(chan model.User)
	user := handler.NewUser(database, userChan, logger)
	srv := rest.NewHTTP(user, serviceConfig)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("listen: %s\n")
		}
	}()

	logger.Info().Msg("HTTP Server started")

	userRPCServer := proto.NewServer(userChan, logger)
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, userRPCServer)

	grpcAddress := fmt.Sprintf("%s:%s", serviceConfig.Grpc.Host, serviceConfig.Grpc.Port)
	lis, err := net.Listen("tcp", grpcAddress)
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
