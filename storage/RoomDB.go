package storage


// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"sci-abo-go/models"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// func InsertRoom(room *models.Room) (string, error) {

// 	collection := GetCollection(os.Getenv("ROOM_COLLECTION"))

// 	result, err := collection.InsertOne(context.Background(), room)
// 	if err != nil {
// 		return "-1", err
// 	}

// 	object_id, ok := result.InsertedID.(primitive.ObjectID)
// 	if !ok {
// 		return "-1", fmt.Errorf("error InsertID")
// 	}

// 	return object_id.Hex(), nil

// }

// func FetchRoomByID(room_id primitive.ObjectID) (*models.Room, error) {

// 	var room models.Room

// 	collection := GetCollection(os.Getenv("ROOM_COLLECTION"))

// 	filter := bson.M{"_id": room_id}
	
// 	err := collection.FindOne(context.Background(), filter).Decode(&room)
// 	if err != nil {
// 		if err == mongo.ErrNilDocument {
// 			return nil, err
// 		}
// 		return nil, err
// 	}

// 	return &room, nil

// }
