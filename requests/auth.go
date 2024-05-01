package requests

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"

	"sci-abo-go/config"
	"sci-abo-go/models"
	db "sci-abo-go/storage" // db is the alias
	"sci-abo-go/utils"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	// Decode the JSON data from the request body into the user variable
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// handel validation against db requirements
	err = models.ValidateUser(&user)
	if err != nil {
		// Handle validation errors
		errors := err.(validator.ValidationErrors)
		// Construct error message
		var errMsg string
		for _, e := range errors {
			errMsg += e.Field() + " is " + e.Tag() + "\n"
		}
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// hash the user password before saving it in the db
	utils.HashPassword(w, &user)

	// Get database and collection names from environment variables
	collection_name := config.GetEnvVar("USER_COLLECTION")

	// Get the MongoDB collection
	userCollection := db.GetCollection(collection_name)

	// Insert the user into the MongoDB collection
	_, err = userCollection.InsertOne(context.Background(), user)
	// checks for errors
	if err != nil {
		// checks if the email is already exist in the db
		if mongo.IsDuplicateKeyError(err) {
			http.Error(w, "Email already exists", http.StatusConflict)
			// other errors
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Send success response
	response := map[string]string{"message": "User registered successfully"}
	w.Header().Set("Content-Type", "application/json")

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
