package repository_test

import (
	"context"
	"errors"
	"testing"

	"emailaddress.horse/thousand/app/models"
	"emailaddress.horse/thousand/repository"
	"github.com/google/uuid"
)

func TestGetMemory(t *testing.T) {
	tests := []struct {
		name          string
		id            func(models.Vampire) (vampireID uuid.UUID, id uuid.UUID)
		expectedError error
	}{
		{
			name:          "successful",
			id:            func(v models.Vampire) (uuid.UUID, uuid.UUID) { return v.ID, v.Memories[0].ID },
			expectedError: nil,
		},
		{
			name:          "vampire not found",
			id:            func(v models.Vampire) (uuid.UUID, uuid.UUID) { return uuid.New(), v.Memories[0].ID },
			expectedError: models.ErrNotFound,
		},
		{
			name:          "memory not found",
			id:            func(v models.Vampire) (uuid.UUID, uuid.UUID) { return v.ID, uuid.New() },
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

				vampireID, id := tt.id(vampire)

				_, err = m.GetMemory(context.Background(), vampireID, id)
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
