package requests

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"sci-abo-go/models"
	db "sci-abo-go/storage" // db is the alias
	"sci-abo-go/utils"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	// Decode the JSON data from the request body into the user variable
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.ErrorResponse("Error in Decode the user request",nil,w)
	}

	// handel validation against db requirements
	utils.ValidateDbRequirements(&user, w)

	// hash the user password before saving it in the db
	utils.HashPassword(w, &user)

	// Get the User collection
	user_collection := db.GetUserCollection()

	// Insert the user into the MongoDB collection
	_, err = user_collection.InsertOne(context.Background(), user)
	// checks for errors
	if err != nil {
		// checks if the email is already exist in the db
		if mongo.IsDuplicateKeyError(err) {
			utils.ErrorResponse("Email already exists", nil, w)
		} else { // other errors
			utils.ErrorResponse(err.Error(), nil, w)
		}
		return
	}

	// Send success response
	utils.SuccessResponse("User created successfully", nil, w)
}
