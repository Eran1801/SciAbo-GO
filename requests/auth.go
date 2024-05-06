package requests

import (
	"net/http"
	"os"
	"sci-abo-go/models"
	"sci-abo-go/storage"
	"sci-abo-go/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
		ErrorResponse(c, "failed to read body")
		return
	}

	// get the user from the db
	var user *models.User
	user, err := storage.GetUserByEmail(login_request.Email)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// compare between login request password and user password in the db
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login_request.Password))
	if err != nil {
		ErrorResponse(c, "invalid email or password")
		return
	}

	// generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	token_string, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		ErrorResponse(c, "failed to cerate a token")
		return
	}

	// send the token to the cookie
	c.SetSameSite(http.SameSiteLaxMode)

	// todo: replace the localhost by url of deploy + change false to true
	c.SetCookie("Authorization", token_string, 3600*24*30, "/", "localhost", false, true)

	// SuccessResponse(c,"token created",token_string)

}

func Validate(c *gin.Context) {
	SuccessResponse(c, "I'm logged in", nil)
}
