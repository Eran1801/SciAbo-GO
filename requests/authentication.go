package requests

import (
    "context"
    "encoding/json"
    "net/http"
    "sci-abo-go/models"
    "sci-abo-go/db"
    "sci-abo-go/config"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
    // Decode the JSON data from the request body into the user variable
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
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
