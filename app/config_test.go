package app

import (
	"errors"
	"fmt"
	"testing"
)

func TestConfigFor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		environmentName    string
		expectedConfigurer Configurer
		expectedError      error
	}{
		{
			name:               "development",
			environmentName:    "development",
			expectedConfigurer: DevelopmentConfig,
			expectedError:      nil,
		},
		{
			name:               "test",
			environmentName:    "test",
			expectedConfigurer: LiveTestConfig,
			expectedError:      nil,
		},
		{
			name:               "unknown environment",
			environmentName:    "scary",
			expectedConfigurer: nil,
			expectedError:      ErrUnrecognisedEnvironment,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualConfigurer, err := ConfigFor(tt.environmentName)

			if fmt.Sprintf("%v", tt.expectedConfigurer) != fmt.Sprintf("%v", actualConfigurer) {
				t.Error("Configurers do not match")
			}

			if !errors.Is(tt.expectedError, err) {
				t.Errorf("expected error %q; got error %q", tt.expectedError, err)
			}
		})
	}
}
