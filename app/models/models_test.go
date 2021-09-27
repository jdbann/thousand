package models

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreateVampire(t *testing.T) {
	t.Parallel()

	ignoreFields := []cmp.Option{
		// Memories when blank have no matchable attributes so need to be ignored
		cmpopts.IgnoreFields(WholeVampire{}, "Memories"),
		// ID and timestamps are defined by the DB and do not require matching
		cmpopts.IgnoreFields(NewVampire{}, "ID", "CreatedAt", "UpdatedAt"),
	}

	tests := []struct {
		name                 string
		vampireName          string
		expectedWholeVampire WholeVampire
	}{
		{
			name:        "successful",
			vampireName: "Gruffudd",
			expectedWholeVampire: WholeVampire{
				NewVampire{
					Name: "Gruffudd",
				},
				[]NewMemory{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newTestModels(t)

			actualWholeVampire, err := m.CreateVampire(context.Background(), tt.vampireName)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.expectedWholeVampire, actualWholeVampire, ignoreFields...); diff != "" {
				t.Error(diff)
			}

			if len(actualWholeVampire.Memories) != vampireMemorySize {
				t.Errorf("expected %d memories; found %d", vampireMemorySize, len(actualWholeVampire.Memories))
			}
		})
	}
}
