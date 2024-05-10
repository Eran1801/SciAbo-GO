package requests

import (
	"log"
	"net/http"
	"os"
	"sci-abo-go/models"
	"sci-abo-go/storage"

	"github.com/gin-gonic/gin"
)

func UploadUserProfilePicture(c *gin.Context) {

	user, _ := c.Get("user") // Extract email from the user that send with the jwt token
	user_email := user.(*models.User).Email

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
	if err := storage.UpdateDocDB(collection_name, user.(*models.User).ID, updates); err != nil {
		ErrorResponse(c, "Error updating user")
		return
	}

	SuccessResponse(c, "Profile image upload successfully", nil)
}