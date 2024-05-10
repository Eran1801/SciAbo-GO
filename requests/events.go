package requests 

import (
	"sci-abo-go/models"

	"github.com/gin-gonic/gin"
)

func AddEvent(c *gin.Context) {

	user, _ := c.Get("user")
	var event models.Event

	err := c.ShouldBindJSON(&event)
	if err != nil {
		ErrorResponse(c, err.Error())
		return		
	}

	id, err = storage.InsertEventDB(event)

	event_ids := user.(*models.User).CreatedEventIDs
	append(event_ids, id)


}