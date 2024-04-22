package requests 

import (
    "github.com/gorilla/mux"
)

func InitializerRoutes() *mux.Router{

	// creating a new instance of the router
	router := mux.NewRouter()

	// define my routes
	router.HandleFunc("/auth/register",registerHandler).Methods("POST")

	return router
}