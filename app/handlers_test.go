package app

import (
	"context"
	"fmt"
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

			app := NewApp(TestConfig(t))

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

func TestListVampires(t *testing.T) {
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

			app := NewApp(TestConfig(t))

			request := httptest.NewRequest(http.MethodGet, "/vampires", nil)
			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)

			err := app.listVampires(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}

func TestNewVampire(t *testing.T) {
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

			app := NewApp(TestConfig(t))

			request := httptest.NewRequest(http.MethodGet, "/vampires/new", nil)
			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)

			err := app.newVampire(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}

func TestCreateVampire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedError  error
	}{
		{
			name: "successful",
			body: url.Values{
				"name": []string{"Gruffudd"},
			}.Encode(),
			expectedStatus: http.StatusSeeOther,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := NewApp(TestConfig(t))

			request := httptest.NewRequest(http.MethodPost, "/vampires", strings.NewReader(tt.body))
			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)

			err := app.createVampire(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}

func TestShowVampire(t *testing.T) {
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

			app := NewApp(TestConfig(t))

			vampire, err := app.Models.CreateVampire(context.Background(), "Gruffudd")
			if err != nil {
				t.Fatal(err)
			}

			url := fmt.Sprintf("/vampires/%s", vampire.ID.String())
			request := httptest.NewRequest(http.MethodGet, url, nil)
			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)
			ctx.SetPath(url)
			ctx.SetParamNames("id")
			ctx.SetParamValues(vampire.ID.String())

			err = app.showVampire(ctx)

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
		name             string
		vampire          *models.Vampire
		body             string
		expectedStatus   int
		expectedLocation string
		expectedVampire  *models.Vampire
		expectedError    error
	}{
		{
			name:    "successful",
			vampire: &models.Vampire{},
			body: url.Values{
				"name": []string{"Gruffudd"},
			}.Encode(),
			expectedStatus:   http.StatusFound,
			expectedLocation: "/",
			expectedVampire: &models.Vampire{
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

			app := NewApp(TestConfig(t))
			app.Vampire = tt.vampire

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

			if diff := cmp.Diff(tt.expectedVampire, app.Vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}

func TestDeleteMemory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		vampire          *models.Vampire
		body             string
		expectedStatus   int
		expectedLocation string
		expectedVampire  *models.Vampire
		expectedError    error
	}{
		{
			name: "successful",
			vampire: &models.Vampire{
				Memories: []*models.Memory{
					{
						ID: 1,
						Experiences: []models.Experience{
							"one",
						},
					},
				},
			},
			body: url.Values{
				"_method": []string{"DELETE"},
			}.Encode(),
			expectedStatus:   http.StatusFound,
			expectedLocation: "/",
			expectedVampire: &models.Vampire{
				Memories: []*models.Memory{
					{
						ID:          2,
						Experiences: []models.Experience{},
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

			app := NewApp(TestConfig(t))
			app.Vampire = tt.vampire

			request := httptest.NewRequest(http.MethodPost, "/memories/1", strings.NewReader(tt.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)
			ctx.SetPath("/memories/1")
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := app.deleteMemory(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if location := response.Header().Get(echo.HeaderLocation); tt.expectedLocation != location {
				t.Errorf("expected %q; got %q", tt.expectedLocation, location)
			}

			if diff := cmp.Diff(tt.expectedVampire, app.Vampire); diff != "" {
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
		name             string
		vampire          *models.Vampire
		body             string
		expectedStatus   int
		expectedLocation string
		expectedVampire  *models.Vampire
		expectedError    error
	}{
		{
			name: "successful",
			vampire: &models.Vampire{
				Memories: []*models.Memory{
					{
						ID: 1,
					},
				},
			},
			body: url.Values{
				"experience": []string{"one"},
			}.Encode(),
			expectedStatus:   http.StatusFound,
			expectedLocation: "/",
			expectedVampire: &models.Vampire{
				Memories: []*models.Memory{
					{
						ID: 1,
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

			app := NewApp(TestConfig(t))
			app.Vampire = tt.vampire

			request := httptest.NewRequest(http.MethodPost, "/memories/1/experiences", strings.NewReader(tt.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)
			ctx.SetPath("/memories/1/experiences")
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := app.createExperience(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if location := response.Header().Get(echo.HeaderLocation); tt.expectedLocation != location {
				t.Errorf("expected %q; got %q", tt.expectedLocation, location)
			}

			if diff := cmp.Diff(tt.expectedVampire, app.Vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}

func TestCreateSkill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		vampire          *models.Vampire
		body             string
		expectedStatus   int
		expectedLocation string
		expectedVampire  *models.Vampire
		expectedError    error
	}{
		{
			name:    "successful",
			vampire: &models.Vampire{},
			body: url.Values{
				"description": []string{"Basic agricultural practices"},
			}.Encode(),
			expectedStatus:   http.StatusFound,
			expectedLocation: "/",
			expectedVampire: &models.Vampire{
				Skills: []*models.Skill{
					{
						ID:          1,
						Description: "Basic agricultural practices",
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

			app := NewApp(TestConfig(t))
			app.Vampire = tt.vampire

			request := httptest.NewRequest(http.MethodPost, "/skills", strings.NewReader(tt.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)

			err := app.createSkill(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if location := response.Header().Get(echo.HeaderLocation); tt.expectedLocation != location {
				t.Errorf("expected %q; got %q", tt.expectedLocation, location)
			}

			if diff := cmp.Diff(tt.expectedVampire, app.Vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}

func TestUpdateSkill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		vampire          *models.Vampire
		body             string
		expectedStatus   int
		expectedLocation string
		expectedVampire  *models.Vampire
		expectedError    error
	}{
		{
			name: "successful",
			vampire: &models.Vampire{
				Skills: []*models.Skill{
					{
						ID:          1,
						Description: "Basic agricultural practices",
						Checked:     false,
					},
				},
			},
			body: url.Values{
				"_method": []string{"PATCH"},
				"checked": []string{"1"},
			}.Encode(),
			expectedStatus:   http.StatusFound,
			expectedLocation: "/",
			expectedVampire: &models.Vampire{
				Skills: []*models.Skill{
					{
						ID:          1,
						Description: "Basic agricultural practices",
						Checked:     true,
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

			app := NewApp(TestConfig(t))
			app.Vampire = tt.vampire

			request := httptest.NewRequest(http.MethodPost, "/skills/1", strings.NewReader(tt.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)
			ctx.SetPath("/skills/1")
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := app.updateSkill(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if location := response.Header().Get(echo.HeaderLocation); tt.expectedLocation != location {
				t.Errorf("expected %q; got %q", tt.expectedLocation, location)
			}

			if diff := cmp.Diff(tt.expectedVampire, app.Vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}

func TestCreateResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		vampire          *models.Vampire
		body             string
		expectedStatus   int
		expectedLocation string
		expectedVampire  *models.Vampire
		expectedError    error
	}{
		{
			name:    "successful",
			vampire: &models.Vampire{},
			body: url.Values{
				"description": []string{"Calweddyn Farm, rich but challenging soils"},
				"stationary":  []string{"1"},
			}.Encode(),
			expectedStatus:   http.StatusFound,
			expectedLocation: "/",
			expectedVampire: &models.Vampire{
				Resources: []*models.Resource{
					{
						ID:          1,
						Description: "Calweddyn Farm, rich but challenging soils",
						Stationary:  true,
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

			app := NewApp(TestConfig(t))
			app.Vampire = tt.vampire

			request := httptest.NewRequest(http.MethodPost, "/resources", strings.NewReader(tt.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)

			err := app.createResource(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if location := response.Header().Get(echo.HeaderLocation); tt.expectedLocation != location {
				t.Errorf("expected %q; got %q", tt.expectedLocation, location)
			}

			if diff := cmp.Diff(tt.expectedVampire, app.Vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}

func TestUpdateResource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		vampire          *models.Vampire
		body             string
		expectedStatus   int
		expectedLocation string
		expectedVampire  *models.Vampire
		expectedError    error
	}{
		{
			name: "successful",
			vampire: &models.Vampire{
				Resources: []*models.Resource{
					{
						ID:          1,
						Description: "Basic agricultural practices",
						Stationary:  true,
						Lost:        false,
					},
				},
			},
			body: url.Values{
				"_method": []string{"PATCH"},
				"lost":    []string{"1"},
			}.Encode(),
			expectedStatus:   http.StatusFound,
			expectedLocation: "/",
			expectedVampire: &models.Vampire{
				Resources: []*models.Resource{
					{
						ID:          1,
						Description: "Basic agricultural practices",
						Stationary:  true,
						Lost:        true,
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

			app := NewApp(TestConfig(t))
			app.Vampire = tt.vampire

			request := httptest.NewRequest(http.MethodPost, "/resources/1", strings.NewReader(tt.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)
			ctx.SetPath("/resources/1")
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := app.updateResource(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if location := response.Header().Get(echo.HeaderLocation); tt.expectedLocation != location {
				t.Errorf("expected %q; got %q", tt.expectedLocation, location)
			}

			if diff := cmp.Diff(tt.expectedVampire, app.Vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}

func TestCreateCharacter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		vampire          *models.Vampire
		body             string
		expectedStatus   int
		expectedLocation string
		expectedVampire  *models.Vampire
		expectedError    error
	}{
		{
			name:    "successful",
			vampire: &models.Vampire{},
			body: url.Values{
				"name":         []string{"Lord Othian"},
				"descriptor[]": []string{"English gentry visiting a cathedral in St. Davids"},
				"type":         []string{"immortal"},
			}.Encode(),
			expectedStatus:   http.StatusFound,
			expectedLocation: "/",
			expectedVampire: &models.Vampire{
				Characters: []*models.Character{
					{
						ID:   1,
						Name: "Lord Othian",
						Descriptors: []string{
							"English gentry visiting a cathedral in St. Davids",
						},
						Type: "immortal",
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

			app := NewApp(TestConfig(t))
			app.Vampire = tt.vampire

			request := httptest.NewRequest(http.MethodPost, "/characters", strings.NewReader(tt.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)

			err := app.createCharacter(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if location := response.Header().Get(echo.HeaderLocation); tt.expectedLocation != location {
				t.Errorf("expected %q; got %q", tt.expectedLocation, location)
			}

			if diff := cmp.Diff(tt.expectedVampire, app.Vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}

func TestUpdateCharacter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		vampire          *models.Vampire
		body             string
		expectedStatus   int
		expectedLocation string
		expectedVampire  *models.Vampire
		expectedError    error
	}{
		{
			name: "successful",
			vampire: &models.Vampire{
				Characters: []*models.Character{
					{
						ID:       1,
						Deceased: false,
					},
				},
			},
			body: url.Values{
				"_method":  []string{"PATCH"},
				"deceased": []string{"1"},
			}.Encode(),
			expectedStatus:   http.StatusFound,
			expectedLocation: "/",
			expectedVampire: &models.Vampire{
				Characters: []*models.Character{
					{
						ID:       1,
						Deceased: true,
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

			app := NewApp(TestConfig(t))
			app.Vampire = tt.vampire

			request := httptest.NewRequest(http.MethodPost, "/characters/1", strings.NewReader(tt.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)
			ctx.SetPath("/characters/1")
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := app.updateCharacter(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if location := response.Header().Get(echo.HeaderLocation); tt.expectedLocation != location {
				t.Errorf("expected %q; got %q", tt.expectedLocation, location)
			}

			if diff := cmp.Diff(tt.expectedVampire, app.Vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}

func TestCreateDescriptor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		vampire          *models.Vampire
		body             string
		expectedStatus   int
		expectedLocation string
		expectedVampire  *models.Vampire
		expectedError    error
	}{
		{
			name: "successful",
			vampire: &models.Vampire{
				Characters: []*models.Character{
					{
						ID:          1,
						Descriptors: []string{},
					},
				},
			},
			body: url.Values{
				"descriptor": []string{"one"},
			}.Encode(),
			expectedStatus:   http.StatusFound,
			expectedLocation: "/",
			expectedVampire: &models.Vampire{
				Characters: []*models.Character{
					{
						ID:          1,
						Descriptors: []string{"one"},
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

			app := NewApp(TestConfig(t))
			app.Vampire = tt.vampire

			request := httptest.NewRequest(http.MethodPost, "/characters/1/descriptor", strings.NewReader(tt.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)
			ctx.SetPath("/characters/1/descriptor")
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			err := app.createDescriptor(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if location := response.Header().Get(echo.HeaderLocation); tt.expectedLocation != location {
				t.Errorf("expected %q; got %q", tt.expectedLocation, location)
			}

			if diff := cmp.Diff(tt.expectedVampire, app.Vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}

func TestCreateMark(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		vampire          *models.Vampire
		body             string
		expectedStatus   int
		expectedLocation string
		expectedVampire  *models.Vampire
		expectedError    error
	}{
		{
			name:    "successful",
			vampire: &models.Vampire{},
			body: url.Values{
				"description": []string{"Muddy footprints, muddy handprints, muddy sheets - I leave a trail of dirt wherever I travel."},
			}.Encode(),
			expectedStatus:   http.StatusFound,
			expectedLocation: "/",
			expectedVampire: &models.Vampire{
				Marks: []*models.Mark{
					{
						ID:          1,
						Description: "Muddy footprints, muddy handprints, muddy sheets - I leave a trail of dirt wherever I travel.",
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

			app := NewApp(TestConfig(t))
			app.Vampire = tt.vampire

			request := httptest.NewRequest(http.MethodPost, "/marks", strings.NewReader(tt.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			response := httptest.NewRecorder()
			ctx := app.NewContext(request, response)

			err := app.createMark(ctx)

			if tt.expectedStatus != response.Code {
				t.Errorf("expected %d; got %d", tt.expectedStatus, response.Code)
			}

			if location := response.Header().Get(echo.HeaderLocation); tt.expectedLocation != location {
				t.Errorf("expected %q; got %q", tt.expectedLocation, location)
			}

			if diff := cmp.Diff(tt.expectedVampire, app.Vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %v; got %v", tt.expectedError, err)
			}
		})
	}
}
