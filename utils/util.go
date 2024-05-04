package utils

import (
	"fmt"
	"net/http"
	"sci-abo-go/models"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(user *models.User) error {

	hash_password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("field to hash password")
	}
	user.Password = string(hash_password)
	return nil
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
