package userapi

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	userpb "github.com/openfms/protos/gen/user/v1"
	"regexp"
)

// ValidateSignupUser validates signup user request
//func (us *UserService) ValidateSignupUser(req *userpb.SignupRequest) error {
//	return validation.ValidateStruct(req,
//		validation.Field(&req.Email, validation.Required, is.Email, validation.Length(8, 255)),
//		validation.Field(&req.Password, validation.Required, validation.Length(1, 255)),
//		validation.Field(&req.Name, validation.Required, validation.Length(1, 255)),
//	)
//}

// ValidateSignInUser validates signIn user request
func (us *UserService) ValidateSignInUser(req *userpb.SignInRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.UserNameEmail, validation.Required, validation.Length(1, 255), validation.By(isValidEmailOrUsername)),
		validation.Field(&req.Password, validation.Required, validation.Length(8, 255)),
	)
}

// isValidEmailOrUsername checks whether the given value is a valid email or a valid username.
func isValidEmailOrUsername(value any) error {
	str, ok := value.(string)
	if !ok {
		return validation.NewError("validation_invalid_type", "invalid type")
	}

	// Define the regular expression patterns for email and username.
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	usernamePattern := `^[a-zA-Z0-9_]{4,}$`

	// Check if it's a valid email
	if match, _ := regexp.MatchString(emailPattern, str); match {
		return nil
	}

	// Check if it's a valid username
	if match, _ := regexp.MatchString(usernamePattern, str); match {
		return nil
	}

	return validation.NewError("validation_invalid_userid", "invalid email or username")
}
