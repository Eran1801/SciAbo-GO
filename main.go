package main

import (
	"log"
	"net/http"
	"sci-abo-go/storage"
	"github.com/joho/godotenv"

)

func main(){

	// load env file
	if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

	// db connection and creating indexes
	storage.InitializeDB()

	// init routes
	router := InitializerRoutes()

	// start the HTTP server using the router
	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080",router)
	if err != nil {
		log.Fatalf("Server failed to start %v ", err)
	}

}

