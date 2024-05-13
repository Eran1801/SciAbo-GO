package utils

type LoginRequest struct {

	Email 		string 		`json:"email"`
	Password 	string 		`json:"password"`

}

type ForgetPassword struct {
	Email string `json:"email" validate:"email"`
}

type ResetPassword struct {

	Email 				string 		`json:"email" bson:"email"`
	Password 			string 		`json:"password" bson:"password"`
	ConfirmPassword 	string 		`json:"confirm_password" bson:"confirm_password"`
}

type ValidateResetCode struct {

	ID 		 	string 	 `json:"id" bson:"id"`
	UserCode 	string 	 `json:"user_code" bson:"user_code"`
}

type ResendCode struct {

	Email string `json:"id" bson:"id"`

}

type Participants struct {

	Participants  []string `json:"participants" bson:"participants"`
}

type ChangePassword struct { 

	CurrentPassword 		string 	`json:"current_password" bson:"current_password"`
	NewPassword 			string	`json:"new_password" bson:"new_password"`
	ConfirmNewPassword 		string	`json:"confirm_new_password" bson:"confirm_new_password"`

}