package app

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"emailaddress.horse/thousand/app/models"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/echo/v4"
)

func TestRoot(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		expectedStatus int
		expectedError  error
	}{
		{
			name:           "successful",
			expectedStatus: http.StatusOK,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := NewApp()

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)

			err := app.root(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}

func TestCreateDetails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		character         *models.Character
		body              string
		expectedStatus    int
		expectedLocation  string
		expectedCharacter *models.Character
		expectedError     error
	}{
		{
			name:      "successful",
			character: &models.Character{},
			body: url.Values{
				"name": []string{"Gruffudd"},
			}.Encode(),
			expectedStatus:   http.StatusFound,
			expectedLocation: "/",
			expectedCharacter: &models.Character{
				Details: &models.Details{
					Name: "Gruffudd",
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := NewApp()
			app.Character = tt.character

			request := httptest.NewRequest(http.MethodPost, "/details", strings.NewReader(tt.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)

			err := app.createDetails(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if location := response.Header().Get(echo.HeaderLocation); tt.expectedLocation != location {
				t.Errorf("expected %q; got %q", tt.expectedLocation, location)
			}

			if diff := cmp.Diff(tt.expectedCharacter, app.Character); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}

func TestCreateExperience(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		character         *models.Character
		body              string
		expectedStatus    int
		expectedLocation  string
		expectedCharacter *models.Character
		expectedError     error
	}{
		{
			name:      "successful",
			character: &models.Character{},
			body: url.Values{
				"experience": []string{"one"},
			}.Encode(),
			expectedStatus:   http.StatusFound,
			expectedLocation: "/",
			expectedCharacter: &models.Character{
				Memories: [5]models.Memory{
					{
						Experiences: []models.Experience{
							"one",
						},
					},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := NewApp()
			app.Character = tt.character

			request := httptest.NewRequest(http.MethodPost, "/memories/0/experiences", strings.NewReader(tt.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)
			ctx.SetPath("/memories/0/experiences")
			ctx.SetParamNames("id")
			ctx.SetParamValues("0")

			err := app.createExperience(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if location := response.Header().Get(echo.HeaderLocation); tt.expectedLocation != location {
				t.Errorf("expected %q; got %q", tt.expectedLocation, location)
			}

			if diff := cmp.Diff(tt.expectedCharacter, app.Character); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}
