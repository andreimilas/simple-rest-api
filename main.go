package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"sample-rest-api/app/api"
	"sample-rest-api/app/user"
	"sample-rest-api/config"
	"sample-rest-api/database"
)

// AddRoutes - add routes to router
func AddRoutes(router *mux.Router, apiHandler *api.Handler) {
	v1Router := router.PathPrefix("/v1").Subrouter()

	// Add Routes
	user.AddRoutes(v1Router, apiHandler)

	// Pretty print available routes to the CLI
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		if route.GetHandler() != nil {
			fmt.Println(methods, path)
		}
		return nil
	})
}

func main() {
	// Load config
	config.Load("config.yml")
	if config.Config == nil {
		os.Exit(2)
	}
	log.Println("Loaded configuration")

	// Prepare DSN for MySQL connection
	dbConfig := config.Config.Database
	mysqlConfig := mysql.Config{
		User:                 dbConfig.Username,
		Passwd:               dbConfig.Password,
		Net:                  "tcp",
		Addr:                 dbConfig.Hostname + ":" + strconv.Itoa(dbConfig.Port),
		DBName:               dbConfig.Name,
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	// Connect to MySQL
	dbHandle := database.Connect(mysqlConfig.FormatDSN())
	if dbHandle == nil {
		os.Exit(2)
	}
	defer dbHandle.Close()
	log.Println("MySQL connection established")

	// Initialize API handler
	apiHandler := api.Init(dbHandle)

	// Initialize router
	router := mux.NewRouter().StrictSlash(true)
	log.Println("Loading routes...")
	AddRoutes(router, apiHandler)

	// Start the HTTP server
	serverConfig := config.Config.Server
	log.Println("Running HTTP server on port " + serverConfig.Port + "...")
	log.Fatal(http.ListenAndServe(serverConfig.Hostname+":"+serverConfig.Port, router))
}
