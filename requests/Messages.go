package requests

import (
	"os"
	"sci-abo-go/models"
	"sci-abo-go/storage"
	"sci-abo-go/utils"
	"time"

	// "sci-abo-go/utils/html"

	"github.com/gin-gonic/gin"
)

func SendFirstMessage(c *gin.Context) {

	user_model := utils.GetUserFromCookie(c)
	var message_request utils.FirstMessageRequest

	// valid request body
	err := c.ShouldBindJSON(&message_request)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// get the receiver User
	receiver_user, err := storage.GetUserByEmail(message_request.ReceiverEmail)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// send email to the user that receive the message on the new message
	// receiver_email := message_request.ReceiverEmail
	// subject := "Science Aboard - New Message"
	// url_link := "http://localhost:8080/api/messages/my_messages"

	// utils.SendEmailWithGoMail(receiver_email, subject, html.GetEmailTemplate("new_message"), url_link)

	_, room := utils.PopulateMessageAndRoomStruct(message_request, user_model, receiver_user)

	// insert room to the db and get it's room_id
	room_id, err := storage.InsertRoom(&room)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// for sender
	err = storage.AddRoomIdToUser(user_model, room_id)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// for receiver
	err = storage.AddRoomIdToUser(receiver_user, room_id)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "success", nil)
}

func GetAllRoomsByUserID(c *gin.Context) {
	// The user will activate this endpoint when he enter the "My Messages" section

	user_model := utils.GetUserFromCookie(c)

	rooms_ids := user_model.RoomsIDs

	user_rooms, err := storage.FetchUserRooms(utils.StringToPrimitiveList(rooms_ids))
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "success", user_rooms)

}

func GetMessagesByRoomID(c *gin.Context) {

	room_id := c.Query("room_id")
	room, err := storage.FetchRoomByID(utils.StringToPrimitive(room_id))
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "success", room.Messages)
}

func SendAndInsertNewMessageToRoom(c *gin.Context) {
	
	user_model := utils.GetUserFromCookie(c) // sender
	var message_request utils.MessageRequest // message content and room ID

	err := c.ShouldBindJSON(&message_request)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// Get the room that needs to be update
	room, err := storage.FetchRoomByID(utils.StringToPrimitive(message_request.RoomID))
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// create a message struct
	message := models.Message {

		Content: message_request.MessageContent,
		CreatedAt: time.Now(),
		SenderId: user_model.ID.Hex(),
		SenderFullName: user_model.FirstName + " " + user_model.LastName,
	}

	// add to the array of messages the new message
	room.Messages = append(room.Messages, message)
	updates := map[string]interface{}{
		"messages" : room.Messages,
	}

	// update the room doc in the db with the new array
	err = storage.UpdateDocDB(os.Getenv("ROOM_COLLECTION"), room.ID, updates)
	if err != nil{
		ErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "success", nil)

}
