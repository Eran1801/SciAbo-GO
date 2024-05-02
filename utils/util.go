package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"sci-abo-go/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-playground/validator/v10"

)


func HashPassword(w http.ResponseWriter, user *models.User) {

	hash_password, err := bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
	if err != nil{
		http.Error(w,"Field to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hash_password)
}


func GetObjectIdByStringId(id string) primitive.ObjectID {

	obj_id, err := primitive.ObjectIDFromHex(id) // id in the db is ObjectId and not string, needs to convert it.
	if err != nil {
		log.Println("Error converting string to ObjectId: ", err)
		return primitive.ObjectID{}
	}

	return obj_id
}

func SuccessResponse(message string,data *string, w http.ResponseWriter) {

	response := make(map[string]string)
	response["message"] = message

	if data != nil {
		response["data"] = *data
	}

	w.Header().Set("Content-Type", "application/json")

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func ErrorResponse(message string, data *string, w http.ResponseWriter) {

	response := make(map[string]string)
	response["message"] = message

	if data != nil {
		response["data"] = *data
	}

	w.Header().Set("Content-Type", "application/json")

	json_response, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write(json_response)

}

func ValidateDbRequirements(user *models.User, w http.ResponseWriter) {

	err := models.ValidateUser(user)
	if err != nil {
		// Handle validation errors
		errors := err.(validator.ValidationErrors)
		// Construct error message
		var errMsg string
		for _, e := range errors {
			errMsg += e.Field() + " is " + e.Tag() + "\n"
		}
		ErrorResponse(errMsg, nil, w)
		return
	}

}