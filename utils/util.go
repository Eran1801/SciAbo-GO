package utils

import (
	"fmt"
	"log"
	"net/http"
	"sci-abo-go/models"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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

func ValidateDbRequirements(user *models.User, w http.ResponseWriter) error {
    err := models.ValidateUser(user)
    if err != nil {
        // Handle validation errors
        errors := err.(validator.ValidationErrors)
        // Construct error message
        var errMsg string
        for _, e := range errors {
            errMsg += e.Field() + " is " + e.Tag() + "; "
			if len(errMsg) > 0 {
				return fmt.Errorf(errMsg)
			}
        }

    }
    return nil
}