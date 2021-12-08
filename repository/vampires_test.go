package repository_test

import (
	"context"
	"errors"
	"testing"

	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/repository"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestCreateVampire(t *testing.T) {
	tests := []struct {
		name            string
		vampireName     string
		expectedVampire models.Vampire
	}{
		{
			name:        "successful",
			vampireName: "Gruffudd",
			expectedVampire: models.Vampire{
				Name:       "Gruffudd",
				Memories:   []models.Memory{},
				Skills:     []models.Skill{},
				Resources:  []models.Resource{},
				Characters: []models.Character{},
				Marks:      []models.Mark{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			m := newTestRepository(t)

			actualVampire, err := m.CreateVampire(context.Background(), tt.vampireName)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.expectedVampire, actualVampire, cmpopts.IgnoreFields(models.Vampire{}, "ID", "Memories")); diff != "" {
				t.Error(diff)
			}

			if len(actualVampire.Memories) != models.VampireMemorySize {
				t.Errorf("expected %d memories; found %d", models.VampireMemorySize, len(actualVampire.Memories))
			}
		})
	}
}

func TestGetVampire(t *testing.T) {
	tests := []struct {
		name          string
		id            func(models.Vampire) uuid.UUID
		expectedError error
	}{
		{
			name:          "successful",
			id:            func(v models.Vampire) uuid.UUID { return v.ID },
			expectedError: nil,
		},
		{
			name:          "not found",
			id:            func(v models.Vampire) uuid.UUID { return uuid.New() },
			expectedError: models.ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			m := newTestRepository(t)
			err := m.WithSavepoint(func(m *repository.Repository) error {
				vampire, err := m.CreateVampire(context.Background(), "test vampire")
				if err != nil {
					t.Fatal(err)
				}

				id := tt.id(vampire)

				_, err = m.GetVampire(context.Background(), id)
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("expected %q; received %q", tt.expectedError, err)
				}

				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
