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

type mockNewExperienceRenderer struct {
	err error
}

func (m *mockNewExperienceRenderer) NewExperience(w http.ResponseWriter, memory models.Memory) error {
	if m.err != nil {
		return m.err
	}

	_, err := w.Write([]byte(memory.ID.String()))
	if err != nil {
		panic(err)
	}

	return nil
}

type mockMemoryGetter struct {
	vampireID uuid.UUID
	memoryID  uuid.UUID
	memory    models.Memory
	err       error
}

func (m *mockMemoryGetter) GetMemory(_ context.Context, vampireID, memoryID uuid.UUID) (models.Memory, error) {
	m.vampireID = vampireID
	m.memoryID = memoryID
	return m.memory, m.err
}

func TestNewExperience(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		renderer          *mockNewExperienceRenderer
		getter            *mockMemoryGetter
		path              string
		expectedStatus    int
		expectedBody      string
		expectedVampireID uuid.UUID
		expectedMemoryID  uuid.UUID
	}{
		{
			name:     "successful",
			renderer: &mockNewExperienceRenderer{},
			getter: &mockMemoryGetter{
				memory: models.Memory{
					ID:        uuid.MustParse("22222222-2222-2222-2222-222222222222"),
					VampireID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				},
			},
			path:              "/vampires/11111111-1111-1111-1111-111111111111/memories/22222222-2222-2222-2222-222222222222/experiences/new",
			expectedStatus:    http.StatusOK,
			expectedBody:      "22222222-2222-2222-2222-222222222222",
			expectedVampireID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			expectedMemoryID:  uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		},
		{
			name:           "error parsing vampire id",
			renderer:       &mockNewExperienceRenderer{},
			getter:         &mockMemoryGetter{},
			path:           "/vampires/unknown/memories/22222222-2222-2222-2222-222222222222/experiences/new",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name:           "error parsing vampire id",
			renderer:       &mockNewExperienceRenderer{},
			getter:         &mockMemoryGetter{},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/memories/unknown/experiences/new",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name:     "not found from getter",
			renderer: &mockNewExperienceRenderer{},
			getter: &mockMemoryGetter{
				err: models.ErrNotFound,
			},
			path:              "/vampires/11111111-1111-1111-1111-111111111111/memories/22222222-2222-2222-2222-222222222222/experiences/new",
			expectedStatus:    http.StatusNotFound,
			expectedBody:      "404: Not Found",
			expectedVampireID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			expectedMemoryID:  uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		},
		{
			name:     "error from getter",
			renderer: &mockNewExperienceRenderer{},
			getter: &mockMemoryGetter{
				err: errors.New("mock error"),
			},
			path:              "/vampires/11111111-1111-1111-1111-111111111111/memories/22222222-2222-2222-2222-222222222222/experiences/new",
			expectedStatus:    http.StatusInternalServerError,
			expectedBody:      "500: Internal Server Error",
			expectedVampireID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			expectedMemoryID:  uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		},
		{
			name: "error from renderer",
			renderer: &mockNewExperienceRenderer{
				err: errors.New("mock error"),
			},
			getter: &mockMemoryGetter{
				memory: models.Memory{
					ID:        uuid.MustParse("22222222-2222-2222-2222-222222222222"),
					VampireID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				},
			},
			path:              "/vampires/11111111-1111-1111-1111-111111111111/memories/22222222-2222-2222-2222-222222222222/experiences/new",
			expectedStatus:    http.StatusInternalServerError,
			expectedBody:      "500: Internal Server Error",
			expectedVampireID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			expectedMemoryID:  uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewMux()

			handlers.NewExperience(r, testLogger(t), tt.renderer, tt.getter)

			status, _, body := get(r, tt.path)

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedBody != body {
				t.Errorf("expected body %q; got %q", tt.expectedBody, body)
			}

			if tt.expectedVampireID != tt.getter.vampireID {
				t.Errorf("expected getter to receive vampire ID %s; got %s", tt.expectedVampireID, tt.getter.vampireID)
			}

			if tt.expectedMemoryID != tt.getter.memoryID {
				t.Errorf("expected getter to receive memory ID %s; got %s", tt.expectedMemoryID, tt.getter.memoryID)
			}
		})
	}
}

