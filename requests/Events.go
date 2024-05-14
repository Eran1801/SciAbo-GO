package requests

import (
	"net/http"
	"os"
	"sci-abo-go/models"
	"sci-abo-go/storage"
	"sci-abo-go/utils"

	"github.com/gin-gonic/gin"
)

func AddEvent(c *gin.Context) {

	var event models.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// initialize empty list of participants before saving
	event.Participants = make([]string, 0)

	_, err = storage.InsertEventDB(&event)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "event added successful", nil)

}

func GetAllUserEvents(c *gin.Context) {
	/*
		It returns all the events that the user is signed-in
		This function returns two dict, one is past events the second is future events.
	*/

	user, _ := c.Get("user")
	user_model, exists := user.(*models.User)
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// get all the events id's that the user is sign-in to
	var user_events_ids []string = user_model.JoinedEventIDs

	// convert event_ids(string) to ObjectID
	event_ids := utils.FromStringListToPrimitiveList(user_events_ids)

	// fetch all the events(models.Event) that the user is join to
	events := storage.FetchUserEvents(event_ids)

	// order by past and future events
	events_by_past_and_future := utils.DivideEventsToPastFuture(events)

	SuccessResponse(c, "success", events_by_past_and_future)

}

func GetAllParticipatesInEvent(c *gin.Context) {
	/*
		A function that get's as a prams a list of string that holds all id's of the users
		that is signed to this event and return all the Users to the client.
	*/

	var participants_ids utils.Participants
	err := c.ShouldBindJSON(&participants_ids)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	participants_ids_primitive := utils.FromStringListToPrimitiveList(participants_ids.Participants)

	users, _ := storage.GetUsersByIDs(participants_ids_primitive)

	SuccessResponse(c, "success", users)

}

func GetEventByID(c *gin.Context) {

	event_id := c.Query("event_id")
	if event_id == "" {
		ErrorResponse(c, "event id is not given")
		return
	}

	event, err := storage.FetchEventByID(event_id)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "success", event)

}

func JoinEvent(c *gin.Context) {
	/*
	In this endpoint a user can join to an event.
	This function do 3 main things:
		* add user id to the event participants
		* add event id to the user events
		* add the details of the user on the event (when he will arrive and etc)
	*/ 

	// get the user
	user, _ := c.Get("user")
	user_model, exists := user.(*models.User)
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// initialize the entity of the user details about join to an event
	var event_participant models.EventParticipant

	// valid data from body
	err := c.ShouldBindJSON(&event_participant)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// set the ID of the user in the EventParticipant entity
	event_participant.UserID = user_model.ID.Hex()

	// get the event id from query params
	event_id := c.Query("event_id")
	if event_id == "" {
		ErrorResponse(c, "event_id is empty")
		return
	}

	// update event in the db
	obj_event_id := utils.StringToPrimitive(event_id)
	event_collection_name := os.Getenv("EVENTS_COLLECTION")
	err = storage.AddParticipantToEvent(event_collection_name, obj_event_id, user_model.ID.Hex())
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}	
	
	// add event id to the events ids array of the user (the ones he join in to)
	user_collection_name := os.Getenv("USER_COLLECTION")
	err = storage.AddEventIdToUserEvents(user_collection_name, user_model.ID, event_id)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}		

	// set the user id to the join event details to assign the participant(user) to his id 
	event_participant.UserID = user_model.ID.Hex()

	err = storage.InsertEventParticipant(event_participant)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c,"success",nil)
}