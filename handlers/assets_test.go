package handlers_test

import (
	"io/fs"
	"net/http"
	"testing"
	"testing/fstest"

	"emailaddress.horse/thousand/handlers"
	"github.com/go-chi/chi/v5"
)

func TestAssets(t *testing.T) {
	t.Parallel()

	testFS := fstest.MapFS{
		"css": {
			Mode: fs.ModeDir,
			Data: []byte("beautiful css"),
		},
		"css/main.css": {
			Data: []byte("beautiful css"),
		},
	}

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful",
			path:           "/assets/css/main.css",
			expectedStatus: http.StatusOK,
			expectedBody:   "beautiful css",
		},
		{
			name:           "not found for directories",
			path:           "/assets/css",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404: Not Found",
		},
		{
			name:           "not found for unknown files",
			path:           "/assets/css/ugly.css",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404: Not Found",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewMux()

			handlers.Assets(r, testFS)

			status, _, body := get(r, tt.path)

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedBody != body {
				t.Errorf("expected body %q; got %q", tt.expectedBody, body)
			}
		})
	}
}
