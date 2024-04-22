package requests

import (
	"context"
	"encoding/json"
	"net/http"
	"sci-abo-go/config"
	"sci-abo-go/db"
	"sci-abo-go/models"

	"github.com/go-playground/validator/v10"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
    // Decode the JSON data from the request body into the user variable
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

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

    // Get database and collection names from environment variables
    db_name := config.GetEnvVar("DB_NAME")
    collection := config.GetEnvVar("USER_COLLECTION")

    // Get the MongoDB collection
    userCollection := db.GetCollection(db_name, collection)

    // Insert the user into the MongoDB collection
    _, err = userCollection.InsertOne(context.Background(), user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
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
