package utils

type LoginRequest struct {

	Email 		string 		`json:"email"`
	Password 	string 		`json:"password"`

}

type ForgetPassword struct {
	Email string `json:"email"`
}

type ResetPassword struct {

	Email 				string 		`json:"email" bson:"email"`
	Password 			string 		`json:"password"`
	ConfirmPassword 	string 		`json:"confirm_password" validate:"eqfield=Password"`
}

type ValidateResetCode struct {

	ID 		 	string 	 `json:"id" bson:"id"`
	UserCode 	string 	 `json:"user_code" bson:"user_code"`
}

