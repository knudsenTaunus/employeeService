package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/database"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/handler"
	"github.com/janPhil/mySQLHTTPRestGolang/internal/server"
	"log"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("Failed to create database %s", err)
	}
	router := mux.NewRouter()
	employeeHandler := handler.New(db)
	s := server.NewServer(employeeHandler, router)
	/*if err != nil {
		log.Fatalf("Couldn't start the server: %v", err)
	}

	 */
	s.StartServer()

}
