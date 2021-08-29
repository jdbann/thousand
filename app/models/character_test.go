package models

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCharacter_AddExperience(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		character         *Character
		memoryID          int
		experienceString  string
		expectedCharacter *Character
		expectedError     error
	}{
		{
			name:             "success with recognised memoryID",
			character:        &Character{},
			memoryID:         0,
			experienceString: "one",
			expectedCharacter: &Character{
				Memories: [5]Memory{
					Memory{Experiences: []Experience{Experience("one")}},
				},
			},
			expectedError: nil,
		},
		{
			name: "failure with memoryID for full memory",
			character: &Character{
				Memories: [5]Memory{
					Memory{Experiences: []Experience{
						Experience("one"),
						Experience("two"),
						Experience("three"),
					}},
				},
			},
			memoryID:         0,
			experienceString: "four",
			expectedCharacter: &Character{
				Memories: [5]Memory{
					Memory{Experiences: []Experience{
						Experience("one"),
						Experience("two"),
						Experience("three"),
					}},
				},
			},
			expectedError: ErrMemoryFull,
		},
		{
			name:              "failure with negative memoryID",
			character:         &Character{},
			memoryID:          -1,
			experienceString:  "one",
			expectedCharacter: &Character{},
			expectedError:     ErrNotFound,
		},
		{
			name:              "failure with unrecognised memoryID",
			character:         &Character{},
			memoryID:          5,
			experienceString:  "one",
			expectedCharacter: &Character{},
			expectedError:     ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.character.AddExperience(tt.memoryID, tt.experienceString)

			if diff := cmp.Diff(tt.expectedCharacter, tt.character); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %s; actual %s", tt.expectedError, err)
			}
		})
	}
}
