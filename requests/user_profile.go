package requests

import (
	"log"
	"net/http"
	"sci-abo-go/storage"
)

func UploadUserProfilePicture(w http.ResponseWriter, r *http.Request) {
	
	email := r.FormValue("email")// needs to extract from the request
	user, err := storage.GetUserByEmail(email)

	if err != nil {
		ErrorResponse("Error fetching user from db", w)
		return

	} else if user == nil {
		ErrorResponse("No user found with this email", w)
		return
	}

	// parse the multipart form data with a max size of 10 MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		ErrorResponse("File too large or incorrect data", w)
		return
	}

	// extract file from the parsed form
	file, header, err := r.FormFile("profile_image")
	if err != nil {
		ErrorResponse("No user found with this id", w)
		log.Println("Error retrieving file from form: ", err)
		return
	}
	defer file.Close()
	

	// Upload file to S3 and get the URL
	image_url, err := storage.UploadFileToS3(file, header.Filename, email)
	if err != nil {
		ErrorResponse("Failed to upload file: ",w)
		log.Println("Error uploading file to S3: ", err)
		return
	}

	// update the user db and add the profile url image
	updates := map[string]interface{}{
		"profile_image_url": image_url,
	}
	err = storage.UpdateUser(email, updates)
	if err != nil {
		ErrorResponse("Error updating user ",w)
		return
	}

	SuccessResponse("Profile image uploaded/updated successfully", nil, w)
}
