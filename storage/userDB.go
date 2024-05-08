package storage

import (
	"context"
	"fmt"
	"os"
	"sci-abo-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertUserDB(model *models.User) error {

	// get the User collection
	collection := GetCollection(os.Getenv("USER_COLLECTION"))

	// insert the user into the MongoDB collection
	_, err := collection.InsertOne(context.Background(), model)
	// checks for errors
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("email already exists")
		} else { // other errors
			return fmt.Errorf(err.Error())
		}
	}
	return nil
}


func GetUserByEmail(email string) (*models.User, error) {

	var user models.User

	collection := GetCollection(os.Getenv("USER_COLLECTION"))

	filter := bson.M{"email": email} // set the filter to retrieve data from the db

	err := collection.FindOne(context.TODO(), filter).Decode(&user) // check in the db and decode to the user variable
	if err != nil {
		if err == mongo.ErrNilDocument {
			return nil, nil // no docs was found
		}
		return nil, err
	}

	return &user, nil
}

