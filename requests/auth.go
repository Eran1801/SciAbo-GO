package requests

import (
	"sci-abo-go/models"
	"sci-abo-go/storage"
	"sci-abo-go/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	// "github.com/golang-jwt/jwt/v5"
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

	if err := storage.SaveUserInDatabase(&user); err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "User created successfully", user.Email)
}

func Login(c *gin.Context) {

    var login_request utils.LoginRequest 

    // get the email
    if c.Bind(&login_request) != nil {
        ErrorResponse(c, "Failed to read body")
        return
    }

    // get the user from the db
    var user *models.User
    user ,err := storage.GetUserByEmail(login_request.Email)
    if err != nil {
        ErrorResponse(c, err.Error())
        return
    }
    

    // compare between login request pass and user password in the db
    err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(login_request.Password))
    if err != nil {
        ErrorResponse(c,"invalid email or password")
        return
    }

    



}
