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

type mockNewSessionRenderer struct {
	err error
}

func (m *mockNewSessionRenderer) NewSession(w http.ResponseWriter, _ *http.Request, form *form.NewSessionForm) error {
	if m.err != nil {
		return m.err
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(form); err != nil {
		panic(err)
	}

	return nil
}

func TestNewSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		renderer       *mockNewSessionRenderer
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful",
			renderer:       &mockNewSessionRenderer{},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Email":{"Message":"","Value":""},"Password":{"Message":"","Value":""}}`,
		},
		{
			name: "error from renderer",
			renderer: &mockNewSessionRenderer{
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

			handlers.NewSession(r, testLogger(t), tt.renderer)

			status, _, body := get(r, "/session/new")

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedBody != body {
				t.Errorf("expected body %q; got %q", tt.expectedBody, body)
			}
		})
	}
}

type mockUserAuthenticator struct {
	form *form.NewSessionForm
	user models.User
	err  error
}

func (m *mockUserAuthenticator) AuthenticateUser(_ context.Context, form *form.NewSessionForm) (models.User, error) {
	m.form = form
	return m.user, m.err
}

func TestCreateSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		body             url.Values
		authenticator    *mockUserAuthenticator
		renderer         *mockNewSessionRenderer
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
			authenticator: &mockUserAuthenticator{
				user: models.User{ID: uuid.New()},
			},
			renderer:         &mockNewSessionRenderer{},
			setter:           &mockSetter{&mockCurrentUserIDSetter{}, &mockFlashSetter{}},
			expectedStatus:   http.StatusSeeOther,
			expectedLocation: "/",
		},
		{
			name: "form invalid",
			body: url.Values{
				"email": []string{"john@bannister.com"},
			},
			authenticator:  &mockUserAuthenticator{},
			renderer:       &mockNewSessionRenderer{},
			setter:         &mockSetter{&mockCurrentUserIDSetter{}, &mockFlashSetter{}},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"Email":{"Message":"","Value":"john@bannister.com"},"Password":{"Message":"Please provide a password.","Value":""}}`,
		},
		{
			name: "no user from authenticator",
			body: url.Values{
				"email":    []string{"john@bannister.com"},
				"password": []string{"password"},
			},
			authenticator: &mockUserAuthenticator{
				user: models.User{},
			},
			renderer:       &mockNewSessionRenderer{},
			setter:         &mockSetter{&mockCurrentUserIDSetter{}, &mockFlashSetter{}},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"Email":{"Message":"No user found with this email address and password.","Value":"john@bannister.com"},"Password":{"Message":"","Value":"password"}}`,
		},
		{
			name: "error from authenticator",
			body: url.Values{
				"email":    []string{"john@bannister.com"},
				"password": []string{"password"},
			},
			authenticator: &mockUserAuthenticator{
				err: errors.New("mock error"),
			},
			renderer:       &mockNewSessionRenderer{},
			setter:         &mockSetter{&mockCurrentUserIDSetter{}, &mockFlashSetter{}},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name: "error from setter",
			body: url.Values{
				"email":    []string{"john@bannister.com"},
				"password": []string{"password"},
			},
			authenticator: &mockUserAuthenticator{
				user: models.User{ID: uuid.New()},
			},
			renderer: &mockNewSessionRenderer{},
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

			handlers.CreateSession(r, testLogger(t), tt.authenticator, tt.renderer, tt.setter)

			status, headers, body := post(r, "/session", tt.body.Encode())

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
