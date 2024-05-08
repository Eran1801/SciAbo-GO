package requests

import (
	"log"
	"net/http"
	"os"
	"sci-abo-go/models"
	"sci-abo-go/storage"
	"sci-abo-go/utils"
	"strings"
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

	if err := storage.InsertUserDB(&user); err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	SuccessResponse(c, "User created successfully", nil)
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
}

func ForgotPassword(c *gin.Context) {

	var forget_password utils.ForgetPassword
	var reset_code models.ResetCode

	if err := c.ShouldBindJSON(&forget_password); err != nil {
		ErrorResponse(c, "error in binding JSON")
		return
	}

	email := forget_password.Email
	if strings.Contains(email, "@") {

		code := utils.Create4DigitCode() // create 4 digit code
		reset_code.Code = code
		reset_code.Time = time.Now()

		// save the code in the db
		id, err := storage.InsertResetCodeDB(&reset_code)
		if err != nil {
			ErrorResponse(c, err.Error())
			return
		}

		data := map[string]interface{}{
			"reset_code_id":   id,
		}

		// send email to user with the code
		// err := utils.SendEmailWithGoMail(email, "utils/html/forget_password.html", code)
		// if err != nil {
		// 	ErrorResponse(c, "Failed to send email")
		// } else {
		SuccessResponse(c, "Mail send successfully", data)
		// }

	} else {
		ErrorResponse(c, "Email not valid")
		return
	}

}

func ValidateResetCode(c *gin.Context) {
	
	var validate_reset utils.ValidateResetCode

	if err := c.ShouldBindJSON(&validate_reset); err != nil {
		ErrorResponse(c, "error in binding JSON")
		return
	}

	// retrieve the reset code model from db using the id
	model, err := storage.GetResetCodeByID(utils.StringToPrimitive(validate_reset.ID))
	if err != nil {
		log.Println(err.Error())
		ErrorResponse(c,err.Error())
		return
	}

	// check equal codes
	if model.Code != validate_reset.UserCode {
		ErrorResponse(c,"Code not match")
		return
	}

	SuccessResponse(c,"Valid code",nil)
}

func ResetPassword(c *gin.Context) {

	var reset_password utils.ResetPassword

	if err := c.ShouldBindJSON(&reset_password); err != nil {
		ErrorResponse(c, "error in binding JSON")
		return
	}

	// todo: encrypt,verify user, lower the email before saving.

	if reset_password.Password == reset_password.ConfirmPassword {
		
		// set the updates to know which fields to update in the db
		updates := map[string]interface{}{
			"password": reset_password.Password,
		}

		// retrieve user to update
		user, _ := storage.GetUserByEmail(reset_password.Email)

		// update doc
		storage.UpdateDocDB(os.Getenv("USER_COLLECTION"), user.ID, updates)

		SuccessResponse(c, "passwords match", nil)
		return
	} else {
		ErrorResponse(c, "passwords won't match")
	}

}
