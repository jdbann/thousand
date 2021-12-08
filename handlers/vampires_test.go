package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/templates"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type mockShowVampiresRenderer struct {
	err error
}

func (m *mockShowVampiresRenderer) ShowVampires(w http.ResponseWriter, vampires []models.Vampire) error {
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

func (m *mockNewVampireRenderer) NewVampire(w http.ResponseWriter) error {
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
	name    string
	vampire models.Vampire
	err     error
}

func (m *mockVampireCreator) CreateVampire(_ context.Context, name string) (models.Vampire, error) {
	m.name = name
	return m.vampire, m.err
}

func TestCreateVampire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		body             url.Values
		creator          *mockVampireCreator
		expectedStatus   int
		expectedName     string
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
			expectedStatus:   http.StatusSeeOther,
			expectedName:     "Gruffudd",
			expectedLocation: "/vampires/12345678-90ab-cdef-1234-567890abcdef",
		},
		{
			name: "error from creator",
			body: url.Values{"name": []string{"Gruffudd"}},
			creator: &mockVampireCreator{
				err: errors.New("mock error"),
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedName:     "Gruffudd",
			expectedLocation: "",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewMux()

			handlers.CreateVampire(r, testLogger(t), tt.creator)

			status, headers, _ := post(r, "/vampires", tt.body.Encode())

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedName != tt.creator.name {
				t.Errorf("expected creator to receive name %q; got %q", tt.expectedName, tt.creator.name)
			}

			actualLocation := headers.Get("Location")
			if tt.expectedLocation != actualLocation {
				t.Errorf("expected location %q; got %q", tt.expectedLocation, actualLocation)
			}
		})
	}
}

type mockVampireGetter struct {
	vampire    models.Vampire
	err        error
	receivedID uuid.UUID
}

func (m *mockVampireGetter) GetVampire(_ context.Context, id uuid.UUID) (models.Vampire, error) {
	m.receivedID = id
	return m.vampire, m.err
}

func TestShowVampire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		vampireGetter  *mockVampireGetter
		expectedStatus int
	}{
		{
			name:           "successful",
			vampireGetter:  &mockVampireGetter{vampire: models.Vampire{}},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not found",
			vampireGetter:  &mockVampireGetter{err: models.ErrNotFound},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Renderer = templates.NewEchoRenderer(e)

			request := httptest.NewRequest(http.MethodGet, "/vampires/12345678-90ab-cdef-1234-567890abcdef", nil)
			response := httptest.NewRecorder()

			handlers.ShowVampire(e, tt.vampireGetter)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}
		})
	}
}
