package requests

import (
	"net/http"
	"os"
	"sci-abo-go/models"
	"sci-abo-go/storage"
	"sci-abo-go/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func AddEvent(c *gin.Context) {

	var event models.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	event.CreateTime = time.Now() 	// time creation of event
	event.Verified = "0" // new event needs to be approved. in search event it's extract only events with Verified value of "1"

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

	user_model := utils.GetUserFromCookie(c)

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
	user_model := utils.GetUserFromCookie(c)

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
	err = storage.AddParticipantToEvent(obj_event_id, user_model.ID.Hex())
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

	err = storage.InsertEventParticipantDB(event_participant)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "success", nil)
}


func SearchEvent(c *gin.Context){

	var filters utils.SearchFilters

	// parse query params inti the filters struct
	err := c.ShouldBindQuery(&filters)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	query := utils.CheckFilters(filters) 

	query["verified"] = "1" // return only approved events by admin
	filter_events, err := storage.FetchEventByFilters(query)
	if err != nil { 
		ErrorResponse(c, err.Error())
		return
	}
	
	SuccessResponse(c, "success", filter_events)

}

func UploadEventPic(c *gin.Context) {

	event, err := storage.FetchEventByID(c.Query("event_id"))
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// Limit the size of the form to 10 MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)

	// Parse the multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		ErrorResponse(c, "File too large or incorrect data")
		return
	}

	// Retrieve the file from the form
	file, header, err := c.Request.FormFile("event_pic")
	if err != nil {
		ErrorResponse(c, "Error retrieving file from form")
		return
	}
	defer file.Close()

	// set the path to save in the bucket
	file_path := "Events/" + event.ID.Hex() + "/event pic" + header.Filename

	// Upload file to S3 and get the URL
	image_url, err := storage.UploadFileToS3(file, file_path)
	if err != nil {
		ErrorResponse(c, "Failed to upload file")
		return
	}

	// set what to update
	updates := map[string]interface{}{
		"event_image_url": image_url,
	}

	// Update the user in the database with the new profile image URL
	collection_name := os.Getenv("USER_COLLECTION")
	if err := storage.UpdateDocDB(collection_name, event.ID, updates); err != nil {
		ErrorResponse(c, "Error updating user")
		return
	}

	SuccessResponse(c, "Profile image upload successfully", nil)


}

