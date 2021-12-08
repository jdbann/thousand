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

func TestNewCharacter(t *testing.T) {
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
			e.Renderer = templates.NewEchoRenderer(e)

			request := httptest.NewRequest(http.MethodGet, "/vampires/12345678-90ab-cdef-1234-567890abcdef/characters/new", nil)
			response := httptest.NewRecorder()

			handlers.NewCharacter(e, tt.vampireGetter)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if tt.expectedID != tt.vampireGetter.id {
				t.Errorf("expected %q; got %q", tt.expectedID, tt.vampireGetter.id)
			}
		})
	}
}

type mockCharacterCreator struct {
	receivedVampireID uuid.UUID
	receivedParams    models.CreateCharacterParams
	err               error
}

func (m *mockCharacterCreator) CreateCharacter(_ context.Context, vampireID uuid.UUID, params models.CreateCharacterParams) (models.Character, error) {
	m.receivedVampireID = vampireID
	m.receivedParams = params
	return models.Character{}, m.err
}

func TestCreateCharacter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		body              url.Values
		characterCreator  *mockCharacterCreator
		expectedStatus    int
		expectedVampireID uuid.UUID
		expectedParams    models.CreateCharacterParams
		expectedLocation  string
	}{
		{
			name: "successful",
			body: url.Values{
				"name": []string{"A name"},
				"type": []string{"mortal"},
			},
			characterCreator:  &mockCharacterCreator{},
			expectedStatus:    http.StatusSeeOther,
			expectedVampireID: uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedParams: models.CreateCharacterParams{
				Name: "A name",
				Type: "mortal",
			},
			expectedLocation: "/vampires/12345678-90ab-cdef-1234-567890abcdef",
		},
		{
			name: "not found",
			body: url.Values{
				"name": []string{"A name"},
				"type": []string{"mortal"},
			},
			characterCreator:  &mockCharacterCreator{err: models.ErrNotFound},
			expectedStatus:    http.StatusNotFound,
			expectedVampireID: uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedParams: models.CreateCharacterParams{
				Name: "A name",
				Type: "mortal",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Renderer = templates.NewEchoRenderer(e)

			request := httptest.NewRequest(http.MethodPost, "/vampires/12345678-90ab-cdef-1234-567890abcdef/characters", strings.NewReader(tt.body.Encode()))
			request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationForm)
			response := httptest.NewRecorder()

			handlers.CreateCharacter(e, tt.characterCreator)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if tt.expectedVampireID != tt.characterCreator.receivedVampireID {
				t.Errorf("expected %q; got %q", tt.expectedVampireID, tt.characterCreator.receivedVampireID)
			}

			if tt.expectedParams != tt.characterCreator.receivedParams {
				t.Errorf("expected %+v; got %+v", tt.expectedParams, tt.characterCreator.receivedParams)
			}

			actualLocation := response.Header().Get(echo.HeaderLocation)
			if tt.expectedLocation != actualLocation {
				t.Errorf("expected %q; got %q", tt.expectedLocation, actualLocation)
			}
		})
	}
}
