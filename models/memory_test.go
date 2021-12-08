package models

import (
	"testing"
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
