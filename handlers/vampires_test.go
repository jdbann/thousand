package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"emailaddress.horse/thousand/app/models"
	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/templates"
	"github.com/labstack/echo/v4"
)

type mockVampiresGetter struct{}

func (m mockVampiresGetter) GetVampires(_ context.Context) ([]models.Vampire, error) {
	return []models.Vampire{}, nil
}

func TestListVampires(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "successful",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Renderer = templates.NewRenderer(e)

			request := httptest.NewRequest(http.MethodGet, "/vampires", nil)
			response := httptest.NewRecorder()

			handlers.ListVampires(e, mockVampiresGetter{})

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}
		})
	}
}

func TestNewVampire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "successful",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Renderer = templates.NewRenderer(e)

			request := httptest.NewRequest(http.MethodGet, "/vampires/new", nil)
			response := httptest.NewRecorder()

			handlers.NewVampire(e)

			e.ServeHTTP(response, request)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}
		})
	}
}
