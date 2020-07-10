package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/handler"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/server"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/storage"
	"log"
)

func main() {
	db, err := storage.New()
	if err != nil {
		log.Fatalf("Failed to create storage %s", err)
	}
	router := mux.NewRouter()
	employeeHandler := handler.New(db)
	s := server.NewServer(employeeHandler, router)
	s.StartServer()
}
