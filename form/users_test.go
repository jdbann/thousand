package form_test

import (
	"testing"

	"emailaddress.horse/thousand/form"
)

func TestNewUserForm_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		email               string
		password            string
		wantResult          bool
		wantEmailMessage    string
		wantPasswordMessage string
	}{
		{
			name:                "valid",
			email:               "john@bannister.com",
			password:            "password",
			wantResult:          true,
			wantEmailMessage:    "",
			wantPasswordMessage: "",
		},
		{
			name:                "email must be present",
			email:               "",
			password:            "password",
			wantResult:          false,
			wantEmailMessage:    "Please provide an email address.",
			wantPasswordMessage: "",
		},
		{
			name:                "email must be valid",
			email:               "john",
			password:            "password",
			wantResult:          false,
			wantEmailMessage:    "This doesn't look like a valid email address. Please take another look.",
			wantPasswordMessage: "",
		},
		{
			name:                "password must be present",
			email:               "john@bannister.com",
			password:            "",
			wantResult:          false,
			wantEmailMessage:    "",
			wantPasswordMessage: "Please provide a password.",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			form := form.NewUser(tt.email, tt.password)

			result := form.Valid()

			if tt.wantResult != result {
				t.Errorf("expected result %t; got %t", tt.wantResult, result)
			}

			if tt.wantEmailMessage != form.Email.Message {
				t.Errorf("expected email message %q; got %q", tt.wantEmailMessage, form.Email.Message)
			}

			if tt.wantPasswordMessage != form.Password.Message {
				t.Errorf("expected password message %q; got %q", tt.wantPasswordMessage, form.Password.Message)
			}
		})
	}
}
