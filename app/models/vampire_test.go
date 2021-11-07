package models

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestForgetMemory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		vampire         *OldVampire
		memoryID        int
		expectedVampire *OldVampire
		expectedError   error
	}{
		{
			name: "success with recognised memoryID",
			vampire: &OldVampire{
				Memories: []*OldMemory{
					{
						ID:          1,
						Experiences: []OldExperience{OldExperience("one")},
					},
				},
			},
			memoryID: 1,
			expectedVampire: &OldVampire{
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
			vampire: &OldVampire{
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
			expectedVampire: &OldVampire{
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
			vampire:         &OldVampire{},
			memoryID:        1,
			expectedVampire: &OldVampire{},
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
		vampire          *OldVampire
		memoryID         int
		experienceString string
		expectedVampire  *OldVampire
		expectedError    error
	}{
		{
			name: "success with recognised memoryID",
			vampire: &OldVampire{
				Memories: []*OldMemory{
					{
						ID: 1,
					},
				},
			},
			memoryID:         1,
			experienceString: "one",
			expectedVampire: &OldVampire{
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
			vampire: &OldVampire{
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
			expectedVampire: &OldVampire{
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
			vampire:          &OldVampire{},
			memoryID:         1,
			experienceString: "one",
			expectedVampire:  &OldVampire{},
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
		vampire         *OldVampire
		skill           *OldSkill
		expectedVampire *OldVampire
	}{
		{
			name:    "success with no skills",
			vampire: &OldVampire{},
			skill: &OldSkill{
				Description: "one",
			},
			expectedVampire: &OldVampire{
				Skills: []*OldSkill{
					{
						ID:          1,
						Description: "one",
					},
				},
			},
		},
		{
			name: "success with existing skills",
			vampire: &OldVampire{
				Skills: []*OldSkill{
					{
						ID:          1,
						Description: "one",
					},
				},
			},
			skill: &OldSkill{
				Description: "two",
			},
			expectedVampire: &OldVampire{
				Skills: []*OldSkill{
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
		vampire       *OldVampire
		skillID       int
		expectedSkill *OldSkill
		expectedError error
	}{
		{
			name: "success",
			vampire: &OldVampire{
				Skills: []*OldSkill{
					{
						ID:          1,
						Description: "one",
					},
				},
			},
			skillID: 1,
			expectedSkill: &OldSkill{
				ID:          1,
				Description: "one",
			},
		},
		{
			name:          "failure with unknown ID",
			vampire:       &OldVampire{},
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
		vampire         *OldVampire
		skill           *OldSkill
		expectedVampire *OldVampire
		expectedError   error
	}{
		{
			name: "success",
			vampire: &OldVampire{
				Skills: []*OldSkill{
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
			skill: &OldSkill{
				ID:          1,
				Description: "one",
				Checked:     true,
			},
			expectedVampire: &OldVampire{
				Skills: []*OldSkill{
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
			vampire: &OldVampire{},
			skill: &OldSkill{
				ID:          1,
				Description: "one",
				Checked:     true,
			},
			expectedVampire: &OldVampire{},
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
		vampire         *OldVampire
		resource        *OldResource
		expectedVampire *OldVampire
	}{
		{
			name:    "success with no resources",
			vampire: &OldVampire{},
			resource: &OldResource{
				Description: "one",
			},
			expectedVampire: &OldVampire{
				Resources: []*OldResource{
					{
						ID:          1,
						Description: "one",
					},
				},
			},
		},
		{
			name: "success with existing resources",
			vampire: &OldVampire{
				Resources: []*OldResource{
					{
						ID:          1,
						Description: "one",
					},
				},
			},
			resource: &OldResource{
				Description: "two",
				Stationary:  true,
			},
			expectedVampire: &OldVampire{
				Resources: []*OldResource{
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
		vampire          *OldVampire
		resourceID       int
		expectedResource *OldResource
		expectedError    error
	}{
		{
			name: "success",
			vampire: &OldVampire{
				Resources: []*OldResource{
					{
						ID:          1,
						Description: "one",
					},
				},
			},
			resourceID: 1,
			expectedResource: &OldResource{
				ID:          1,
				Description: "one",
			},
		},
		{
			name:          "failure with unknown ID",
			vampire:       &OldVampire{},
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
		vampire         *OldVampire
		resource        *OldResource
		expectedVampire *OldVampire
		expectedError   error
	}{
		{
			name: "success",
			vampire: &OldVampire{
				Resources: []*OldResource{
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
			resource: &OldResource{
				ID:          1,
				Description: "one",
				Lost:        true,
			},
			expectedVampire: &OldVampire{
				Resources: []*OldResource{
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
			vampire: &OldVampire{},
			resource: &OldResource{
				ID:          1,
				Description: "one",
				Lost:        true,
			},
			expectedVampire: &OldVampire{},
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
		vampire         *OldVampire
		character       *OldCharacter
		expectedVampire *OldVampire
	}{
		{
			name:    "success with no characters",
			vampire: &OldVampire{},
			character: &OldCharacter{
				Name: "one",
				Type: "mortal",
			},
			expectedVampire: &OldVampire{
				Characters: []*OldCharacter{
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
			vampire: &OldVampire{
				Characters: []*OldCharacter{
					{
						ID:   1,
						Name: "one",
						Type: "mortal",
					},
				},
			},
			character: &OldCharacter{
				Name: "two",
				Type: "immortal",
			},
			expectedVampire: &OldVampire{
				Characters: []*OldCharacter{
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
		vampire           *OldVampire
		characterID       int
		expectedCharacter *OldCharacter
		expectedError     error
	}{
		{
			name: "success",
			vampire: &OldVampire{
				Characters: []*OldCharacter{
					{
						ID:   1,
						Name: "one",
						Type: "mortal",
					},
				},
			},
			characterID: 1,
			expectedCharacter: &OldCharacter{
				ID:   1,
				Name: "one",
				Type: "mortal",
			},
		},
		{
			name:          "failure with unknown ID",
			vampire:       &OldVampire{},
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
		vampire         *OldVampire
		character       *OldCharacter
		expectedVampire *OldVampire
		expectedError   error
	}{
		{
			name: "success",
			vampire: &OldVampire{
				Characters: []*OldCharacter{
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
			character: &OldCharacter{
				ID:       2,
				Name:     "two",
				Type:     "mortal",
				Deceased: true,
			},
			expectedVampire: &OldVampire{
				Characters: []*OldCharacter{
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
			vampire: &OldVampire{},
			character: &OldCharacter{
				ID:       1,
				Name:     "one",
				Type:     "mortal",
				Deceased: true,
			},
			expectedVampire: &OldVampire{},
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
		vampire         *OldVampire
		characterID     int
		descriptor      string
		expectedVampire *OldVampire
		expectedError   error
	}{
		{
			name: "success",
			vampire: &OldVampire{
				Characters: []*OldCharacter{
					{
						ID:          1,
						Descriptors: []string{},
					},
				},
			},
			characterID: 1,
			descriptor:  "one",
			expectedVampire: &OldVampire{
				Characters: []*OldCharacter{
					{
						ID:          1,
						Descriptors: []string{"one"},
					},
				},
			},
		},
		{
			name:            "failure with unknown resource",
			vampire:         &OldVampire{},
			characterID:     1,
			descriptor:      "one",
			expectedVampire: &OldVampire{},
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
		vampire         *OldVampire
		mark            *OldMark
		expectedVampire *OldVampire
	}{
		{
			name:    "success with no characters",
			vampire: &OldVampire{},
			mark: &OldMark{
				Description: "one",
			},
			expectedVampire: &OldVampire{
				Marks: []*OldMark{
					{
						ID:          1,
						Description: "one",
					},
				},
			},
		},
		{
			name: "success with existing characters",
			vampire: &OldVampire{
				Marks: []*OldMark{
					{
						ID:          1,
						Description: "one",
					},
				},
			},
			mark: &OldMark{
				Description: "two",
			},
			expectedVampire: &OldVampire{
				Marks: []*OldMark{
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
