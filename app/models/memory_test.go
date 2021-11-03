package models

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		memory         Memory
		expectedResult bool
	}{
		{
			name: "false with no experiences",
			memory: Memory{
				Experiences: []Experience{},
			},
			expectedResult: false,
		},
		{
			name: "false with less than three experiences",
			memory: Memory{
				Experiences: []Experience{
					{},
					{},
				},
			},
			expectedResult: false,
		},
		{
			name: "true with three experiences",
			memory: Memory{
				Experiences: []Experience{
					{},
					{},
					{},
				},
			},
			expectedResult: true,
		},
		{
			name: "true with more than three experiences",
			memory: Memory{
				Experiences: []Experience{
					{},
					{},
					{},
					{},
				},
			},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualResult := tt.memory.Full()

			if tt.expectedResult != actualResult {
				t.Errorf("expected %t; actual %t", tt.expectedResult, actualResult)
			}
		})
	}
}

func TestOldMemory_Full(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		memory         *OldMemory
		expectedResult bool
	}{
		{
			name: "false with no experiences",
			memory: &OldMemory{
				Experiences: []OldExperience{},
			},
			expectedResult: false,
		},
		{
			name: "false with less than three experiences",
			memory: &OldMemory{
				Experiences: []OldExperience{
					OldExperience("one"),
					OldExperience("two"),
				},
			},
			expectedResult: false,
		},
		{
			name: "true with three experiences",
			memory: &OldMemory{
				Experiences: []OldExperience{
					OldExperience("one"),
					OldExperience("two"),
					OldExperience("three"),
				},
			},
			expectedResult: true,
		},
		{
			name: "true with more than three experiences",
			memory: &OldMemory{
				Experiences: []OldExperience{
					OldExperience("one"),
					OldExperience("two"),
					OldExperience("three"),
					OldExperience("four"),
				},
			},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualResult := tt.memory.Full()

			if tt.expectedResult != actualResult {
				t.Errorf("expected %t; actual %t", tt.expectedResult, actualResult)
			}
		})
	}
}

func TestOldMemory_AddExperience(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		memory         *OldMemory
		experience     OldExperience
		expectedMemory *OldMemory
		expectedError  error
	}{
		{
			name: "success with empty memory",
			memory: &OldMemory{
				Experiences: []OldExperience{},
			},
			experience: OldExperience("one"),
			expectedMemory: &OldMemory{
				Experiences: []OldExperience{
					OldExperience("one"),
				},
			},
			expectedError: nil,
		},
		{
			name: "success with partially filled memory",
			memory: &OldMemory{
				Experiences: []OldExperience{
					OldExperience("one"),
				},
			},
			experience: OldExperience("two"),
			expectedMemory: &OldMemory{
				Experiences: []OldExperience{
					OldExperience("one"),
					OldExperience("two"),
				},
			},
			expectedError: nil,
		},
		{
			name: "failure with full memory",
			memory: &OldMemory{
				Experiences: []OldExperience{
					OldExperience("one"),
					OldExperience("two"),
					OldExperience("three"),
				},
			},
			experience: OldExperience("four"),
			expectedMemory: &OldMemory{
				Experiences: []OldExperience{
					OldExperience("one"),
					OldExperience("two"),
					OldExperience("three"),
				},
			},
			expectedError: ErrMemoryFull,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.memory.AddExperience(tt.experience)

			if diff := cmp.Diff(tt.expectedMemory, tt.memory); diff != "" {
				t.Error(diff)
			}

			if tt.expectedError != err {
				t.Errorf("expected %s; actual %s", tt.expectedError, err)
			}
		})
	}
}
