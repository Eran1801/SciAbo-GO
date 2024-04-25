package main 

import (
    "github.com/gorilla/mux"
	"sci-abo-go/requests"
)

func initializerRoutes() *mux.Router{

	// creating a new instance of the router
	router := mux.NewRouter()

	// auth
	router.HandleFunc("/auth/register",requests.RegisterHandler).Methods("POST")

	// profile
	router.HandleFunc("/profile/upload_profile_image",requests.UploadProfilePictureHandler).Methods("POST")
	

	return router
}