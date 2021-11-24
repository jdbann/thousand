package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"emailaddress.horse/thousand/app/models"
	"emailaddress.horse/thousand/handlers"
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
