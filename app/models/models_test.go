package models

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreateVampire(t *testing.T) {
	t.Parallel()

	ignoreFields := cmp.Options{
		ignoreWholeVampireFields,
		ignoreNewVampireFields,
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
				[]Memory{},
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

func TestAddExperience(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		descriptions        []string
		expectedExperiences []NewExperience
		expectedErrors      []error
	}{
		{
			name: "successful",
			descriptions: []string{
				"Experience 1",
			},
			expectedExperiences: []NewExperience{
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
			expectedExperiences: []NewExperience{
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
			expectedExperiences: []NewExperience{
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
				ErrMemoryFull,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			m := newTestModels(t)

			vampire, err := m.CreateVampire(context.Background(), "test vampire")
			if err != nil {
				t.Fatal(err)
			}

			var actualErrors []error
			memory := vampire.Memories[0]
			for _, description := range tt.descriptions {
				err := m.WithSavepoint(func(m *Models) error {
					_, err := m.AddExperience(context.Background(), vampire.ID, memory.ID, description)
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

			if diff := cmp.Diff(tt.expectedExperiences, actualExperiences, ignoreNewExperienceFields); diff != "" {
				t.Error(diff)
			}
		})
	}
}
