package requests

import (
	"encoding/json"
	"log"
	"net/http"
	"sci-abo-go/storage"
)

func UploadUserProfilePicture(w http.ResponseWriter, r *http.Request) {

	id := "66321220f4d099e0b3d466c"
	user, err := storage.GetUserById(id)

	if err != nil {
		ErrorResponse("Error fetching user from db", w)
		return

	} else if user == nil {
		ErrorResponse("No user found with this id", w)
		return
	}

	// Parse the multipart form data with a max size of 10 MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		ErrorResponse("File too large or incorrect data", w)
		return
	}

	// Extract file from the parsed form
	file, header, err := r.FormFile("profile_image")
	if err != nil {
		ErrorResponse("No user found with this id", w)
		log.Println("Error retrieving file from form: ", err)
		return
	}
	defer file.Close()

	// Upload file to S3 and get the URL
	image_url, err := storage.UploadFileToS3(file, header.Filename, id)
	if err != nil {
		http.Error(w, "Failed to upload file: " + err.Error(), http.StatusInternalServerError)
		log.Println("Error uploading file to S3: ", err)
		return
	}

	// update the user db and add the profile url image
	updates := map[string]interface{}{
		"profile_image_url": image_url,
	}
	err = storage.UpdateUser(id, updates)
	if err != nil {
		http.Error(w, "Error updating user:: " + err.Error(), http.StatusInternalServerError)
		return
	}

	// log the successful upload and respond to the client
	log.Println("Profile image uploaded successfully: ", image_url)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Profile image uploaded successfully"}
	json.NewEncoder(w).Encode(response)

}
