package models

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		memory         *Memory
		expectedResult bool
	}{
		{
			name: "false with no experiences",
			memory: &Memory{
				Experiences: []Experience{},
			},
			expectedResult: false,
		},
		{
			name: "false with less than three experiences",
			memory: &Memory{
				Experiences: []Experience{
					Experience("one"),
					Experience("two"),
				},
			},
			expectedResult: false,
		},
		{
			name: "true with three experiences",
			memory: &Memory{
				Experiences: []Experience{
					Experience("one"),
					Experience("two"),
					Experience("three"),
				},
			},
			expectedResult: true,
		},
		{
			name: "true with more than three experiences",
			memory: &Memory{
				Experiences: []Experience{
					Experience("one"),
					Experience("two"),
					Experience("three"),
					Experience("four"),
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

func TestMemory_AddExperience(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		memory         *Memory
		experience     Experience
		expectedMemory *Memory
		expectedError  error
	}{
		{
			name: "success with empty memory",
			memory: &Memory{
				Experiences: []Experience{},
			},
			experience: Experience("one"),
			expectedMemory: &Memory{
				Experiences: []Experience{
					Experience("one"),
				},
			},
			expectedError: nil,
		},
		{
			name: "success with partially filled memory",
			memory: &Memory{
				Experiences: []Experience{
					Experience("one"),
				},
			},
			experience: Experience("two"),
			expectedMemory: &Memory{
				Experiences: []Experience{
					Experience("one"),
					Experience("two"),
				},
			},
			expectedError: nil,
		},
		{
			name: "failure with full memory",
			memory: &Memory{
				Experiences: []Experience{
					Experience("one"),
					Experience("two"),
					Experience("three"),
				},
			},
			experience: Experience("four"),
			expectedMemory: &Memory{
				Experiences: []Experience{
					Experience("one"),
					Experience("two"),
					Experience("three"),
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
