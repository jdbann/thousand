package form

import "go.uber.org/zap/zapcore"

type NewUserForm struct {
	Email    stringField
	Password stringField
}

var (
	newUserEmailValidations = stringValidations{
		stringPresent("Please provide an email address."),
		stringEmailFormat("This doesn't look like a valid email address. Please take another look."),
	}

	newUserPasswordValidations = stringValidations{
		stringPresent("Please provide a password."),
	}
)

func NewUser(email, password string) *NewUserForm {
	return &NewUserForm{
		Email:    stringField{Value: email},
		Password: stringField{Value: password},
	}
}

func (f *NewUserForm) Valid() bool {
	success := true

	if !newUserEmailValidations.validate(&f.Email) {
		success = false
	}

	if !newUserPasswordValidations.validate(&f.Password) {
		success = false
	}

	return success
}

func (f NewUserForm) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("email", f.Email.Value)
	if f.Password.Value != "" {
		enc.AddString("password", "xxxx")
	} else {
		enc.AddString("password", "")
	}
	return nil
}
