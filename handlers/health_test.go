package handlers_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/health"
	"github.com/go-chi/chi/v5"
)

type mockChecker struct {
	result health.Result
	ok     bool
}

func (m *mockChecker) Check(_ context.Context) (health.Result, bool) {
	return m.result, m.ok
}

func TestHealth(t *testing.T) {
	t.Parallel()

	now := time.Date(2021, time.December, 12, 11, 17, 0, 0, time.UTC)

	tests := []struct {
		name           string
		checker        *mockChecker
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful",
			checker: &mockChecker{
				result: health.Result{
					Details: []health.ComponentResult{
						{
							Name:      "mock component",
							Status:    "ok",
							Timestamp: now,
						},
					},
					Status: "ok",
				},
				ok: true,
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
  "details": [
    {
      "name": "mock component",
      "status": "ok",
      "timestamp": "2021-12-12T11:17:00Z"
    }
  ],
  "status": "ok"
}`,
		},
		{
			name: "failure from checker",
			checker: &mockChecker{
				result: health.Result{
					Details: []health.ComponentResult{
						{
							Name:      "mock component",
							Status:    "failed",
							Error:     "mock error",
							Timestamp: now,
						},
					},
					Status: "failed",
				},
				ok: false,
			},
			expectedStatus: http.StatusBadGateway,
			expectedBody: `{
  "details": [
    {
      "name": "mock component",
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

			handlers.Health(r, testLogger(t), tt.checker)

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
