package requests

import (
	"net/http"
	"os"
	"sci-abo-go/models"
	"sci-abo-go/storage"
	"sci-abo-go/utils"
	// "sci-abo-go/utils/html"
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

	// convert email to lowercase 
	user.Email = strings.ToLower(user.Email)

	if err := utils.ValidateDbRequirements(&user); err != nil {
		ErrorResponse(c, err.Error())
		return
	}
	
	user.Password = utils.EncryptPassword(user.Password)

    user.JoinedEventIDs = make([]string, 0) // init an empty list

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

	// first we need to check if there is any user with this email in our db
	user, _ := storage.GetUserByEmail(forget_password.Email)
	if user == nil {
		ErrorResponse(c,"email not found")
		return
	}
	
	// create reset code instance
	reset := utils.CreateResetCode(&reset_code)

	// inserting reset code to the db for 5 minutes
	code, err := storage.InsertResetCodeDB(&reset)
	if err != nil { 
		ErrorResponse(c, err.Error())
		return
	}

	// todo: needs to remove the comment
	// send email to user with the code
	// err = utils.SendEmailWithGoMail(user.Email, html.GetEmailTemplate("reset_code"), code)
	// if err != nil {
	// 	ErrorResponse(c, "Failed to send email")
	// } else {

		SuccessResponse(c, "Mail send successfully", code)
	// }
}


func ValidateResetCode(c *gin.Context) {
	
	var validate_reset utils.ValidateResetCode

	if err := c.ShouldBindJSON(&validate_reset); err != nil {
		ErrorResponse(c, "error in binding JSON")
		return
	}

	// retrieve the reset code model from db using the id
	model, _ := storage.GetResetCodeByID(utils.StringToPrimitive(validate_reset.ID))
	if model == nil {
		ErrorResponse(c, "code expired")
		return
	}

	// check if the code is correct
	if model.Code != validate_reset.UserCode {
		ErrorResponse(c,"Code not match")
		return
	}

	SuccessResponse(c,"Valid code",nil)
}


func ResetPassword(c *gin.Context) {

	var reset_password utils.ResetPassword

	if err := c.ShouldBindJSON(&reset_password); err != nil {
		ErrorResponse(c, "`error in binding JSON`")
		return
	}

	if reset_password.Password == reset_password.ConfirmPassword {
		
		// set the updates to know which fields to update in the db
		encrypted_password := utils.EncryptPassword(reset_password.Password)
		updates := map[string]interface{}{
			"password": string(encrypted_password),
		}

		// retrieve user to update
		user, err := storage.GetUserByEmail(reset_password.Email)
		if err != nil {
			ErrorResponse(c, err.Error())
			return
		}

		// update user password in the db
		err = storage.UpdateDocDB(os.Getenv("USER_COLLECTION"), user.ID, updates)
		if err != nil { 
			ErrorResponse(c,err.Error())
			return
		}

		SuccessResponse(c, "passwords change successfully", nil)
		return
	} else {
		ErrorResponse(c, "passwords won't match")
		return
	}

}

func ResendResetCode(c *gin.Context) {

	var reset_code_entity utils.ResendCode

	err := c.ShouldBindJSON(&reset_code_entity)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}

	// Initialize the reset variable
	reset := &models.ResetCode{}

	code := utils.Create4DigitCode()
	reset.Code = code
	reset.Time = time.Now()

	id, err := storage.InsertResetCodeDB(reset)
	if err != nil {
		ErrorResponse(c, err.Error())
		return
	}
	// send reset code to the user email
	// utils.SendEmailWithGoMail(reset_code_entity.Email, html.GetEmailTemplate("reset_code"), code)

	SuccessResponse(c, "code save and send successfully", id)
}

func ChangePassword(c *gin.Context) { 

	user, _ := c.Get("user")
	user_model, exists := user.(*models.User)
	if !exists { 
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var request_data utils.ChangePassword

	err := c.ShouldBindJSON(&request_data)
	if err != nil { 
		ErrorResponse(c, err.Error())
		return
	}

	// check if the current password is true
	err = bcrypt.CompareHashAndPassword([]byte(user_model.Password), []byte(request_data.CurrentPassword))
	if err != nil {
		ErrorResponse(c, "wrong current password")
		return
	}

	// check match new password and confirm password
	if request_data.NewPassword != request_data.ConfirmNewPassword {
		ErrorResponse(c, "passwords not match")
		return
	}

	hash_password := utils.EncryptPassword(request_data.NewPassword)
	updates := map[string]interface{}{
		"password": hash_password}

	err = storage.UpdateDocDB(os.Getenv("USER_COLLECTION"), user_model.ID, updates)
	if err != nil { 
		ErrorResponse(c,err.Error())
		return
	}

	SuccessResponse(c, "password changed successfully", nil)

}
