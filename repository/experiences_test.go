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

func TestCreateExperience_LimitToThree(t *testing.T) {
	tests := []struct {
		name                string
		descriptions        []string
		expectedExperiences []models.Experience
		expectedErrors      []error
	}{
		{
			name: "successful",
			descriptions: []string{
				"Experience 1",
			},
			expectedExperiences: []models.Experience{
				{
					Description: "Experience 1",
				},
			},
			expectedErrors: []error{
				nil,
			},
		},
		{
			name: "successful three times",
			descriptions: []string{
				"Experience 1",
				"Experience 2",
				"Experience 3",
			},
			expectedExperiences: []models.Experience{
				{
					Description: "Experience 1",
				},
				{
					Description: "Experience 2",
				},
				{
					Description: "Experience 3",
				},
			},
			expectedErrors: []error{
				nil,
				nil,
				nil,
			},
		},
		{
			name: "failure on fourth experience",
			descriptions: []string{
				"Experience 1",
				"Experience 2",
				"Experience 3",
				"Experience 4",
			},
			expectedExperiences: []models.Experience{
				{
					Description: "Experience 1",
				},
				{
					Description: "Experience 2",
				},
				{
					Description: "Experience 3",
				},
			},
			expectedErrors: []error{
				nil,
				nil,
				nil,
				models.ErrMemoryFull,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			m := newTestRepository(t)

			vampire, err := m.CreateVampire(context.Background(), "test vampire")
			if err != nil {
				t.Fatal(err)
			}

			var actualErrors []error
			memory := vampire.Memories[0]
			for _, description := range tt.descriptions {
				err := m.WithSavepoint(func(m *repository.Repository) error {
					_, err := m.CreateExperience(context.Background(), vampire.ID, memory.ID, description)
					return err
				})
				actualErrors = append(actualErrors, err)
			}

			if diff := cmp.Diff(tt.expectedErrors, actualErrors, cmpopts.EquateErrors()); diff != "" {
				t.Error(diff)
			}

			actualExperiences, err := m.GetExperiences(context.Background(), vampire.ID)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.expectedExperiences, actualExperiences, cmpopts.IgnoreFields(models.Experience{}, "ID", "MemoryID")); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestCreateExperience(t *testing.T) {
	tests := []struct {
		name          string
		ids           func(models.Vampire) (vampireID uuid.UUID, memoryID uuid.UUID)
		expectedError error
	}{
		{
			name:          "successful",
			ids:           func(v models.Vampire) (uuid.UUID, uuid.UUID) { return v.ID, v.Memories[0].ID },
			expectedError: nil,
		},
		{
			name:          "vampire not found",
			ids:           func(v models.Vampire) (uuid.UUID, uuid.UUID) { return uuid.New(), v.Memories[0].ID },
			expectedError: models.ErrNotFound,
		},
		{
			name:          "memory not found",
			ids:           func(v models.Vampire) (uuid.UUID, uuid.UUID) { return v.ID, uuid.New() },
			expectedError: models.ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			m := newTestRepository(t)

			vampire, err := m.CreateVampire(context.Background(), "test vampire")
			if err != nil {
				t.Fatal(err)
			}

			vampireID, memoryID := tt.ids(vampire)

			_, err = m.CreateExperience(context.Background(), vampireID, memoryID, "test description")
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected %q; received %q", tt.expectedError, err)
			}
		})
	}
}
