package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"emailaddress.horse/thousand/app/models"
	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/templates"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type mockVampiresGetter struct{}

func (m mockVampiresGetter) GetVampires(_ context.Context) ([]models.Vampire, error) {
	return []models.Vampire{}, nil
}

func TestListVampires(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "successful",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Renderer = templates.NewRenderer(e)

			request := httptest.NewRequest(http.MethodGet, "/vampires", nil)
			response := httptest.NewRecorder()

			handlers.ListVampires(e, mockVampiresGetter{})

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}
		})
	}
}

func TestNewVampire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "successful",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Renderer = templates.NewRenderer(e)

			request := httptest.NewRequest(http.MethodGet, "/vampires/new", nil)
			response := httptest.NewRecorder()

			handlers.NewVampire(e)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}
		})
	}
}

type mockVampireCreator struct {
	vampire      models.Vampire
	receivedName string
}

func (m *mockVampireCreator) CreateVampire(_ context.Context, name string) (models.Vampire, error) {
	m.receivedName = name
	return m.vampire, nil
}

func TestCreateVampire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		body             url.Values
		vampireCreator   *mockVampireCreator
		expectedStatus   int
		expectedName     string
		expectedLocation string
	}{
		{
			name: "successful",
			body: url.Values{"name": []string{"Gruffudd"}},
			vampireCreator: &mockVampireCreator{
				vampire: models.Vampire{
					ID: uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
				},
			},
			expectedStatus:   http.StatusSeeOther,
			expectedName:     "Gruffudd",
			expectedLocation: "/vampires/12345678-90ab-cdef-1234-567890abcdef",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Renderer = templates.NewRenderer(e)

			request := httptest.NewRequest(http.MethodPost, "/vampires", strings.NewReader(tt.body.Encode()))
			request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationForm)
			response := httptest.NewRecorder()

			handlers.CreateVampire(e, tt.vampireCreator)
			handlers.ShowVampire(e, nil)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if tt.expectedName != tt.vampireCreator.receivedName {
				t.Errorf("expected %q; got %q", tt.expectedName, tt.vampireCreator.receivedName)
			}

			actualLocation := response.Header().Get(echo.HeaderLocation)
			if tt.expectedLocation != actualLocation {
				t.Errorf("expected %q; got %q", tt.expectedLocation, actualLocation)
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
			e.Renderer = templates.NewRenderer(e)

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
