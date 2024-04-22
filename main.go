package main

import(
	"log"
	"net/http"
	"sci-abo-go/requests"
)

func main(){

	// init routes
	router := requests.InitializerRoutes()

	// start the HTTP server using the router
	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080",router)
	if err != nil {
		log.Fatalf("Server failed to start %v ", err)
	}

}

