package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"emailaddress.horse/thousand/handlers"
	"github.com/go-chi/chi/v5"
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

			r := chi.NewMux()

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()

			handlers.Root(r)

			r.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			actualLocation := response.Header().Get("Location")
			if tt.expectedLocation != actualLocation {
				t.Errorf("expected %q; got %q", tt.expectedLocation, actualLocation)
			}
		})
	}
}
