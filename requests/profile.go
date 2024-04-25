package requests

import (
	"encoding/json"
	"log"
	"net/http"
	"sci-abo-go/storage"
)

func UploadProfilePictureHandler(w http.ResponseWriter, r *http.Request) {

    email := "jbghvbjhbjohedfdg@example.com"
    user,err := storage.GetUserByEmail(email)
    if err != nil{
        http.Error(w, "Error fetching user", http.StatusBadRequest)
    }
    log.Println(user)

    // Parse the multipart form data with a max size of 10 MB
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        http.Error(w, "File too large or incorrect data.", http.StatusBadRequest)
        log.Println("Error parsing form data: ", err)
        return
    }

    // Extract file from the parsed form
    file, header, err := r.FormFile("profile_image")
    if err != nil {
        http.Error(w, "Invalid file upload.", http.StatusBadRequest)
        log.Println("Error retrieving file from form: ", err)
        return
    }
    defer file.Close()

    // Upload file to S3 and get the URL
    image_url, err := storage.UploadFileToS3(file, header.Filename, user.ID.String())
    if err != nil {
        http.Error(w, "Failed to upload file: "+err.Error(), http.StatusInternalServerError)
        log.Println("Error uploading file to S3: ", err)
        return
    }

    if user == nil{
        http.Error(w, "No user found with this email:", http.StatusBadRequest)
    }else{
        updates := map[string]interface{}{
            "profile_image_url": image_url,
        }
        
        err := storage.UpdateUser(email, updates)
        if err != nil {
            log.Println("Error updating user:", err)
        }
    }

    // Log the successful upload and respond to the client
    log.Println("Profile image uploaded successfully: ", image_url)
    w.Header().Set("Content-Type", "application/json")
    response := map[string]string{"message": "Profile image uploaded successfully"}
    json.NewEncoder(w).Encode(response)

}
