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

type mockNewCharacterRenderer struct {
	err error
}

func (m *mockNewCharacterRenderer) NewCharacter(w http.ResponseWriter, _ *http.Request, vampire models.Vampire) error {
	if m.err != nil {
		return m.err
	}

	_, err := w.Write([]byte(vampire.Name))
	if err != nil {
		panic(err)
	}

	return nil
}

func TestNewCharacter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		renderer       *mockNewCharacterRenderer
		getter         *mockVampireGetter
		path           string
		expectedStatus int
		expectedBody   string
		expectedID     uuid.UUID
	}{
		{
			name:     "successful",
			renderer: &mockNewCharacterRenderer{},
			getter: &mockVampireGetter{
				vampire: models.Vampire{
					Name: "a vampire",
				},
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/characters/new",
			expectedStatus: http.StatusOK,
			expectedBody:   "a vampire",
			expectedID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		{
			name:           "error parsing vampire id",
			getter:         &mockVampireGetter{},
			path:           "/vampires/unknown/characters/new",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name:     "not found from getter",
			renderer: &mockNewCharacterRenderer{},
			getter: &mockVampireGetter{
				err: models.ErrNotFound,
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/characters/new",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404: Not Found",
			expectedID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		{
			name:     "error from getter",
			renderer: &mockNewCharacterRenderer{},
			getter: &mockVampireGetter{
				err: models.ErrNotFound,
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/characters/new",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404: Not Found",
			expectedID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		{
			name: "error from renderer",
			renderer: &mockNewCharacterRenderer{
				err: models.ErrNotFound,
			},
			getter: &mockVampireGetter{
				vampire: models.Vampire{
					Name: "a vampire",
				},
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/characters/new",
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

			handlers.NewCharacter(r, testLogger(t), tt.renderer, tt.getter)

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

type mockCharacterCreator struct {
	vampireID uuid.UUID
	params    models.CreateCharacterParams
	err       error
}

func (m *mockCharacterCreator) CreateCharacter(_ context.Context, vampireID uuid.UUID, params models.CreateCharacterParams) (models.Character, error) {
	m.vampireID = vampireID
	m.params = params
	return models.Character{}, m.err
}

func TestCreateCharacter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		body              url.Values
		creator           *mockCharacterCreator
		path              string
		expectedStatus    int
		expectedBody      string
		expectedLocation  string
		expectedVampireID uuid.UUID
		expectedParams    models.CreateCharacterParams
	}{
		{
			name: "successful",
			body: url.Values{
				"name": []string{"a name"},
				"type": []string{"mortal"},
			},
			creator:           &mockCharacterCreator{},
			path:              "/vampires/12345678-90ab-cdef-1234-567890abcdef/characters",
			expectedStatus:    http.StatusSeeOther,
			expectedLocation:  "/vampires/12345678-90ab-cdef-1234-567890abcdef",
			expectedVampireID: uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedParams: models.CreateCharacterParams{
				Name: "a name",
				Type: "mortal",
			},
		},
		{
			name: "error parsing vampire ID",
			body: url.Values{
				"name": []string{"a name"},
				"type": []string{"mortal"},
			},
			creator:        &mockCharacterCreator{},
			path:           "/vampires/unknown/characters",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name: "not found from creator",
			body: url.Values{
				"name": []string{"a name"},
				"type": []string{"mortal"},
			},
			creator: &mockCharacterCreator{
				err: models.ErrNotFound,
			},
			path:              "/vampires/12345678-90ab-cdef-1234-567890abcdef/characters",
			expectedStatus:    http.StatusNotFound,
			expectedBody:      "404: Not Found",
			expectedVampireID: uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedParams: models.CreateCharacterParams{
				Name: "a name",
				Type: "mortal",
			},
		},
		{
			name: "error from creator",
			body: url.Values{
				"name": []string{"a name"},
				"type": []string{"mortal"},
			},
			creator: &mockCharacterCreator{
				err: errors.New("mock error"),
			},
			path:              "/vampires/12345678-90ab-cdef-1234-567890abcdef/characters",
			expectedStatus:    http.StatusInternalServerError,
			expectedBody:      "500: Internal Server Error",
			expectedVampireID: uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedParams: models.CreateCharacterParams{
				Name: "a name",
				Type: "mortal",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewMux()

			handlers.CreateCharacter(r, testLogger(t), tt.creator)

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

			if tt.expectedParams != tt.creator.params {
				t.Errorf("expected creator to receive params %+v; got %+v", tt.expectedParams, tt.creator.params)
			}

			location := headers.Get("Location")
			if tt.expectedLocation != location {
				t.Errorf("expected location %q; got %q", tt.expectedLocation, location)
			}
		})
	}
}
