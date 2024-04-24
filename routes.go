package main 

import (
    "github.com/gorilla/mux"
	"sci-abo-go/requests"
)

func initializerRoutes() *mux.Router{

	// creating a new instance of the router
	router := mux.NewRouter()

	// define my routes
	router.HandleFunc("/auth/register",requests.RegisterHandler).Methods("POST")

	return router
}