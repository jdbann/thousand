package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"emailaddress.horse/thousand/handlers"
	"github.com/go-chi/chi/v5"
)

type mockPinger struct {
	err error
}

func (m *mockPinger) Ping(_ context.Context) error {
	return m.err
}

func TestHealth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		pinger         *mockPinger
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful",
			pinger:         &mockPinger{},
			expectedStatus: http.StatusOK,
			expectedBody: `{
  "details": [
    {
      "name": "repository",
      "status": "ok",
      "timestamp": "2021-12-12T11:17:00Z"
    }
  ],
  "status": "ok"
}`,
		},
		{
			name: "error from pinger",
			pinger: &mockPinger{
				err: errors.New("mock error"),
			},
			expectedStatus: http.StatusBadGateway,
			expectedBody: `{
  "details": [
    {
      "name": "repository",
      "status": "failed",
      "error": "mock error",
      "timestamp": "2021-12-12T11:17:00Z"
    }
  ],
  "status": "failed"
}`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewMux()

			now := func() time.Time {
				return time.Date(2021, time.December, 12, 11, 17, 0, 0, time.UTC)
			}

			handlers.Health(r, testLogger(t), tt.pinger, now)

			status, _, body := get(r, "/health")

			if tt.expectedStatus != status {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, status)
			}

			if tt.expectedBody != body {
				t.Errorf("expected body:\n%s\n\ngot:\n%s", tt.expectedBody, body)
			}
		})
	}
}
