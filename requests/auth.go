package requests

import (
	"encoding/json"
	"net/http"

	"sci-abo-go/models"
	"sci-abo-go/storage"
	"sci-abo-go/utils"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	// Decode the JSON data from the request body into the user variable
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ErrorResponse("Error in Decode the user request", w)
		return
	}

	// Handle validation against db requirements
	err = utils.ValidateDbRequirements(&user, w)
	if err != nil {
		ErrorResponse(err.Error(), w)
		return
	}

	// hash the user password before saving it in the db
	utils.HashPassword(w, &user)

	// save the user in the db
	err = storage.SaveUserInDB(&user, w)
	if err != nil {
		ErrorResponse(err.Error(), w)
		return
	}

	// Send success response
	SuccessResponse("User created successfully", string(user.Email), w)
}
