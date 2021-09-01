package models

import "testing"

func TestDescription(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		character           *Character
		expectedDescription string
	}{
		{
			name: "with name",
			character: &Character{
				Name: "John",
			},
			expectedDescription: "John.",
		},
		{
			name: "with name and descriptor",
			character: &Character{
				Name: "John",
				Descriptors: []string{
					"a programmer",
				},
			},
			expectedDescription: "John, a programmer.",
		},
		{
			name: "with name and multiple descriptors",
			character: &Character{
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
			character: &Character{
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
			character: &Character{
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
