package requests

import (

	"sci-abo-go/models"
	"sci-abo-go/storage"
	"sci-abo-go/utils"

	"github.com/gin-gonic/gin"

)

func CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        ErrorResponse(c, err.Error())
        return
    }

    if err := utils.ValidateDbRequirements(&user); err != nil {
        ErrorResponse(c, err.Error())
        return
    }

    if err := utils.EncryptPassword(&user); err != nil {
        ErrorResponse(c, err.Error())
        return
    }

    if err := storage.SaveUserInDB(&user); err != nil {
        ErrorResponse(c, err.Error())
        return
    }

    SuccessResponse(c, "User created successfully", user.Email)
}
