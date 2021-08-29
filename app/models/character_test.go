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

func TestAddSkill(t *testing.T) {
	tests := []struct {
		name              string
		character         *Character
		skill             *Skill
		expectedCharacter *Character
	}{
		{
			name:      "success with no skills",
			character: &Character{},
			skill: &Skill{
				Description: "one",
			},
			expectedCharacter: &Character{
				Skills: []Skill{
					Skill{
						ID:          1,
						Description: "one",
					},
				},
			},
		},
		{
			name: "success with existing skills",
			character: &Character{
				Skills: []Skill{
					Skill{
						ID:          1,
						Description: "one",
					},
				},
			},
			skill: &Skill{
				Description: "two",
			},
			expectedCharacter: &Character{
				Skills: []Skill{
					Skill{
						ID:          1,
						Description: "one",
					},
					Skill{
						ID:          2,
						Description: "two",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.character.AddSkill(tt.skill)

			if diff := cmp.Diff(tt.expectedCharacter, tt.character); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestFindSkill(t *testing.T) {
	tests := []struct {
		name          string
		character     *Character
		skillID       int
		expectedSkill *Skill
		expectedError error
	}{
		{
			name: "success",
			character: &Character{
				Skills: []Skill{
					Skill{
						ID:          1,
						Description: "one",
					},
				},
			},
			skillID: 1,
			expectedSkill: &Skill{
				ID:          1,
				Description: "one",
			},
		},
		{
			name:          "failure with unknown ID",
			character:     &Character{},
			skillID:       1,
			expectedError: ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			skill, err := tt.character.FindSkill(tt.skillID)

			if diff := cmp.Diff(tt.expectedSkill, skill); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %s; actual %s", tt.expectedError, err)
			}
		})
	}
}

func TestUpdateSkill(t *testing.T) {
	tests := []struct {
		name              string
		character         *Character
		skill             *Skill
		expectedCharacter *Character
		expectedError     error
	}{
		{
			name: "success",
			character: &Character{
				Skills: []Skill{
					Skill{
						ID:          1,
						Description: "one",
						Checked:     false,
					},
					Skill{
						ID:          2,
						Description: "one",
						Checked:     true,
					},
				},
			},
			skill: &Skill{
				ID:          1,
				Description: "one",
				Checked:     true,
			},
			expectedCharacter: &Character{
				Skills: []Skill{
					Skill{
						ID:          1,
						Description: "one",
						Checked:     true,
					},
					Skill{
						ID:          2,
						Description: "one",
						Checked:     true,
					},
				},
			},
		},
		{
			name:      "failure with unknown skill",
			character: &Character{},
			skill: &Skill{
				ID:          1,
				Description: "one",
				Checked:     true,
			},
			expectedCharacter: &Character{},
			expectedError:     ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.character.UpdateSkill(tt.skill)

			if diff := cmp.Diff(tt.expectedCharacter, tt.character); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %s; actual %s", tt.expectedError, err)
			}
		})
	}
}
