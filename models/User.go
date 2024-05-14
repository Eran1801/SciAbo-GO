package models

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID                              primitive.ObjectID      `json:"_id" bson:"_id,omitempty"`
    FirstName                       string                  `json:"first_name" bson:"first_name" validate:"required"` 
    LastName                        string                  `json:"last_name" bson:"last_name" validate:"required"`
    Email                           string                  `json:"email" bson:"email" validate:"required,email"`
    Password                        string                  `json:"password" bson:"password" validate:"required,passwordPattern"`
    ConfirmPassword                 string                  `json:"confirm_password" bson:"-" validate:"required,eqfield=Password"`
    ProfileImageURL                 string                  `json:"profile_image" bson:"profile_image_url"`
    LinkedinProfile                 string                  `json:"linkedin_profile" bson:"linkedin_profile" validate:"url"`
    Country                         string                  `json:"user_country" bson:"user_country" validate:"required"`
    AcademicInstitutionOrCompany    string                  `json:"academic_institution_or_company" bson:"academic_institution_or_company" validate:"required"`
    Role                            string                  `json:"role" bson:"role" validate:"required"`
    PrincipalInvestigator           string                  `json:"principal_investigator" bson:"principal_investigator"`
    Industry                        string                  `json:"industry" bson:"industry" validate:"required"`
    About                           string                  `json:"about" bson:"about" validate:"required"`
    JoinedEventIDs                  []string                `json:"joined_event_ids" bson:"joined_event_ids"`
}

// define a custom validator instance
var validate *validator.Validate

// init() function initialize automatic when our program run
func init(){
    validate = validator.New()

    // add custom validate for password pattern
    validate.RegisterValidation("passwordPattern", ValidatePassword)
}

func ValidateUser(user *User) error {
    err := validate.Struct(user)
    if err != nil {
        return err
    }
    return nil
}

func ValidatePassword(f1 validator.FieldLevel) bool {
    password := f1.Field().String()
    
    // Compile separate regex for each requirement
    var (
        hasLower     = regexp.MustCompile(`[a-z]`).MatchString(password)
        hasUpper     = regexp.MustCompile(`[A-Z]`).MatchString(password)
        hasDigit     = regexp.MustCompile(`[0-9]`).MatchString(password)
        hasSymbol    = regexp.MustCompile(`[!@#$%^&*()]`).MatchString(password)
        hasMinLength = len(password) >= 8
    )

    // Check all conditions are met
    return hasLower && hasUpper && hasDigit && hasSymbol && hasMinLength
}
