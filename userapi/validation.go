package userapi

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	userpb "github.com/irisco88/protos/gen/user/v1"
	"regexp"
)

// ValidateCreateUser validates create user request
func (us *UserService) ValidateCreateUser(req *userpb.CreateUserRequest) error {
	if e := validation.Validate(req.User, validation.Required); e != nil {
		return e
	}
	user := req.User
	return validation.ValidateStruct(user,
		validation.Field(&user.Id, validation.Empty),
		validation.Field(&user.Email, validation.When(len(user.Email) > 0, validation.Length(8, 255), is.Email)),
		validation.Field(&user.Password, validation.Required, validation.Length(8, 255)),
		validation.Field(&user.UserName, validation.Required, validation.Length(3, 255)),
		validation.Field(&user.FirstName, validation.Required, validation.Length(3, 255)),
		validation.Field(&user.LastName, validation.When(len(user.LastName) > 0, validation.Length(3, 255))),
		validation.Field(&user.Avatar, validation.Length(0, 255)),
	)
}

// ValidateUpdateUser validates update user request
func (us *UserService) ValidateUpdateUser(req *userpb.UpdateUserRequest) error {
	if e := validation.Validate(req.User, validation.Required); e != nil {
		return e
	}
	user := req.User
	return validation.ValidateStruct(user,
		validation.Field(&user.Id, validation.Required),
		validation.Field(&user.Email, validation.When(len(user.Email) > 0, validation.Length(8, 255), is.Email)),
		validation.Field(&user.Password, validation.Required, validation.Length(8, 255)),
		validation.Field(&user.UserName, validation.Required, validation.Length(3, 255)),
		validation.Field(&user.FirstName, validation.Required, validation.Length(3, 255)),
		validation.Field(&user.LastName, validation.When(len(user.LastName) > 0, validation.Length(3, 255))),
		validation.Field(&user.Avatar, validation.Length(0, 255)),
	)
}

// ValidateDeleteUser validates delete user request
func (us *UserService) ValidateDeleteUser(req *userpb.DeleteUserRequest) error {
	return validation.Validate(req.UserId, validation.Required)

}

// ValidateSignInUser validates signIn user request
func (us *UserService) ValidateSignInUser(req *userpb.SignInRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.UserNameEmail, validation.Required, validation.Length(1, 255), validation.By(isValidEmailOrUsername)),
		validation.Field(&req.Password, validation.Required, validation.Length(8, 255)),
	)
}

// ValidateSignUpUser validates signup user request
func (us *UserService) ValidateSignUpUser(req *userpb.SignUpRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Email, validation.When(len(req.Email) > 0, validation.Length(8, 255), is.Email)),
		validation.Field(&req.Password, validation.Required, validation.Length(8, 255)),
		validation.Field(&req.UserName, validation.Required, validation.Length(3, 255)),
		validation.Field(&req.FirstName, validation.Required, validation.Length(3, 255)),
		validation.Field(&req.LastName, validation.When(len(req.LastName) > 0, validation.Length(3, 255))),
		validation.Field(&req.Avatar, validation.Length(0, 255)),
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
