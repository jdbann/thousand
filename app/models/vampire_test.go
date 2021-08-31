package models

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestVampire_AddExperience(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		vampire          *Vampire
		memoryID         int
		experienceString string
		expectedVampire  *Vampire
		expectedError    error
	}{
		{
			name:             "success with recognised memoryID",
			vampire:          &Vampire{},
			memoryID:         0,
			experienceString: "one",
			expectedVampire: &Vampire{
				Memories: [5]Memory{
					Memory{Experiences: []Experience{Experience("one")}},
				},
			},
			expectedError: nil,
		},
		{
			name: "failure with memoryID for full memory",
			vampire: &Vampire{
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
			expectedVampire: &Vampire{
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
			name:             "failure with negative memoryID",
			vampire:          &Vampire{},
			memoryID:         -1,
			experienceString: "one",
			expectedVampire:  &Vampire{},
			expectedError:    ErrNotFound,
		},
		{
			name:             "failure with unrecognised memoryID",
			vampire:          &Vampire{},
			memoryID:         5,
			experienceString: "one",
			expectedVampire:  &Vampire{},
			expectedError:    ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.vampire.AddExperience(tt.memoryID, tt.experienceString)

			if diff := cmp.Diff(tt.expectedVampire, tt.vampire); diff != "" {
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
		name            string
		vampire         *Vampire
		skill           *Skill
		expectedVampire *Vampire
	}{
		{
			name:    "success with no skills",
			vampire: &Vampire{},
			skill: &Skill{
				Description: "one",
			},
			expectedVampire: &Vampire{
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
			vampire: &Vampire{
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
			expectedVampire: &Vampire{
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

			tt.vampire.AddSkill(tt.skill)

			if diff := cmp.Diff(tt.expectedVampire, tt.vampire); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestFindSkill(t *testing.T) {
	tests := []struct {
		name          string
		vampire       *Vampire
		skillID       int
		expectedSkill *Skill
		expectedError error
	}{
		{
			name: "success",
			vampire: &Vampire{
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
			vampire:       &Vampire{},
			skillID:       1,
			expectedError: ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			skill, err := tt.vampire.FindSkill(tt.skillID)

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
		name            string
		vampire         *Vampire
		skill           *Skill
		expectedVampire *Vampire
		expectedError   error
	}{
		{
			name: "success",
			vampire: &Vampire{
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
			expectedVampire: &Vampire{
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
			name:    "failure with unknown skill",
			vampire: &Vampire{},
			skill: &Skill{
				ID:          1,
				Description: "one",
				Checked:     true,
			},
			expectedVampire: &Vampire{},
			expectedError:   ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.vampire.UpdateSkill(tt.skill)

			if diff := cmp.Diff(tt.expectedVampire, tt.vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %s; actual %s", tt.expectedError, err)
			}
		})
	}
}

func TestAddResource(t *testing.T) {
	tests := []struct {
		name            string
		vampire         *Vampire
		resource        *Resource
		expectedVampire *Vampire
	}{
		{
			name:    "success with no resources",
			vampire: &Vampire{},
			resource: &Resource{
				Description: "one",
			},
			expectedVampire: &Vampire{
				Resources: []Resource{
					Resource{
						ID:          1,
						Description: "one",
					},
				},
			},
		},
		{
			name: "success with existing resources",
			vampire: &Vampire{
				Resources: []Resource{
					Resource{
						ID:          1,
						Description: "one",
					},
				},
			},
			resource: &Resource{
				Description: "two",
				Stationary:  true,
			},
			expectedVampire: &Vampire{
				Resources: []Resource{
					Resource{
						ID:          1,
						Description: "one",
					},
					Resource{
						ID:          2,
						Description: "two",
						Stationary:  true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.vampire.AddResource(tt.resource)

			if diff := cmp.Diff(tt.expectedVampire, tt.vampire); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestFindResource(t *testing.T) {
	tests := []struct {
		name             string
		vampire          *Vampire
		resourceID       int
		expectedResource *Resource
		expectedError    error
	}{
		{
			name: "success",
			vampire: &Vampire{
				Resources: []Resource{
					Resource{
						ID:          1,
						Description: "one",
					},
				},
			},
			resourceID: 1,
			expectedResource: &Resource{
				ID:          1,
				Description: "one",
			},
		},
		{
			name:          "failure with unknown ID",
			vampire:       &Vampire{},
			resourceID:    1,
			expectedError: ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource, err := tt.vampire.FindResource(tt.resourceID)

			if diff := cmp.Diff(tt.expectedResource, resource); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %s; actual %s", tt.expectedError, err)
			}
		})
	}
}

func TestUpdateResource(t *testing.T) {
	tests := []struct {
		name            string
		vampire         *Vampire
		resource        *Resource
		expectedVampire *Vampire
		expectedError   error
	}{
		{
			name: "success",
			vampire: &Vampire{
				Resources: []Resource{
					Resource{
						ID:          1,
						Description: "one",
						Lost:        false,
					},
					Resource{
						ID:          2,
						Description: "one",
						Lost:        true,
					},
				},
			},
			resource: &Resource{
				ID:          1,
				Description: "one",
				Lost:        true,
			},
			expectedVampire: &Vampire{
				Resources: []Resource{
					Resource{
						ID:          1,
						Description: "one",
						Lost:        true,
					},
					Resource{
						ID:          2,
						Description: "one",
						Lost:        true,
					},
				},
			},
		},
		{
			name:    "failure with unknown resource",
			vampire: &Vampire{},
			resource: &Resource{
				ID:          1,
				Description: "one",
				Lost:        true,
			},
			expectedVampire: &Vampire{},
			expectedError:   ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.vampire.UpdateResource(tt.resource)

			if diff := cmp.Diff(tt.expectedVampire, tt.vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %s; actual %s", tt.expectedError, err)
			}
		})
	}
}

func TestAddCharacter(t *testing.T) {
	tests := []struct {
		name            string
		vampire         *Vampire
		character       *Character
		expectedVampire *Vampire
	}{
		{
			name:    "success with no characters",
			vampire: &Vampire{},
			character: &Character{
				Name: "one",
				Type: "mortal",
			},
			expectedVampire: &Vampire{
				Characters: []Character{
					Character{
						ID:   1,
						Name: "one",
						Type: "mortal",
					},
				},
			},
		},
		{
			name: "success with existing characters",
			vampire: &Vampire{
				Characters: []Character{
					Character{
						ID:   1,
						Name: "one",
						Type: "mortal",
					},
				},
			},
			character: &Character{
				Name: "two",
				Type: "immortal",
			},
			expectedVampire: &Vampire{
				Characters: []Character{
					Character{
						ID:   1,
						Name: "one",
						Type: "mortal",
					},
					Character{
						ID:   2,
						Name: "two",
						Type: "immortal",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.vampire.AddCharacter(tt.character)

			if diff := cmp.Diff(tt.expectedVampire, tt.vampire); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestAddMark(t *testing.T) {
	tests := []struct {
		name            string
		vampire         *Vampire
		mark            *Mark
		expectedVampire *Vampire
	}{
		{
			name:    "success with no characters",
			vampire: &Vampire{},
			mark: &Mark{
				Description: "one",
			},
			expectedVampire: &Vampire{
				Marks: []Mark{
					Mark{
						ID:          1,
						Description: "one",
					},
				},
			},
		},
		{
			name: "success with existing characters",
			vampire: &Vampire{
				Marks: []Mark{
					Mark{
						ID:          1,
						Description: "one",
					},
				},
			},
			mark: &Mark{
				Description: "two",
			},
			expectedVampire: &Vampire{
				Marks: []Mark{
					Mark{
						ID:          1,
						Description: "one",
					},
					Mark{
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

			tt.vampire.AddMark(tt.mark)

			if diff := cmp.Diff(tt.expectedVampire, tt.vampire); diff != "" {
				t.Error(diff)
			}
		})
	}
}
