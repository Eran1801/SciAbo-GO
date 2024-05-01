package utils

import (
	"net/http"
	"sci-abo-go/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	
)


func HashPassword(w http.ResponseWriter, user *models.User) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
	if err != nil{
		http.Error(w,"Field to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashPassword)
}


func GetObjectIdByStringId(id string) primitive.ObjectID {

	objID, err := primitive.ObjectIDFromHex(id) // id in the db is ObjectId and not string, needs to convert it.
	if err != nil {
		log.Println("Error converting string to ObjectId: ", err)
		return primitive.ObjectID{}
	}

	return objID
}