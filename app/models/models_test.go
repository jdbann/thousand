package models

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestCreateVampire(t *testing.T) {
	t.Parallel()

	ignoreFields := cmp.Options{
		ignoreVampireFields,
	}

	tests := []struct {
		name            string
		vampireName     string
		expectedVampire Vampire
	}{
		{
			name:        "successful",
			vampireName: "Gruffudd",
			expectedVampire: Vampire{
				Name:       "Gruffudd",
				Memories:   []Memory{},
				Skills:     []Skill{},
				Resources:  []Resource{},
				Characters: []Character{},
				Marks:      []Mark{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newTestModels(t)

			actualVampire, err := m.CreateVampire(context.Background(), tt.vampireName)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.expectedVampire, actualVampire, ignoreFields...); diff != "" {
				t.Error(diff)
			}

			if len(actualVampire.Memories) != vampireMemorySize {
				t.Errorf("expected %d memories; found %d", vampireMemorySize, len(actualVampire.Memories))
			}
		})
	}
}

func TestGetVampire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		id            func(Vampire) uuid.UUID
		expectedError error
	}{
		{
			name:          "successful",
			id:            func(v Vampire) uuid.UUID { return v.ID },
			expectedError: nil,
		},
		{
			name:          "not found",
			id:            func(v Vampire) uuid.UUID { return uuid.New() },
			expectedError: ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := newTestModels(t)
			err := m.WithSavepoint(func(m *Models) error {
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

func TestAddExperience(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		descriptions        []string
		expectedExperiences []Experience
		expectedErrors      []error
	}{
		{
			name: "successful",
			descriptions: []string{
				"Experience 1",
			},
			expectedExperiences: []Experience{
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
			expectedExperiences: []Experience{
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
			expectedExperiences: []Experience{
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

			if diff := cmp.Diff(tt.expectedExperiences, actualExperiences, ignoreExperienceFields); diff != "" {
				t.Error(diff)
			}
		})
	}
}
