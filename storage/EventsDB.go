package storage

import (
	"context"
	"fmt"
	"os"
	"sci-abo-go/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertEventDB(model *models.Event) (string, error) {

	// get the reset code collection
	collection := GetCollection(os.Getenv("EVENTS_COLLECTION"))

	// insert the reset code data to the MongoDB collection
	result, err := collection.InsertOne(context.Background(), model)
	if err != nil{
		return "-1", fmt.Errorf(err.Error())
	}

	// retrieve the id to send to the client 
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "-1", fmt.Errorf("error asserting InsertedID to ObjectID")
	}

	// convert the id to string
	return oid.Hex(), nil
}