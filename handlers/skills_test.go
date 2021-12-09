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
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type mockNewSkillRenderer struct {
	err error
}

func (m *mockNewSkillRenderer) NewSkill(w http.ResponseWriter, vampire models.Vampire) error {
	if m.err != nil {
		return m.err
	}

	_, err := w.Write([]byte(vampire.Name))
	if err != nil {
		panic(err)
	}

	return nil
}

func TestNewSkill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		renderer       *mockNewSkillRenderer
		getter         *mockVampireGetter
		path           string
		expectedStatus int
		expectedBody   string
		expectedID     uuid.UUID
	}{
		{
			name:     "successful",
			renderer: &mockNewSkillRenderer{},
			getter: &mockVampireGetter{
				vampire: models.Vampire{
					Name: "a vampire",
				},
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/skills/new",
			expectedStatus: http.StatusOK,
			expectedBody:   "a vampire",
			expectedID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		{
			name:           "error parsing vampire id",
			getter:         &mockVampireGetter{},
			path:           "/vampires/unknown/skills/new",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name:     "not found from getter",
			renderer: &mockNewSkillRenderer{},
			getter: &mockVampireGetter{
				err: models.ErrNotFound,
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/skills/new",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404: Not Found",
			expectedID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		{
			name:     "error from getter",
			renderer: &mockNewSkillRenderer{},
			getter: &mockVampireGetter{
				err: models.ErrNotFound,
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/skills/new",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404: Not Found",
			expectedID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		{
			name: "error from renderer",
			renderer: &mockNewSkillRenderer{
				err: models.ErrNotFound,
			},
			getter: &mockVampireGetter{
				vampire: models.Vampire{
					Name: "a vampire",
				},
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/skills/new",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
			expectedID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewMux()

			handlers.NewSkill(r, testLogger(t), tt.renderer, tt.getter)

			status, _, body := get(r, tt.path)

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedBody != body {
				t.Errorf("expected body %q; got %q", tt.expectedBody, body)
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
