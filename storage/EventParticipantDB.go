package storage

import (
	"context"
	"fmt"
	"os"
	"sci-abo-go/models"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertEventParticipantDB(model models.EventParticipant) error {

	collection := GetCollection(os.Getenv("EVENT_PARTICIPANT_COLLECTION"))

	_, err := collection.InsertOne(context.Background(), model)
	if err != nil {
		return err
	}

	return nil

}

func DeleteEventParticipantByUserID(user_id string) error{

	collection := GetCollection(os.Getenv("EVENT_PARTICIPANT_COLLECTION"))

	filter := bson.M{"user_id": user_id}

	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if result.DeletedCount == 0{
		return fmt.Errorf("no Participant was found with %v id", user_id)
	}


	return nil
}
