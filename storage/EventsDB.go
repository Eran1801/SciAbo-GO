package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"sci-abo-go/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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


func AddParticipantToEvent(collection_name string, event_id primitive.ObjectID, participant_id string) error {

	// get event collection
	collection := GetCollection(collection_name)

	log.Printf("event_id %s", event_id)
	log.Printf("participant_id %s", participant_id)

	filter := bson.M{"_id": event_id}
	update := bson.M{"$push": bson.M{"participants": participant_id}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
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


func FetchUserEvents(events_ids []primitive.ObjectID) map[string][]models.Event {
	/*
		Given a list of id's of events, this function returns all the events
		that the user is signed-in to, separate to past and future events. 
	*/

	event_collection := GetCollection("events")
    var events []models.Event
    filter := bson.M{"_id": bson.M{"$in": events_ids}}
    cursor, err := event_collection.Find(context.Background(), filter)
    if err != nil {
        log.Printf("Error fetching events: %v", err)
        return nil
    }
    err = cursor.All(context.Background(), &events)
	if err != nil {
		log.Printf("Error decoding events: %v", err)
	 	return nil
	}

	process_events := DivideEventsToPastFuture(events)

	return process_events
}


func DivideEventsToPastFuture(events []models.Event) map[string][]models.Event {
	// Prepare to separate the events into past and future
    now := time.Now()
    past_events := make([]models.Event, 0)
    future_events := make([]models.Event, 0)

    for _, event := range events {
        startDate, err := time.Parse("2006-01-02", event.StartDate) // Assuming date format is YYYY-MM-DD
        if err != nil {
            log.Printf("Error parsing start date for event ID %s: %v", event.ID.Hex(), err)
            continue
        }

        if startDate.Before(now) {
            past_events = append(past_events, event)
        } else {
            future_events = append(future_events, event)
        }
    }

    // Categorize events into a map of past and future
    categorizedEvents := map[string][]models.Event{
        "past_events":   past_events,
        "future_events": future_events,
    }

    return categorizedEvents
}
