package repository_test

import (
	"context"
	"errors"
	"testing"

	"emailaddress.horse/thousand/models"
	"github.com/google/uuid"
)

func TestCreateMark(t *testing.T) {
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
			name:          "vampire not found",
			id:            func(v models.Vampire) uuid.UUID { return uuid.New() },
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

			_, err = m.CreateMark(context.Background(), tt.id(vampire), "A description")
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected %q; received %q", tt.expectedError, err)
			}
		})
	}
}
