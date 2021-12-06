package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/templates"
	"github.com/labstack/echo/v4"
)

func TestRoot(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		expectedStatus   int
		expectedLocation string
	}{
		{
			name:             "successful",
			expectedStatus:   http.StatusSeeOther,
			expectedLocation: "/vampires",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Renderer = templates.NewRenderer(e)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()

			handlers.Root(e)
			handlers.ListVampires(e, nil)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			actualLocation := response.Header().Get(echo.HeaderLocation)
			if tt.expectedLocation != actualLocation {
				t.Errorf("expected %q; got %q", tt.expectedLocation, actualLocation)
			}
		})
	}
}
