package storage

import (
	"context"
	"fmt"
	"os"
	"sci-abo-go/models"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	filter := bson.M{"email": strings.ToLower(email)} // set the filter to retrieve data from the db

	err := collection.FindOne(context.TODO(), filter).Decode(&user) // check in the db and decode to the user variable
	if err != nil {
		if err == mongo.ErrNilDocument {
			return nil, nil // no docs was found
		}
		return nil, err
	}

	return &user, nil
}


func GetUsersByIDs(user_ids []primitive.ObjectID) ([](models.User), error) {
	/*
	Given strings of id's of users
	this function returns all those users
	*/
	event_collection := GetCollection(os.Getenv("USER_COLLECTION"))
    var users []models.User
    filter := bson.M{"_id": bson.M{"$in": user_ids}}
    cursor, err := event_collection.Find(context.Background(), filter)
    if err != nil {
        return nil, err
    }

    err = cursor.All(context.Background(), &users)
	if err != nil {
	 	return nil, err
	}

	return users, nil

}

func DeleteUser(user_id primitive.ObjectID) error { 

	collection := GetCollection(os.Getenv("USER_COLLECTION"))
	
	filter := bson.M{"_id": user_id}

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil{
		return fmt.Errorf(err.Error())
	}
	return nil

}