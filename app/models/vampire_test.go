package models

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestForgetMemory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		vampire         *Vampire
		memoryID        int
		expectedVampire *Vampire
		expectedError   error
	}{
		{
			name: "success with recognised memoryID",
			vampire: &Vampire{
				Memories: []*OldMemory{
					{
						ID:          1,
						Experiences: []OldExperience{OldExperience("one")},
					},
				},
			},
			memoryID: 1,
			expectedVampire: &Vampire{
				Memories: []*OldMemory{
					{
						ID:          2,
						Experiences: []OldExperience{},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "success with recognised memoryID for non-incremental ID",
			vampire: &Vampire{
				Memories: []*OldMemory{
					{
						ID:          1,
						Experiences: []OldExperience{OldExperience("one")},
					},
					{
						ID:          2,
						Experiences: []OldExperience{OldExperience("two")},
					},
				},
			},
			memoryID: 1,
			expectedVampire: &Vampire{
				Memories: []*OldMemory{
					{
						ID:          3,
						Experiences: []OldExperience{},
					},
					{
						ID:          2,
						Experiences: []OldExperience{OldExperience("two")},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:            "failure with unrecognised memoryID",
			vampire:         &Vampire{},
			memoryID:        1,
			expectedVampire: &Vampire{},
			expectedError:   ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.vampire.ForgetMemory(tt.memoryID)

			if diff := cmp.Diff(tt.expectedVampire, tt.vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %s; actual %s", tt.expectedError, err)
			}
		})
	}
}

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
			name: "success with recognised memoryID",
			vampire: &Vampire{
				Memories: []*OldMemory{
					{
						ID: 1,
					},
				},
			},
			memoryID:         1,
			experienceString: "one",
			expectedVampire: &Vampire{
				Memories: []*OldMemory{
					{
						ID:          1,
						Experiences: []OldExperience{OldExperience("one")},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "failure with memoryID for full memory",
			vampire: &Vampire{
				Memories: []*OldMemory{
					{
						ID: 1,
						Experiences: []OldExperience{
							OldExperience("one"),
							OldExperience("two"),
							OldExperience("three"),
						},
					},
				},
			},
			memoryID:         1,
			experienceString: "four",
			expectedVampire: &Vampire{
				Memories: []*OldMemory{
					{
						ID: 1,
						Experiences: []OldExperience{
							OldExperience("one"),
							OldExperience("two"),
							OldExperience("three"),
						},
					},
				},
			},
			expectedError: ErrMemoryFull,
		},
		{
			name:             "failure with unrecognised memoryID",
			vampire:          &Vampire{},
			memoryID:         1,
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
				Skills: []*Skill{
					{
						ID:          1,
						Description: "one",
					},
				},
			},
		},
		{
			name: "success with existing skills",
			vampire: &Vampire{
				Skills: []*Skill{
					{
						ID:          1,
						Description: "one",
					},
				},
			},
			skill: &Skill{
				Description: "two",
			},
			expectedVampire: &Vampire{
				Skills: []*Skill{
					{
						ID:          1,
						Description: "one",
					},
					{
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
				Skills: []*Skill{
					{
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
				Skills: []*Skill{
					{
						ID:          1,
						Description: "one",
						Checked:     false,
					},
					{
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
				Skills: []*Skill{
					{
						ID:          1,
						Description: "one",
						Checked:     true,
					},
					{
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
				Resources: []*Resource{
					{
						ID:          1,
						Description: "one",
					},
				},
			},
		},
		{
			name: "success with existing resources",
			vampire: &Vampire{
				Resources: []*Resource{
					{
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
				Resources: []*Resource{
					{
						ID:          1,
						Description: "one",
					},
					{
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
				Resources: []*Resource{
					{
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
				Resources: []*Resource{
					{
						ID:          1,
						Description: "one",
						Lost:        false,
					},
					{
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
				Resources: []*Resource{
					{
						ID:          1,
						Description: "one",
						Lost:        true,
					},
					{
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
				Characters: []*Character{
					{
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
				Characters: []*Character{
					{
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
				Characters: []*Character{
					{
						ID:   1,
						Name: "one",
						Type: "mortal",
					},
					{
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

func TestFindCharacter(t *testing.T) {
	tests := []struct {
		name              string
		vampire           *Vampire
		characterID       int
		expectedCharacter *Character
		expectedError     error
	}{
		{
			name: "success",
			vampire: &Vampire{
				Characters: []*Character{
					{
						ID:   1,
						Name: "one",
						Type: "mortal",
					},
				},
			},
			characterID: 1,
			expectedCharacter: &Character{
				ID:   1,
				Name: "one",
				Type: "mortal",
			},
		},
		{
			name:          "failure with unknown ID",
			vampire:       &Vampire{},
			characterID:   1,
			expectedError: ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resource, err := tt.vampire.FindCharacter(tt.characterID)

			if diff := cmp.Diff(tt.expectedCharacter, resource); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %s; actual %s", tt.expectedError, err)
			}
		})
	}
}

func TestUpdateCharacter(t *testing.T) {
	tests := []struct {
		name            string
		vampire         *Vampire
		character       *Character
		expectedVampire *Vampire
		expectedError   error
	}{
		{
			name: "success",
			vampire: &Vampire{
				Characters: []*Character{
					{
						ID:       1,
						Name:     "one",
						Type:     "mortal",
						Deceased: true,
					},
					{
						ID:       2,
						Name:     "two",
						Type:     "mortal",
						Deceased: false,
					},
				},
			},
			character: &Character{
				ID:       2,
				Name:     "two",
				Type:     "mortal",
				Deceased: true,
			},
			expectedVampire: &Vampire{
				Characters: []*Character{
					{
						ID:       1,
						Name:     "one",
						Type:     "mortal",
						Deceased: true,
					},
					{
						ID:       2,
						Name:     "two",
						Type:     "mortal",
						Deceased: true,
					},
				},
			},
		},
		{
			name:    "failure with unknown resource",
			vampire: &Vampire{},
			character: &Character{
				ID:       1,
				Name:     "one",
				Type:     "mortal",
				Deceased: true,
			},
			expectedVampire: &Vampire{},
			expectedError:   ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.vampire.UpdateCharacter(tt.character)

			if diff := cmp.Diff(tt.expectedVampire, tt.vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %s; actual %s", tt.expectedError, err)
			}
		})
	}
}

func TestVampire_AddDescriptor(t *testing.T) {
	tests := []struct {
		name            string
		vampire         *Vampire
		characterID     int
		descriptor      string
		expectedVampire *Vampire
		expectedError   error
	}{
		{
			name: "success",
			vampire: &Vampire{
				Characters: []*Character{
					{
						ID:          1,
						Descriptors: []string{},
					},
				},
			},
			characterID: 1,
			descriptor:  "one",
			expectedVampire: &Vampire{
				Characters: []*Character{
					{
						ID:          1,
						Descriptors: []string{"one"},
					},
				},
			},
		},
		{
			name:            "failure with unknown resource",
			vampire:         &Vampire{},
			characterID:     1,
			descriptor:      "one",
			expectedVampire: &Vampire{},
			expectedError:   ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.vampire.AddDescriptor(tt.characterID, tt.descriptor)

			if diff := cmp.Diff(tt.expectedVampire, tt.vampire); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %s; actual %s", tt.expectedError, err)
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
				Marks: []*Mark{
					{
						ID:          1,
						Description: "one",
					},
				},
			},
		},
		{
			name: "success with existing characters",
			vampire: &Vampire{
				Marks: []*Mark{
					{
						ID:          1,
						Description: "one",
					},
				},
			},
			mark: &Mark{
				Description: "two",
			},
			expectedVampire: &Vampire{
				Marks: []*Mark{
					{
						ID:          1,
						Description: "one",
					},
					{
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
