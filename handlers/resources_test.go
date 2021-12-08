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

func TestNewResource(t *testing.T) {
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

			request := httptest.NewRequest(http.MethodGet, "/vampires/12345678-90ab-cdef-1234-567890abcdef/resources/new", nil)
			response := httptest.NewRecorder()

			handlers.NewResource(e, tt.vampireGetter)

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

type mockResourceCreator struct {
	receivedVampireID uuid.UUID
	receivedParams    models.CreateResourceParams
	err               error
}

func (m *mockResourceCreator) CreateResource(_ context.Context, vampireID uuid.UUID, params models.CreateResourceParams) (models.Resource, error) {
	m.receivedVampireID = vampireID
	m.receivedParams = params
	return models.Resource{
		ID:          uuid.New(),
		VampireID:   vampireID,
		Description: m.receivedParams.Description,
		Stationary:  m.receivedParams.Stationary,
	}, m.err
}

func TestCreateResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		body              url.Values
		resourceCreator   *mockResourceCreator
		expectedStatus    int
		expectedVampireID uuid.UUID
		expectedParams    models.CreateResourceParams
		expectedLocation  string
	}{
		{
			name: "successful",
			body: url.Values{
				"description": []string{"A description"},
				"stationary":  []string{"1"},
			},
			resourceCreator:   &mockResourceCreator{},
			expectedStatus:    http.StatusSeeOther,
			expectedVampireID: uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedParams: models.CreateResourceParams{
				Description: "A description",
				Stationary:  true,
			},
			expectedLocation: "/vampires/12345678-90ab-cdef-1234-567890abcdef",
		},
		{
			name: "not found",
			body: url.Values{
				"description": []string{"A description"},
				"stationary":  []string{"1"},
			},
			resourceCreator:   &mockResourceCreator{err: models.ErrNotFound},
			expectedStatus:    http.StatusNotFound,
			expectedVampireID: uuid.MustParse("12345678-90ab-cdef-1234-567890abcdef"),
			expectedParams: models.CreateResourceParams{
				Description: "A description",
				Stationary:  true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Renderer = templates.NewRenderer(e)

			request := httptest.NewRequest(http.MethodPost, "/vampires/12345678-90ab-cdef-1234-567890abcdef/resources", strings.NewReader(tt.body.Encode()))
			request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationForm)
			response := httptest.NewRecorder()

			handlers.CreateResource(e, tt.resourceCreator)
			handlers.ShowVampire(e, nil)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if tt.expectedVampireID != tt.resourceCreator.receivedVampireID {
				t.Errorf("expected %q; got %q", tt.expectedVampireID, tt.resourceCreator.receivedVampireID)
			}

			if tt.expectedParams != tt.resourceCreator.receivedParams {
				t.Errorf("expected %+v; got %+v", tt.expectedParams, tt.resourceCreator.receivedParams)
			}

			actualLocation := response.Header().Get(echo.HeaderLocation)
			if tt.expectedLocation != actualLocation {
				t.Errorf("expected %q; got %q", tt.expectedLocation, actualLocation)
			}
		})
	}
}
