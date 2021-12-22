package form

type NewSessionForm struct {
	Email    stringField
	Password stringField
}

var (
	newSessionEmailValidations = stringValidations{
		stringPresent("Please provide an email address."),
	}

	newSessionPasswordValidations = stringValidations{
		stringPresent("Please provide a password."),
	}
)

func NewSession(email, password string) *NewSessionForm {
	return &NewSessionForm{
		Email:    stringField{Value: email},
		Password: stringField{Value: password},
	}
}

func (f *NewSessionForm) Valid() bool {
	success := true

	if !newSessionEmailValidations.validate(&f.Email) {
		success = false
	}

	if !newSessionPasswordValidations.validate(&f.Password) {
		success = false
	}

	return success
}
