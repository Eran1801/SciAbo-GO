package requests

import (
	"log"
	"os"
	"sci-abo-go/models"
	"sci-abo-go/storage"
	"sci-abo-go/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddEvent(c *gin.Context) {

	user, _ := c.Get("user")
	user_model := user.(*models.User)

	var event models.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// initialize empty list of participants before saving
	event.Participants = make([]string, 0)

	event_id, err := storage.InsertEventDB(&event)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// convert event_id(string) to primitive.ObjectID
	eventObjectID, _ := primitive.ObjectIDFromHex(event_id)

	// add user id to participants in event
	storage.AddParticipantToEvent(os.Getenv("EVENTS_COLLECTION"), eventObjectID, user_model.ID.Hex())

	// add event id to the joined_event_ids of the user
	storage.AddEventIdToUserEvents(os.Getenv("USER_COLLECTION"), user_model.ID, event_id)

	SuccessResponse(c, "event added successful", nil)

}

func GetAllUserEvents(c *gin.Context) {

	user, exists := c.Get("user")
	if !exists {
		ErrorResponse(c, "User not authenticated")
		return
	}
	user_model := user.(*models.User)

	// get all the events id's that the user is sign-in to
	var user_events_ids []string = user_model.JoinedEventIDs

	// convert event_ids(string) to ObjectID
	event_ids := utils.StringToObjectID(user_events_ids)

	// fetch all the events(models.Event) that the user is sign-in to order by past and future events
	events := storage.FetchUserEvents(event_ids)

	SuccessResponse(c, "success", events)

}

func GetAllParticipatesInEvent(c *gin.Context) {
	/*
	A function that get's as a prams a list of string that holds all id's of the users 
	that is signed to this event and return all the Users to the client.
	*/

	participants_ids, _ := c.Get("participants")
	log.Printf("participants_ids % v", participants_ids)

	SuccessResponse(c,"success",participants_ids)

}
