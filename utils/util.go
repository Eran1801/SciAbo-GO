package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"sci-abo-go/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(user *models.User) error {

	hash_password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("field to hash password")
	}
	user.Password = string(hash_password)
	return nil
}

func ValidateDbRequirements(user *models.User) error {
	err := models.ValidateUser(user)
	if err != nil {
		// Handle validation errors
		errors := err.(validator.ValidationErrors)
		// Construct error message
		var errMsg string
		for _, e := range errors {
			errMsg += e.Field() + " is " + e.Tag() + "; "
			if len(errMsg) > 0 {
				return fmt.Errorf(errMsg)
			}
		}

	}
	return nil
}

func Create4DigitCode() string {

	result := ""

	for i:= 0; i < 4; i++ {
		random_number := fmt.Sprint(rand.Intn(9) + 1)
		result += random_number
	}

	return result
}

func SendEmailWithGoMail(to string, templatePath string, code string) error {
    // Prepare the email template
    var body bytes.Buffer
    t, _ := template.ParseFiles(templatePath)

    // Execute the template with the provided data
    err := t.Execute(&body, struct{ Code string }{Code: code})
    if err != nil {
        log.Fatalf("Error executing template: %v", err)
        return err
    }

    // Create a new message
    m := gomail.NewMessage()

    // Set the sender and recipient(s)
    m.SetHeader("From", "SciAboConferenceHub@outlook.com")
    m.SetHeader("To", to)

    // Set the subject
    m.SetHeader("Subject", "Forget Password Code")

    // Set the email body as HTML
    m.SetBody("text/html", body.String())

    // Configure the SMTP dialer
    d := gomail.NewDialer("smtp.office365.com", 587, "SciAboConferenceHub@outlook.com", "Eran1302")

    // Now send the email
    if err := d.DialAndSend(m); err != nil {
        log.Fatalf("Error sending email: %v", err)
    }

	return nil
}

func StringToPrimitive(hex string) primitive.ObjectID {
	// convert from primitive.ObjectID to string

    oid, _ := primitive.ObjectIDFromHex(hex)
    return oid
}