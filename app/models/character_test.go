package models

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDescription(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		character           *OldCharacter
		expectedDescription string
	}{
		{
			name: "with name",
			character: &OldCharacter{
				Name: "John",
			},
			expectedDescription: "John.",
		},
		{
			name: "with name and descriptor",
			character: &OldCharacter{
				Name: "John",
				Descriptors: []string{
					"a programmer",
				},
			},
			expectedDescription: "John, a programmer.",
		},
		{
			name: "with name and multiple descriptors",
			character: &OldCharacter{
				Name: "John",
				Descriptors: []string{
					"a programmer",
					"struggling to sleep",
				},
			},
			expectedDescription: "John, a programmer, struggling to sleep.",
		},
		{
			name: "with name, multiple descriptors and a type",
			character: &OldCharacter{
				Name: "John",
				Descriptors: []string{
					"a programmer",
					"struggling to sleep",
				},
				Type: "mortal",
			},
			expectedDescription: "John, a programmer, struggling to sleep. (Mortal)",
		},
		{
			name: "with name and a type",
			character: &OldCharacter{
				Name: "John",
				Type: "mortal",
			},
			expectedDescription: "John. (Mortal)",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			description := tt.character.Description()

			if tt.expectedDescription != description {
				t.Errorf("wanted %q; got %q", tt.expectedDescription, description)
			}
		})
	}
}

func TestAddDescriptor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		character         *OldCharacter
		descriptor        string
		expectedCharacter *OldCharacter
	}{
		{
			name:       "success",
			character:  &OldCharacter{},
			descriptor: "one",
			expectedCharacter: &OldCharacter{
				Descriptors: []string{"one"},
			},
		},
		{
			name: "success with existing descriptors",
			character: &OldCharacter{
				Descriptors: []string{"one"},
			},
			descriptor: "two",
			expectedCharacter: &OldCharacter{
				Descriptors: []string{"one", "two"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			tt.character.AddDescriptor(tt.descriptor)

			if diff := cmp.Diff(tt.expectedCharacter, tt.character); diff != "" {
				t.Error(diff)
			}
		})
	}
}
