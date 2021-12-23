package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/middleware"
	"emailaddress.horse/thousand/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type mockShowVampiresRenderer struct {
	err error
}

func (m *mockShowVampiresRenderer) ShowVampires(w http.ResponseWriter, _ *http.Request, vampires []models.Vampire) error {
	if m.err != nil {
		return m.err
	}

	names := make([]string, len(vampires))
	for i, v := range vampires {
		names[i] = v.Name
	}

	_, err := w.Write([]byte(strings.Join(names, ", ")))
	if err != nil {
		panic(err)
	}

	return nil
}

type mockVampiresGetter struct {
	vampires []models.Vampire
	err      error
}

func (m *mockVampiresGetter) GetVampires(_ context.Context) ([]models.Vampire, error) {
	return m.vampires, m.err
}

func TestListVampires(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		renderer       *mockShowVampiresRenderer
		getter         *mockVampiresGetter
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "successful",
			renderer: &mockShowVampiresRenderer{},
			getter: &mockVampiresGetter{
				vampires: []models.Vampire{
					{Name: "one"},
					{Name: "two"},
					{Name: "three"},
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "one, two, three",
		},
		{
			name: "error from getter",
			getter: &mockVampiresGetter{
				err: errors.New("mock error"),
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name: "error from renderer",
			renderer: &mockShowVampiresRenderer{
				err: errors.New("mock error"),
			},
			getter: &mockVampiresGetter{
				vampires: []models.Vampire{
					{Name: "one"},
					{Name: "two"},
					{Name: "three"},
				},
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

			handlers.ListVampires(r, testLogger(t), tt.renderer, tt.getter)

			status, _, body := get(r, "/vampires")

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedBody != body {
				t.Errorf("expected body %q; got %q", tt.expectedBody, body)
			}
		})
	}
}

type mockNewVampireRenderer struct {
	err error
}

func (m *mockNewVampireRenderer) NewVampire(w http.ResponseWriter, _ *http.Request) error {
	if m.err != nil {
		return m.err
	}

	_, err := w.Write([]byte("new vampire"))
	if err != nil {
		panic(err)
	}

	return nil
}

func TestNewVampire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		renderer       *mockNewVampireRenderer
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful",
			renderer:       &mockNewVampireRenderer{},
			expectedStatus: http.StatusOK,
			expectedBody:   "new vampire",
		},
		{
			name: "error from renderer",
			renderer: &mockNewVampireRenderer{
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

			handlers.NewVampire(r, testLogger(t), tt.renderer)

			status, _, body := get(r, "/vampires/new")

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedBody != body {
				t.Errorf("expected body %q; got %q", tt.expectedBody, body)
			}
		})
	}
}

type mockVampireCreator struct {
	userID  uuid.UUID
	name    string
	vampire models.Vampire
	err     error
}

func (m *mockVampireCreator) CreateVampire(_ context.Context, id uuid.UUID, name string) (models.Vampire, error) {
	m.userID = id
	m.name = name
	return m.vampire, m.err
}

func TestCreateVampire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		body             url.Values
		creator          *mockVampireCreator
		user             models.User
		expectedStatus   int
		expectedName     string
		expectedUserID   uuid.UUID
		expectedLocation string
	}{
		{
			name: "successful",
			body: url.Values{"name": []string{"Gruffudd"}},
			creator: &mockVampireCreator{
				vampire: models.Vampire{
					ID: uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
				},
			},
			user: models.User{
				ID: uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			},
			expectedStatus:   http.StatusSeeOther,
			expectedName:     "Gruffudd",
			expectedUserID:   uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedLocation: "/vampires/12345678-90ab-cdef-1234-567890abcdef",
		},
		{
			name: "error from creator",
			body: url.Values{"name": []string{"Gruffudd"}},
			creator: &mockVampireCreator{
				err: errors.New("mock error"),
			},
			user: models.User{
				ID: uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedName:     "Gruffudd",
			expectedUserID:   uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedLocation: "",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewMux()

			handlers.CreateVampire(r, testLogger(t), tt.creator)

			req := postRequest("/vampires", tt.body.Encode())
			req.request = middleware.RequestWithCurrentUser(req.request, tt.user)

			status, headers, _ := req.perform(r)

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedName != tt.creator.name {
				t.Errorf("expected creator to receive name %q; got %q", tt.expectedName, tt.creator.name)
			}

			if tt.expectedUserID != tt.creator.userID {
				t.Errorf("expected creator to receive name %q; got %q", tt.expectedUserID, tt.creator.userID)
			}

			actualLocation := headers.Get("Location")
			if tt.expectedLocation != actualLocation {
				t.Errorf("expected location %q; got %q", tt.expectedLocation, actualLocation)
			}
		})
	}
}

type mockShowVampireRenderer struct {
	err error
}

func (m *mockShowVampireRenderer) ShowVampire(w http.ResponseWriter, _ *http.Request, vampire models.Vampire) error {
	if m.err != nil {
		return m.err
	}

	_, err := w.Write([]byte(vampire.Name))
	if err != nil {
		panic(err)
	}

	return nil
}

type mockVampireGetter struct {
	vampire models.Vampire
	err     error
	id      uuid.UUID
}

func (m *mockVampireGetter) GetVampire(_ context.Context, id uuid.UUID) (models.Vampire, error) {
	m.id = id
	return m.vampire, m.err
}

func TestShowVampire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		renderer       *mockShowVampireRenderer
		getter         *mockVampireGetter
		path           string
		expectedStatus int
		expectedBody   string
		expectedID     uuid.UUID
	}{
		{
			name:     "successful",
			renderer: &mockShowVampireRenderer{},
			getter: &mockVampireGetter{
				vampire: models.Vampire{Name: "A vampire"},
			},
			path:           "/vampires/12345678-90ab-cdef-1234-567890abcdef",
			expectedStatus: http.StatusOK,
			expectedBody:   "A vampire",
			expectedID:     uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
		},
		{
			name:           "error parsing id",
			renderer:       &mockShowVampireRenderer{},
			getter:         &mockVampireGetter{},
			path:           "/vampires/unknown",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name: "not found error from getter",
			getter: &mockVampireGetter{
				err: models.ErrNotFound,
			},
			path:           "/vampires/12345678-90ab-cdef-1234-567890abcdef",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404: Not Found",
			expectedID:     uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
		},
		{
			name: "error from getter",
			getter: &mockVampireGetter{
				err: errors.New("mock error"),
			},
			path:           "/vampires/12345678-90ab-cdef-1234-567890abcdef",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
			expectedID:     uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
		},
		{
			name: "error from renderer",
			renderer: &mockShowVampireRenderer{
				err: errors.New("mock error"),
			},
			path: "/vampires/12345678-90ab-cdef-1234-567890abcdef",
			getter: &mockVampireGetter{
				vampire: models.Vampire{Name: "A vampire"},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
			expectedID:     uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewMux()

			handlers.ShowVampire(r, testLogger(t), tt.renderer, tt.getter)

			status, _, body := get(r, tt.path)

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedBody != body {
				t.Errorf("expected body %q; got %q", tt.expectedBody, body)
			}

			if tt.expectedID != tt.getter.id {
				t.Errorf("expected getter to receive ID %s; got %s", tt.expectedID, tt.getter.id)
			}
		})
	}
}
