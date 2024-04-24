package utils

import (
	"net/http"
	"sci-abo-go/models"
	"golang.org/x/crypto/bcrypt"
)


func HashPassword(w http.ResponseWriter, user *models.User) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
	if err != nil{
		http.Error(w,"Field to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashPassword)
}