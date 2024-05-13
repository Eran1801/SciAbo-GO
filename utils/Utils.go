package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"sci-abo-go/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)


func EncryptPassword(password string) string {

	hash_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "field to hash password"
	}
	return string(hash_password)
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

	for i := 0; i < 4; i++ {
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

func CreateResetCode(reset *models.ResetCode) models.ResetCode {
	// create the entity of the ResetCode struct 
	code := Create4DigitCode() // create 4 digit code
	reset.Code = code
	reset.Time = time.Now()

	return *reset

}

func FromStringListToPrimitiveList(ids []string) []primitive.ObjectID {

	events_ids := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		obj_id, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		events_ids = append(events_ids, obj_id)
	}

	return events_ids
}

func DivideEventsToPastFuture(events []models.Event) map[string][]models.Event {
	// when the user enter the section of 'my events' the events will shown divided to past and future events

    now := time.Now()
	
	// init 2 arrays of events to store past and future events
    past_events := make([]models.Event, 0)
    future_events := make([]models.Event, 0)

	// iterate all the events and categorize them according to their date
    for _, event := range events {
        startDate, err := time.Parse("2006-01-02", event.StartDate)
        if err != nil {
            log.Printf("Error parsing start date for event ID %s: %v", event.ID.Hex(), err)
            continue
        }

        if startDate.Before(now) {
            past_events = append(past_events, event)
        } else {
            future_events = append(future_events, event)
        }
    }

    // categorize events into a map of past and future
    categorizedEvents := map[string][]models.Event{
        "past_events":   past_events,
        "future_events": future_events,
    }

    return categorizedEvents
}