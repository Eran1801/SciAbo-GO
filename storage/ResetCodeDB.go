package storage

import (
	"context"
	"fmt"
	"os"
	"sci-abo-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertResetCodeDB(model *models.ResetCode) (string, error) {

	// get the reset code collection
	collection := GetCollection(os.Getenv("RESET_CODE_COLLECTION"))

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

func GetResetCodeByID(id primitive.ObjectID) (*models.ResetCode, error) {

	var model models.ResetCode

	collection := GetCollection(os.Getenv("RESET_CODE_COLLECTION"))

	filter := bson.M{"_id": id} // based on what to search in the db

	err := collection.FindOne(context.TODO(),filter).Decode(&model)
	if err != nil{
		if err == mongo.ErrNilDocument {
			return nil, nil
		}
		return nil, err
	}

	return &model,nil

}
