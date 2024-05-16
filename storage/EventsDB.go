package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"sci-abo-go/models"
	"sci-abo-go/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertEventDB(model *models.Event) (string, error) {

	// get the event collection
	collection := GetCollection(os.Getenv("EVENTS_COLLECTION"))

	// insert the reset code data to the MongoDB collection
	result, err := collection.InsertOne(context.Background(), model)
	if err != nil {
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


func AddParticipantToEvent(event_id primitive.ObjectID, participant_id string) error {

	// get event collection
	collection := GetCollection(os.Getenv("EVENTS_COLLECTION"))

	// filter also verify that the participant_id in not already in the event participants
	filter := bson.M{"_id": event_id, "participants": bson.M{"$ne": participant_id}}
	update := bson.M{"$push": bson.M{"participants": participant_id}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user already in the event participants")
	}
	return nil

}


func AddEventIdToUserEvents(collection_name string, user_id primitive.ObjectID, event_id string) error {

	// get the user collection
	collection := GetCollection(collection_name)

	filter := bson.M{"_id": user_id} // find user by id
	update := bson.M{"$push": bson.M{"joined_event_ids": event_id}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil

}


func FetchUserEvents(events_ids []primitive.ObjectID) []models.Event {
	/*
		Given a list of id's of events, this function returns all the events
		that the user is join to
	*/

	// gets the event collection
	event_collection := GetCollection(os.Getenv("EVENTS_COLLECTION"))

    var events []models.Event // init an empty array of Events

	// set the filter to be the array of events ids
    filter := bson.M{"_id": bson.M{"$in": events_ids}} 

    cursor, err := event_collection.Find(context.Background(), filter)
    if err != nil {
        log.Printf("Error fetching events: %v", err)
        return nil
    }
	
	// decodes all the events 
    err = cursor.All(context.Background(), &events)
	if err != nil {
		log.Printf("Error decoding events: %v", err)
	 	return nil
	}

	return events
}


func FetchEventByID(event_id string) (*models.Event, error) { 

	var event models.Event

	collection := GetCollection(os.Getenv("EVENTS_COLLECTION"))

	filter := bson.M{"_id" : utils.StringToPrimitive(event_id)}

	err := collection.FindOne(context.TODO(), filter).Decode(&event)
	if err != nil {
		if err == mongo.ErrNilDocument {
			return nil, err
		}
		return nil, err
	}

	return &event, nil
}


func DeleteParticipantFromEvent(event_id string, user_id string) error{

	log.Printf("event id - %v", event_id)
	log.Printf("user id - %v", user_id)

	collection := GetCollection(os.Getenv("EVENTS_COLLECTION"))

	filter := bson.M{"_id": utils.StringToPrimitive(event_id)}
	update := bson.M{"$pull": bson.M{"participants": user_id}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil{
		return fmt.Errorf(err.Error())
	}

	if result.MatchedCount == 0{
		return fmt.Errorf("no participant found with user_id: %v in event_id: %v", user_id, event_id)
	}
	
	return nil

}

func FetchEventByFilters(query bson.M) ([]models.Event, error) { 

	var events []models.Event

	collection := GetCollection(os.Getenv("EVENTS_COLLECTION"))
	cursor, err := collection.Find(context.Background(), query)
	if err != nil {
		return events, err
	}

	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &events)
	if err != nil {
		return events, err
	}

	return events, nil
}