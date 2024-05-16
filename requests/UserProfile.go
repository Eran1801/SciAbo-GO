package requests

import (
	"log"
	"net/http"
	"os"
	"sci-abo-go/storage"
	"sci-abo-go/utils"

	"github.com/gin-gonic/gin"
)

func UploadUserProfilePicture(c *gin.Context) {

	user_model := utils.GetUserFromCookie(c)

	user_email := user_model.Email

	// Limit the size of the form to 10 MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)

	// Parse the multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		ErrorResponse(c, "File too large or incorrect data")
		return
	}

	// Retrieve the file from the form
	file, header, err := c.Request.FormFile("profile_image")
	if err != nil {
		ErrorResponse(c, "Error retrieving file from form")
		log.Println("Error retrieving file from form: ", err)
		return
	}
	defer file.Close()

	// set the path to save in the bucket
	file_path := "Users/" + user_email + "/profile picture" + header.Filename

	// Upload file to S3 and get the URL
	image_url, err := storage.UploadFileToS3(file, user_email, file_path)
	if err != nil {
		ErrorResponse(c, "Failed to upload file")
		log.Println("Error uploading file to S3: ", err)
		return
	}

	// set what to update
	updates := map[string]interface{}{
		"profile_image_url": image_url,
	}

	// Update the user in the database with the new profile image URL
	collection_name := os.Getenv("USER_COLLECTION")
	if err := storage.UpdateDocDB(collection_name, user_model.ID, updates); err != nil {
		ErrorResponse(c, "Error updating user")
		return
	}

	SuccessResponse(c, "Profile image upload successfully", nil)
}

func UpdateUserDetails(c *gin.Context) {

	user_model := utils.GetUserFromCookie(c)

	var update_profile utils.UpdateUserDetailsRequest

	err := c.ShouldBindJSON(&update_profile)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	if err := utils.ValidateStruct(&update_profile); err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// convert update_profile to map[string]interface{}
	updated_data := utils.StructToMap(update_profile)

	err = storage.UpdateDocDB(os.Getenv("USER_COLLECTION"), user_model.ID, updated_data)
	if err != nil { 
		ErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "success", nil)

}

func DeleteUser(c *gin.Context) {
	/*
	When deleting a user from our website, we need to delete it's account with all it's details
	But first we need to delete the user details from each event he is join.
	*/

	user_model := utils.GetUserFromCookie(c)

	events_ids := user_model.JoinedEventIDs
	log.Printf("events ids - %v", events_ids)

	for _, event_id := range events_ids {
		log.Printf("event id % v", event_id)

		// for each event, delete user id from it's participant array
		err := storage.DeleteParticipantFromEvent(event_id, user_model.ID.Hex())
		if err != nil {
			ErrorResponse(c, err.Error())
			return
		}
	}

	// delete all the records of event_participant that the user is in
	err := storage.DeleteEventParticipantByUserID(user_model.ID.Hex())
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// after deleting the user from the events he was sign-in to, we need to delete the user itself
	err = storage.DeleteUser(user_model.ID)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}
	
	// remove the cookie
	c.SetCookie("Authorization", "", -1, "/", "localhost", false, true)

	SuccessResponse(c, "user deleted successfully", nil)

}
