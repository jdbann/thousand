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

type mockNewResourceRenderer struct {
	err error
}

func (m *mockNewResourceRenderer) NewResource(w http.ResponseWriter, vampire models.Vampire) error {
	if m.err != nil {
		return m.err
	}

	_, err := w.Write([]byte(vampire.Name))
	if err != nil {
		panic(err)
	}

	return nil
}

func TestNewResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		renderer       *mockNewResourceRenderer
		getter         *mockVampireGetter
		path           string
		expectedStatus int
		expectedBody   string
		expectedID     uuid.UUID
	}{
		{
			name:     "successful",
			renderer: &mockNewResourceRenderer{},
			getter: &mockVampireGetter{
				vampire: models.Vampire{
					Name: "a vampire",
				},
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/resources/new",
			expectedStatus: http.StatusOK,
			expectedBody:   "a vampire",
			expectedID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		{
			name:           "error parsing vampire id",
			getter:         &mockVampireGetter{},
			path:           "/vampires/unknown/resources/new",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "500: Internal Server Error",
		},
		{
			name:     "not found from getter",
			renderer: &mockNewResourceRenderer{},
			getter: &mockVampireGetter{
				err: models.ErrNotFound,
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/resources/new",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404: Not Found",
			expectedID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		{
			name:     "error from getter",
			renderer: &mockNewResourceRenderer{},
			getter: &mockVampireGetter{
				err: models.ErrNotFound,
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/resources/new",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404: Not Found",
			expectedID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		{
			name: "error from renderer",
			renderer: &mockNewResourceRenderer{
				err: models.ErrNotFound,
			},
			getter: &mockVampireGetter{
				vampire: models.Vampire{
					Name: "a vampire",
				},
			},
			path:           "/vampires/11111111-1111-1111-1111-111111111111/resources/new",
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

			handlers.NewResource(r, testLogger(t), tt.renderer, tt.getter)

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
			e.Renderer = templates.NewEchoRenderer(e)

			request := httptest.NewRequest(http.MethodPost, "/vampires/12345678-90ab-cdef-1234-567890abcdef/resources", strings.NewReader(tt.body.Encode()))
			request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationForm)
			response := httptest.NewRecorder()

			handlers.CreateResource(e, tt.resourceCreator)

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
