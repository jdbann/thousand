package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"emailaddress.horse/thousand/app/models"
	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/templates"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type mockMemoryGetter struct {
	memory models.Memory
	err    error
}

func (m mockMemoryGetter) GetMemory(_ context.Context, _, _ uuid.UUID) (models.Memory, error) {
	return m.memory, m.err
}

func TestNewExperience(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		memoryGetter   mockMemoryGetter
		expectedStatus int
	}{
		{
			name:           "successful",
			memoryGetter:   mockMemoryGetter{memory: models.Memory{}},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not found",
			memoryGetter:   mockMemoryGetter{err: models.ErrNotFound},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Renderer = templates.NewRenderer(e)

			request := httptest.NewRequest(http.MethodGet, "/vampires/12345678-90ab-cdef-1234-567890abcdef/memories/12345678-90ab-cdef-1234-567890abcdef/experiences/new", nil)
			response := httptest.NewRecorder()

			handlers.NewExperience(e, tt.memoryGetter)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}
		})
	}
}
