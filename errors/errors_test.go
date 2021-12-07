package errors_test

import (
	e "errors"
	"fmt"
	"testing"

	"emailaddress.horse/thousand/errors"
)

func TestIs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		err            error
		targetErr      error
		expectedResult bool
	}{
		{
			name:           "true with same name",
			err:            errors.New("A message"),
			targetErr:      errors.New("A message"),
			expectedResult: true,
		},
		{
			name:           "true with different causes",
			err:            errors.New("A message").Cause(fmt.Errorf("something")),
			targetErr:      errors.New("A message").Cause(fmt.Errorf("something else")),
			expectedResult: true,
		},
		{
			name:           "false with different names",
			err:            errors.New("A message"),
			targetErr:      errors.New("A different message"),
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := e.Is(tt.err, tt.targetErr)

			if result != tt.expectedResult {
				t.Errorf("expected %t; received %t", tt.expectedResult, result)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	t.Parallel()

	var errExample = errors.New("An example")

	tests := []struct {
		name           string
		err            error
		expectedResult error
	}{
		{
			name:           "cause if cause set",
			err:            errors.New("A message").Cause(errExample),
			expectedResult: errExample,
		},
		{
			name:           "nil if cause not set",
			err:            errors.New("A message"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := e.Unwrap(tt.err)

			if result != tt.expectedResult {
				t.Errorf("expected %q; received %q", tt.expectedResult, result)
			}
		})
	}
}
