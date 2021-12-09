package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	vampireID   uuid.UUID
	description string
	err         error
}

func (m *mockSkillCreator) CreateSkill(_ context.Context, vampireID uuid.UUID, description string) (models.Skill, error) {
	m.vampireID = vampireID
	m.description = description
	return models.Skill{}, m.err
}

func TestCreateSkill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		body                url.Values
		creator             *mockSkillCreator
		path                string
		expectedStatus      int
		expectedBody        string
		expectedLocation    string
		expectedVampireID   uuid.UUID
		expectedDescription string
	}{
		{
			name:                "successful",
			body:                url.Values{"description": []string{"A description"}},
			creator:             &mockSkillCreator{},
			path:                "/vampires/12345678-90ab-cdef-1234-567890abcdef/skills",
			expectedStatus:      http.StatusSeeOther,
			expectedVampireID:   uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedDescription: "A description",
			expectedLocation:    "/vampires/12345678-90ab-cdef-1234-567890abcdef",
		},
		{
			name:           "error parsing vampire ID",
			body:           url.Values{"description": []string{"A description"}},
			creator:        &mockSkillCreator{},
			path:           "/vampires/unknown/skills",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name: "not found from creator",
			body: url.Values{"description": []string{"A description"}},
			creator: &mockSkillCreator{
				err: models.ErrNotFound,
			},
			path:                "/vampires/12345678-90ab-cdef-1234-567890abcdef/skills",
			expectedStatus:      http.StatusNotFound,
			expectedBody:        "404: Not Found",
			expectedVampireID:   uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedDescription: "A description",
		},
		{
			name: "error from creator",
			body: url.Values{"description": []string{"A description"}},
			creator: &mockSkillCreator{
				err: errors.New("mock error"),
			},
			path:                "/vampires/12345678-90ab-cdef-1234-567890abcdef/skills",
			expectedStatus:      http.StatusInternalServerError,
			expectedBody:        "500: Internal Server Error",
			expectedVampireID:   uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedDescription: "A description",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewMux()

			handlers.CreateSkill(r, testLogger(t), tt.creator)

			status, headers, body := post(r, tt.path, tt.body.Encode())

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedBody != body {
				t.Errorf("expected body %q; got %q", tt.expectedBody, body)
			}

			location := headers.Get("Location")
			if tt.expectedLocation != location {
				t.Errorf("expected %q; got %q", tt.expectedLocation, location)
			}

			if tt.expectedVampireID != tt.creator.vampireID {
				t.Errorf("expected %q; got %q", tt.expectedVampireID, tt.creator.vampireID)
			}

			if tt.expectedDescription != tt.creator.description {
				t.Errorf("expected %q; got %q", tt.expectedDescription, tt.creator.description)
			}
		})
	}
}
