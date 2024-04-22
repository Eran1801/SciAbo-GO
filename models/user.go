package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user in the system
type User struct {
    ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    FirstName    string `json:"first_name"`
    LastName     string `json:"last_name"`
    Email        string `json:"email"`
    Password     string `json:"password"`
    ConfirmPassword string `json:"confirm_password"`
}

