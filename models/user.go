package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/go-playground/validator/v10"
    "regexp"

)

type User struct {
    ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    FirstName           string `json:"first_name" validate:"required"` 
    LastName            string `json:"last_name" validate:"required"`
    Email               string `json:"email" validate:"required,email"`
    Password            string `json:"password" validate:"required,passwordPattern"`
    ConfirmPassword     string `json:"confirm_password" validate:"required,eqfield=Password" bson:"-"` // bson tell mongo skip this in saving to the db

}

// define a custom validator instance
var validate *validator.Validate

// init() function initialize automatic when our program run
func init(){
    validate = validator.New()

    // add custom validate for password pattern
    validate.RegisterValidation("passwordPattern", ValidatePassword)
}


func ValidateUser(user *User) error {
    err := validate.Struct(user)
    if err != nil {
        return err
    }
    return nil
}

func ValidatePassword(f1 validator.FieldLevel) bool {
    password := f1.Field().String()
    // Compile separate regex for each requirement
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSymbol := regexp.MustCompile(`[!@#$%^&*()]`).MatchString(password)
    hasMinLength := len(password) >= 8

    // Check all conditions are met
    return hasLower && hasUpper && hasDigit && hasSymbol && hasMinLength
}