type mockExperienceCreator struct {
	err         error
	vampireID   uuid.UUID
	memoryID    uuid.UUID
	description string
}

func (m *mockExperienceCreator) CreateExperience(_ context.Context, vampireID, memoryID uuid.UUID, description string) (models.Experience, error) {
	m.vampireID = vampireID
	m.memoryID = memoryID
	m.description = description
	return models.Experience{}, m.err
}

func TestCreateExperience(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		body                url.Values
		creator             *mockExperienceCreator
		path                string
		expectedStatus      int
		expectedBody        string
		expectedLocation    string
		expectedVampireID   uuid.UUID
		expectedMemoryID    uuid.UUID
		expectedDescription string
	}{
		{
			name: "successful",
			body: url.Values{
				"description": []string{"A description"},
			},
			creator:             &mockExperienceCreator{},
			path:                "/vampires/11111111-1111-1111-1111-111111111111/memories/22222222-2222-2222-2222-222222222222/experiences",
			expectedStatus:      http.StatusSeeOther,
			expectedLocation:    "/vampires/11111111-1111-1111-1111-111111111111",
			expectedVampireID:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			expectedMemoryID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			expectedDescription: "A description",
		},
		{
			name: "error parsing vampire id",
			body: url.Values{
				"description": []string{"A description"},
			},
			creator:        &mockExperienceCreator{},
			path:           "/vampires/unknown/memories/22222222-2222-2222-2222-222222222222/experiences",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name: "error parsing memory id",
			body: url.Values{
				"description": []string{"A description"},
			},
			creator:        &mockExperienceCreator{},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/memories/unknown/experiences",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name: "not found from creator",
			body: url.Values{
				"description": []string{"A description"},
			},
			creator: &mockExperienceCreator{
				err: models.ErrNotFound,
			},
			path:                "/vampires/11111111-1111-1111-1111-111111111111/memories/22222222-2222-2222-2222-222222222222/experiences",
			expectedStatus:      http.StatusNotFound,
			expectedBody:        "404: Not Found",
			expectedVampireID:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			expectedMemoryID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			expectedDescription: "A description",
		},
		{
			name: "error from creator",
			body: url.Values{
				"description": []string{"A description"},
			},
			creator: &mockExperienceCreator{
				err: errors.New("mock error"),
			},
			path:                "/vampires/11111111-1111-1111-1111-111111111111/memories/22222222-2222-2222-2222-222222222222/experiences",
			expectedStatus:      http.StatusInternalServerError,
			expectedBody:        "500: Internal Server Error",
			expectedVampireID:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			expectedMemoryID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			expectedDescription: "A description",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewMux()

			handlers.CreateExperience(r, testLogger(t), tt.creator)

			status, headers, body := post(r, tt.path, tt.body.Encode())

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedBody != body {
				t.Errorf("expected body %q; got %q", tt.expectedBody, body)
			}

			location := headers.Get("Location")
			if tt.expectedLocation != location {
				t.Errorf("expected location %q; got %q", tt.expectedLocation, location)
			}

			if tt.expectedVampireID != tt.creator.vampireID {
				t.Errorf("expected %q; got %q", tt.expectedVampireID, tt.creator.vampireID)
			}

			if tt.expectedMemoryID != tt.creator.memoryID {
				t.Errorf("expected %q; got %q", tt.expectedMemoryID, tt.creator.memoryID)
			}

			if tt.expectedDescription != tt.creator.description {
				t.Errorf("expected %q; got %q", tt.expectedDescription, tt.creator.description)
			}
		})
	}
}
