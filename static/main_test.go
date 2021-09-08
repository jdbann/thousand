package static

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		headers        http.Header
		path           string
		expectedStatus int
		expectedError  error
	}{
		{
			name:           "success",
			path:           "/css/main.css",
			expectedStatus: http.StatusOK,
		},
		{
			name: "success with etag",
			headers: http.Header{
				"If-None-Match": []string{etagMap["css/main.css"]},
			},
			path:           "/css/main.css",
			expectedStatus: http.StatusNotModified,
		},
		{
			name: "success with incorrect etag",
			headers: http.Header{
				"If-None-Match": []string{"horses"},
			},
			path:           "/css/main.css",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Use(Middleware())

			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			request.Header = tt.headers

			response := httptest.NewRecorder()

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}
		})
	}
}
