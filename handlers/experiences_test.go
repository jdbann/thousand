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

type mockMemoryGetter struct {
	memory models.Memory
	err    error
}

func (m mockMemoryGetter) GetMemory(_ context.Context, _, _ uuid.UUID) (models.Memory, error) {
	return m.memory, m.err
}

func TestNewExperience(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		memoryGetter   mockMemoryGetter
		expectedStatus int
	}{
		{
			name:           "successful",
			memoryGetter:   mockMemoryGetter{memory: models.Memory{}},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not found",
			memoryGetter:   mockMemoryGetter{err: models.ErrNotFound},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Renderer = templates.NewRenderer(e)

			request := httptest.NewRequest(http.MethodGet, "/vampires/12345678-90ab-cdef-1234-567890abcdef/memories/12345678-90ab-cdef-1234-567890abcdef/experiences/new", nil)
			response := httptest.NewRecorder()

			handlers.NewExperience(e, tt.memoryGetter)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}
		})
	}
}

type mockExperienceCreator struct {
	err                 error
	receivedVampireID   uuid.UUID
	receivedMemoryID    uuid.UUID
	receivedDescription string
}

func (m *mockExperienceCreator) CreateExperience(_ context.Context, vampireID, memoryID uuid.UUID, description string) (models.Experience, error) {
	m.receivedVampireID = vampireID
	m.receivedMemoryID = memoryID
	m.receivedDescription = description
	return models.Experience{}, m.err
}

func TestCreateExperience(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		body                url.Values
		experienceCreator   *mockExperienceCreator
		expectedStatus      int
		expectedVampireID   uuid.UUID
		expectedMemoryID    uuid.UUID
		expectedDescription string
	}{
		{
			name: "successful",
			body: url.Values{
				"description": []string{"A description"},
			},
			experienceCreator:   &mockExperienceCreator{},
			expectedStatus:      http.StatusSeeOther,
			expectedVampireID:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			expectedMemoryID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			expectedDescription: "A description",
		},
		{
			name: "not found",
			body: url.Values{
				"description": []string{"A description"},
			},
			experienceCreator:   &mockExperienceCreator{err: models.ErrNotFound},
			expectedStatus:      http.StatusNotFound,
			expectedVampireID:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			expectedMemoryID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			expectedDescription: "A description",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Renderer = templates.NewRenderer(e)

			request := httptest.NewRequest(http.MethodPost, "/vampires/11111111-1111-1111-1111-111111111111/memories/22222222-2222-2222-2222-222222222222/experiences", strings.NewReader(tt.body.Encode()))
			request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationForm)
			response := httptest.NewRecorder()

			handlers.ShowVampire(e, nil)
			handlers.CreateExperience(e, tt.experienceCreator)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if tt.expectedVampireID != tt.experienceCreator.receivedVampireID {
				t.Errorf("expected %q; got %q", tt.expectedVampireID, tt.experienceCreator.receivedVampireID)
			}

			if tt.expectedMemoryID != tt.experienceCreator.receivedMemoryID {
				t.Errorf("expected %q; got %q", tt.expectedMemoryID, tt.experienceCreator.receivedMemoryID)
			}

			if tt.expectedDescription != tt.experienceCreator.receivedDescription {
				t.Errorf("expected %q; got %q", tt.expectedDescription, tt.experienceCreator.receivedDescription)
			}
		})
	}
}
