package storage

import (
	"context"
	"os"
	"sci-abo-go/models"
)

func InsertEventParticipant(model models.EventParticipant) error {

	collection := GetCollection(os.Getenv("EVENT_PARTICIPANT_COLLECTION"))

	_, err := collection.InsertOne(context.Background(), model)
	if err != nil{
		return err 
	}

	return nil

}
