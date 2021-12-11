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

type mockNewMarkRenderer struct {
	err error
}

func (m *mockNewMarkRenderer) NewMark(w http.ResponseWriter, vampire models.Vampire) error {
	if m.err != nil {
		return m.err
	}

	_, err := w.Write([]byte(vampire.Name))
	if err != nil {
		panic(err)
	}

	return nil
}

func TestNewMark(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		renderer       *mockNewMarkRenderer
		getter         *mockVampireGetter
		path           string
		expectedStatus int
		expectedBody   string
		expectedID     uuid.UUID
	}{
		{
			name:     "successful",
			renderer: &mockNewMarkRenderer{},
			getter: &mockVampireGetter{
				vampire: models.Vampire{
					Name: "a vampire",
				},
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/marks/new",
			expectedStatus: http.StatusOK,
			expectedBody:   "a vampire",
			expectedID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		{
			name:           "error parsing vampire id",
			getter:         &mockVampireGetter{},
			path:           "/vampires/unknown/marks/new",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name:     "not found from getter",
			renderer: &mockNewMarkRenderer{},
			getter: &mockVampireGetter{
				err: models.ErrNotFound,
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/marks/new",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404: Not Found",
			expectedID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		{
			name:     "error from getter",
			renderer: &mockNewMarkRenderer{},
			getter: &mockVampireGetter{
				err: models.ErrNotFound,
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/marks/new",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404: Not Found",
			expectedID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		{
			name: "error from renderer",
			renderer: &mockNewMarkRenderer{
				err: models.ErrNotFound,
			},
			getter: &mockVampireGetter{
				vampire: models.Vampire{
					Name: "a vampire",
				},
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/marks/new",
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

			handlers.NewMark(r, testLogger(t), tt.renderer, tt.getter)

			status, _, body := get(r, tt.path)

			if tt.expectedStatus != status {
				t.Errorf("expected %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedBody != body {
				t.Errorf("expected body %q; got %q", tt.expectedBody, body)
			}

			if tt.expectedID != tt.getter.id {
				t.Errorf("expected %q; got %q", tt.expectedID, tt.getter.id)
			}
		})
	}
}

type mockMarkCreator struct {
	vampireID   uuid.UUID
	description string
	err         error
}

func (m *mockMarkCreator) CreateMark(_ context.Context, vampireID uuid.UUID, description string) (models.Mark, error) {
	m.vampireID = vampireID
	m.description = description
	return models.Mark{}, m.err
}

func TestCreateMark(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		body                url.Values
		creator             *mockMarkCreator
		path                string
		expectedStatus      int
		expectedBody        string
		expectedLocation    string
		expectedVampireID   uuid.UUID
		expectedDescription string
	}{
		{
			name: "successful",
			body: url.Values{
				"description": []string{"a description"},
			},
			creator:             &mockMarkCreator{},
			path:                "/vampires/12345678-90ab-cdef-1234-567890abcdef/marks",
			expectedStatus:      http.StatusSeeOther,
			expectedLocation:    "/vampires/12345678-90ab-cdef-1234-567890abcdef",
			expectedVampireID:   uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedDescription: "a description",
		},
		{
			name: "error parsing vampire ID",
			body: url.Values{
				"description": []string{"a description"},
			},
			creator:        &mockMarkCreator{},
			path:           "/vampires/unknown/marks",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name: "not found from creator",
			body: url.Values{
				"description": []string{"a description"},
			},
			creator: &mockMarkCreator{
				err: models.ErrNotFound,
			},
			path:                "/vampires/12345678-90ab-cdef-1234-567890abcdef/marks",
			expectedStatus:      http.StatusNotFound,
			expectedBody:        "404: Not Found",
			expectedVampireID:   uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedDescription: "a description",
		},
		{
			name: "error from creator",
			body: url.Values{
				"description": []string{"a description"},
			},
			creator: &mockMarkCreator{
				err: errors.New("mock error"),
			},
			path:                "/vampires/12345678-90ab-cdef-1234-567890abcdef/marks",
			expectedStatus:      http.StatusInternalServerError,
			expectedBody:        "500: Internal Server Error",
			expectedVampireID:   uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedDescription: "a description",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewMux()

			handlers.CreateMark(r, testLogger(t), tt.creator)

			status, headers, body := post(r, tt.path, tt.body.Encode())

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedBody != body {
				t.Errorf("expected body %q; got %q", tt.expectedBody, body)
			}

			if tt.expectedVampireID != tt.creator.vampireID {
				t.Errorf("expected creator to receive vampire ID %q; got %q", tt.expectedVampireID, tt.creator.vampireID)
			}

			if tt.expectedDescription != tt.creator.description {
				t.Errorf("expected creator to receive description %q; got %q", tt.expectedDescription, tt.creator.description)
			}

			location := headers.Get("Location")
			if tt.expectedLocation != location {
				t.Errorf("expected location %q; got %q", tt.expectedLocation, location)
			}
		})
	}
}
