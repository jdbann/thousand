package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"emailaddress.horse/thousand/form"
	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type mockNewUserRenderer struct {
	err error
}

func (m *mockNewUserRenderer) NewUser(w http.ResponseWriter, _ *http.Request, form *form.NewUserForm) error {
	if m.err != nil {
		return m.err
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(form); err != nil {
		panic(err)
	}

	return nil
}

func TestNewUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		renderer       *mockNewUserRenderer
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful",
			renderer:       &mockNewUserRenderer{},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Email":{"Message":"","Value":""},"Password":{"Message":"","Value":""}}`,
		},
		{
			name: "error from renderer",
			renderer: &mockNewUserRenderer{
				err: errors.New("mock error"),
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewMux()

			handlers.NewUser(r, testLogger(t), tt.renderer)

			status, _, body := get(r, "/user/new")

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedBody != body {
				t.Errorf("expected body %q; got %q", tt.expectedBody, body)
			}
		})
	}
}

type mockUserCreator struct {
	form *form.NewUserForm
	user models.User
	err  error
}

func (m *mockUserCreator) CreateUser(_ context.Context, form *form.NewUserForm) (models.User, error) {
	m.form = form
	return m.user, m.err
}

type mockCurrentUserIDSetter struct {
	id  uuid.UUID
	err error
}

func (m *mockCurrentUserIDSetter) SetCurrentUserID(_ *http.Request, _ http.ResponseWriter, id uuid.UUID) error {
	m.id = id
	return m.err
}

type mockFlashSetter struct {
	message string
	err     error
}

func (m *mockFlashSetter) SetFlash(_ *http.Request, _ http.ResponseWriter, message string) error {
	m.message = message
	return m.err
}

type mockSetter struct {
	*mockCurrentUserIDSetter
	*mockFlashSetter
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		body             url.Values
		creator          *mockUserCreator
		renderer         *mockNewUserRenderer
		setter           *mockSetter
		expectedStatus   int
		expectedBody     string
		expectedLocation string
	}{
		{
			name: "successful",
			body: url.Values{
				"email":    []string{"john@bannister.com"},
				"password": []string{"password"},
			},
			creator:          &mockUserCreator{},
			renderer:         &mockNewUserRenderer{},
			setter:           &mockSetter{&mockCurrentUserIDSetter{}, &mockFlashSetter{}},
			expectedStatus:   http.StatusSeeOther,
			expectedLocation: "/",
		},
		{
			name: "form invalid",
			body: url.Values{
				"email": []string{"john@bannister.com"},
			},
			creator:        &mockUserCreator{},
			renderer:       &mockNewUserRenderer{},
			setter:         &mockSetter{&mockCurrentUserIDSetter{}, &mockFlashSetter{}},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"Email":{"Message":"","Value":"john@bannister.com"},"Password":{"Message":"Please provide a password.","Value":""}}`,
		},
		{
			name: "email already in use from creator",
			body: url.Values{
				"email":    []string{"john@bannister.com"},
				"password": []string{"password"},
			},
			creator: &mockUserCreator{
				err: models.ErrEmailAlreadyInUse,
			},
			renderer:       &mockNewUserRenderer{},
			setter:         &mockSetter{&mockCurrentUserIDSetter{}, &mockFlashSetter{}},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"Email":{"Message":"Email already in use.","Value":"john@bannister.com"},"Password":{"Message":"","Value":"password"}}`,
		},
		{
			name: "error from creator",
			body: url.Values{
				"email":    []string{"john@bannister.com"},
				"password": []string{"password"},
			},
			creator: &mockUserCreator{
				err: errors.New("mock error"),
			},
			renderer:       &mockNewUserRenderer{},
			setter:         &mockSetter{&mockCurrentUserIDSetter{}, &mockFlashSetter{}},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name: "error from creator",
			body: url.Values{
				"email":    []string{"john@bannister.com"},
				"password": []string{"password"},
			},
			creator:  &mockUserCreator{},
			renderer: &mockNewUserRenderer{},
			setter: &mockSetter{
				&mockCurrentUserIDSetter{
					err: errors.New("mock error"),
				},
				&mockFlashSetter{},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewMux()

			handlers.CreateUser(r, testLogger(t), tt.creator, tt.renderer, tt.setter)

			status, headers, body := post(r, "/user", tt.body.Encode())

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			actualLocation := headers.Get("Location")
			if tt.expectedLocation != actualLocation {
				t.Errorf("expected location %q; got %q", tt.expectedLocation, actualLocation)
			}

			if tt.expectedBody != body {
				t.Errorf("expected body %q; got %q", tt.expectedBody, body)
			}
		})
	}
}
