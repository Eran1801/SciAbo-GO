package requests

import (
    "net/http"   // package http provides HTTP client and server implementations.
    "sci-abo-go/models"  // importing the 'models' package where the User struct is defined.
	"encoding/json"        
)

// RegisterHandler is an HTTP handler function that registers a new user.
func registerHandler(w http.ResponseWriter, r *http.Request) {

    // declare a variable 'user' of type models.User. This will hold the parsed data.
    var user models.User

    // json.NewDecoder creates a new decoder that reads from r.Body (request body).
    // The Decode method decodes the JSON-encoded data into the 'user' variable.
    err := json.NewDecoder(r.Body).Decode(&user)

	// if NewDecoder is fails it's returns an error that it's not nil, this is why we check if err != nill
	if err != nil {
        // Send an HTTP error message; this is similar to returning an HttpResponseBadRequest in Django.
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}

	// Create a response message
    response := map[string]string{"message": "Welcome " + user.FirstName}

    // Set the content type of the response
    w.Header().Set("Content-Type", "application/json")

    // Marshal the response struct to JSON
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Write the JSON response
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)

}
