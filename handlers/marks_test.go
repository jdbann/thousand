package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/templates"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TestNewMark(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		vampireGetter  *mockVampireGetter
		expectedStatus int
		expectedID     uuid.UUID
	}{
		{
			name:           "successful",
			vampireGetter:  &mockVampireGetter{vampire: models.Vampire{}},
			expectedStatus: http.StatusOK,
			expectedID:     uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
		},
		{
			name:           "not found",
			vampireGetter:  &mockVampireGetter{err: models.ErrNotFound},
			expectedStatus: http.StatusNotFound,
			expectedID:     uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Renderer = templates.NewRenderer(e)

			request := httptest.NewRequest(http.MethodGet, "/vampires/12345678-90ab-cdef-1234-567890abcdef/marks/new", nil)
			response := httptest.NewRecorder()

			handlers.NewMark(e, tt.vampireGetter)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if tt.expectedID != tt.vampireGetter.receivedID {
				t.Errorf("expected %q; got %q", tt.expectedID, tt.vampireGetter.receivedID)
			}
		})
	}
}

type mockMarkCreator struct {
	receivedVampireID   uuid.UUID
	receivedDescription string
	err                 error
}

func (m *mockMarkCreator) CreateMark(_ context.Context, vampireID uuid.UUID, description string) (models.Mark, error) {
	m.receivedVampireID = vampireID
	m.receivedDescription = description
	return models.Mark{}, m.err
}

func TestCreateMark(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		body                url.Values
		markCreator         *mockMarkCreator
		expectedStatus      int
		expectedVampireID   uuid.UUID
		expectedDescription string
		expectedLocation    string
	}{
		{
			name:                "successful",
			body:                url.Values{"description": []string{"A description"}},
			markCreator:         &mockMarkCreator{},
			expectedStatus:      http.StatusSeeOther,
			expectedVampireID:   uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedDescription: "A description",
			expectedLocation:    "/vampires/12345678-90ab-cdef-1234-567890abcdef",
		},
		{
			name:                "not found",
			body:                url.Values{"description": []string{"A description"}},
			markCreator:         &mockMarkCreator{err: models.ErrNotFound},
			expectedStatus:      http.StatusNotFound,
			expectedVampireID:   uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedDescription: "A description",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Renderer = templates.NewRenderer(e)

			request := httptest.NewRequest(http.MethodPost, "/vampires/12345678-90ab-cdef-1234-567890abcdef/marks", strings.NewReader(tt.body.Encode()))
			request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationForm)
			response := httptest.NewRecorder()

			handlers.CreateMark(e, tt.markCreator)
			handlers.ShowVampire(e, nil)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if tt.expectedVampireID != tt.markCreator.receivedVampireID {
				t.Errorf("expected %q; got %q", tt.expectedVampireID, tt.markCreator.receivedVampireID)
			}

			if tt.expectedDescription != tt.markCreator.receivedDescription {
				t.Errorf("expected %q; got %q", tt.expectedDescription, tt.markCreator.receivedDescription)
			}

			actualLocation := response.Header().Get(echo.HeaderLocation)
			if tt.expectedLocation != actualLocation {
				t.Errorf("expected %q; got %q", tt.expectedLocation, actualLocation)
			}
		})
	}
}
