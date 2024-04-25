package main

import (
	"sci-abo-go/requests"

	"github.com/gorilla/mux"
)

func initializerRoutes() *mux.Router {

	// creating a new instance of the router
	router := mux.NewRouter()

	// auth
	router.HandleFunc("/auth/register", requests.CreatingUser).Methods("POST")

	// profile
	router.HandleFunc("/profile/upload_profile_image", requests.UploadUserProfilePicture).Methods("POST")

	return router
}
