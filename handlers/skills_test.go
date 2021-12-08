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

func TestNewSkill(t *testing.T) {
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

			request := httptest.NewRequest(http.MethodGet, "/vampires/12345678-90ab-cdef-1234-567890abcdef/skills/new", nil)
			response := httptest.NewRecorder()

			handlers.NewSkill(e, tt.vampireGetter)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}
		})
	}
}

type mockSkillCreator struct {
	receivedVampireID   uuid.UUID
	receivedDescription string
	err                 error
}

func (m *mockSkillCreator) CreateSkill(_ context.Context, vampireID uuid.UUID, description string) (models.Skill, error) {
	m.receivedVampireID = vampireID
	m.receivedDescription = description
	return models.Skill{}, m.err
}

func TestCreateSkill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		body                url.Values
		skillCreator        *mockSkillCreator
		expectedStatus      int
		expectedVampireID   uuid.UUID
		expectedDescription string
		expectedLocation    string
	}{
		{
			name:                "successful",
			body:                url.Values{"description": []string{"A description"}},
			skillCreator:        &mockSkillCreator{},
			expectedStatus:      http.StatusSeeOther,
			expectedVampireID:   uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedDescription: "A description",
			expectedLocation:    "/vampires/12345678-90ab-cdef-1234-567890abcdef",
		},
		{
			name:                "not found",
			body:                url.Values{"description": []string{"A description"}},
			skillCreator:        &mockSkillCreator{err: models.ErrNotFound},
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
			e.Renderer = templates.NewEchoRenderer(e)

			request := httptest.NewRequest(http.MethodPost, "/vampires/12345678-90ab-cdef-1234-567890abcdef/skills", strings.NewReader(tt.body.Encode()))
			request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationForm)
			response := httptest.NewRecorder()

			handlers.CreateSkill(e, tt.skillCreator)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if tt.expectedVampireID != tt.skillCreator.receivedVampireID {
				t.Errorf("expected %q; got %q", tt.expectedVampireID, tt.skillCreator.receivedVampireID)
			}

			if tt.expectedDescription != tt.skillCreator.receivedDescription {
				t.Errorf("expected %q; got %q", tt.expectedDescription, tt.skillCreator.receivedDescription)
			}

			actualLocation := response.Header().Get(echo.HeaderLocation)
			if tt.expectedLocation != actualLocation {
				t.Errorf("expected %q; got %q", tt.expectedLocation, actualLocation)
			}
		})
	}
}
