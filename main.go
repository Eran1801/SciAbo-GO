package main

import (
	"log"
	"net/http"
	"sci-abo-go/config"
	"sci-abo-go/db"
)

func main(){

	// load env vars into the 
	config.LoadEnv()

	// db connection
	db.InitializeDB()

	// init routes
	router := initializerRoutes()

	// start the HTTP server using the router
	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080",router)
	if err != nil {
		log.Fatalf("Server failed to start %v ", err)
	}

}

