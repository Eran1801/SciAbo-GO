package requests

import (
	"log"
	"net/http"
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

	// Upload file to S3 and get the URL
	file_path := "Users/" + user_email + "/profile picture" + header.Filename
	image_url, err := storage.UploadFileToS3(file, user_email, file_path)
	if err != nil {
		ErrorResponse(c, "Failed to upload file")
		log.Println("Error uploading file to S3: ", err)
		return
	}

	// Update the user in the database with the new profile image URL
	updates := map[string]interface{}{
		"profile_image_url": image_url,
	}
	if err := storage.UpdateUser(user_email, updates); err != nil {
		ErrorResponse(c, "Error updating user")
		return
	}

	SuccessResponse(c, "Profile image upload successfully", nil)
}
