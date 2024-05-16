package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"sci-abo-go/models"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"

	"github.com/gin-gonic/gin"
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

func ValidateStruct(model interface{}) error {
	err := models.ValidateModel(model)
	if err != nil {
		// handle validation errors
		errors := err.(validator.ValidationErrors)
		// construct error message
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
	// convert from string to primitive.ObjectID
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
		start_date := event.StartDate

		if start_date.Before(now) {
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

func GetUserFromCookie(c *gin.Context) *models.User {

	user, _ := c.Get("user")
	user_model, exists := user.(*models.User)
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return nil
	}

	return user_model

}

func StructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(data)

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldValue := v.Field(i).Interface()
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			result[jsonTag] = fieldValue
		} else {
			result[field.Name] = fieldValue
		}
	}

	return result
}

func CheckFilters(filters SearchFilters) primitive.M{ 

	query := bson.M{}

	if filters.ConferenceName != "" {
		query["name"] = bson.M{"$regex":filters.ConferenceName, "$options":"i"}
	}

	if filters.Abbreviation != "" {
        query["abbreviation"] = bson.M{"$regex": filters.Abbreviation, "$options": "i"}
    }

    if filters.Country != "" {
        query["country"] = filters.Country
    }

    if filters.City != "" {
        query["city"] = filters.City
    }

    if filters.Year != "" {
        year, err := strconv.Atoi(filters.Year)
        if err == nil {
            start_date := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
            end_date := time.Date(year, time.December, 31, 23, 59, 59, 999, time.UTC)
            query["start_date"] = bson.M{"$gte": start_date, "$lt": end_date}
        }
    }

    return query

}
