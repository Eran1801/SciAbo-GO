package requests

import (
	"github.com/gin-gonic/gin"

)

func GetAllEvents(c *gin.Context) {

	user, _ := c.Get("user")

	SuccessResponse(c,"inside get_all_events function",user)
}