package repository_test

import (
	"context"
	"errors"
	"testing"

	"emailaddress.horse/thousand/form"
	"emailaddress.horse/thousand/models"
)

func TestCreateUser(t *testing.T) {
	m := newTestRepository(t)

	user, err := m.CreateUser(
		context.Background(),
		form.NewUser("john@bannister.com", "password"),
	)
	if err != nil {
		t.Fatal("error creating user:", err)
	}

	_, err = m.CreateUser(
		context.Background(),
		form.NewUser(user.Email, "password"),
	)
	if err != nil {
		if !errors.Is(err, models.ErrEmailAlreadyInUse) {
			t.Fatal("unexpected error creating user with duplicate email:", err)
		}
	} else {
		t.Fatal("expected email already in use error")
	}
}
