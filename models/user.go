package models

import(
    "github.com/go-playground/validator/v10"
    "go.mongodb.org/mongo-driver/bson/primitive"

)

// Define a custom validator instance
var validate *validator.Validate

func init(){
    validate = validator.New()
}

// User represents a user in the system
type User struct {
    ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    FirstName    string `json:"first_name" validate:"required"` 
    LastName     string `json:"last_name" validate:"required"`
    Email        string `json:"email" validate:"required,email"`
    Password     string `json:"password" validate:"required"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

func ValidateUser(user *User) error {
    err := validate.Struct(user)
    if err != nil {
        return err
    }
    return nil
}

