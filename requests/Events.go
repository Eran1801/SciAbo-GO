package requests

import (
	"sci-abo-go/models"
	"sci-abo-go/storage"

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

	id, err = storage.InsertEventDB(&event)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	event_ids := user.(*models.User).JoinedEventIDs
	append(event_ids, id)


}
