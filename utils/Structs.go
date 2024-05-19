package utils

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgetPassword struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ValidateResetCode struct {
	ID       string `json:"id"`
	UserCode string `json:"user_code"`
}

type ResendCode struct {
	Email string `json:"email"`
}

type Participants struct {
	Participants []string `json:"participants"`
}

type ChangePassword struct {
	CurrentPassword    string `json:"current_password"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

type UpdateUserDetailsRequest struct {
	Email                        string `json:"email" validate:"required,email"`
	LinkedinProfile              string `json:"linkedin_profile" validate:"url"`
	Country                      string `json:"user_country" validate:"required"`
	AcademicInstitutionOrCompany string `json:"academic_institution_or_company" validate:"required"`
	Role                         string `json:"role" validate:"required"`
	PrincipalInvestigator        string `json:"principal_investigator"`
	Industry                     string `json:"industry" validate:"required"`
	About                        string `json:"about" validate:"required"`
}

type SearchFilters struct {
	Year            string      `form:"year"`
	ConferenceName  string      `form:"conference_name"`
	Abbreviation    string      `form:"abbreviation"`
	Country         string      `form:"country"`
	City            string     `form:"city"`
}

type FirstMessageRequest struct {
	ReceiverID     string `json:"receiver_id"`
	ReceiverEmail  string `json:"receiver_email"`
	MessageContent string `json:"message_content"`
}

type MessageRequest struct {

    RoomID              string      `json:"room_id"`
	MessageContent      string      `json:"message_content"`
}
