package models


// User represents a user in the system
type User struct {
    ID           uint   `json:"id"`
    FirstName    string `json:"first_name"`
    LastName     string `json:"last_name"`
    Email        string `json:"email"`
    Password     string `json:"password"`
    ConfirmPassword string `json:"confirm_password"`
}

