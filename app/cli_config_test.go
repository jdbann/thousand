package app

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestApplyCLIConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		originalApp App
		cliConfig   CLIConfig
		expectedApp App
	}{
		{
			name: "applies set values",
			originalApp: App{
				DatabaseURL: "old_value",
			},
			cliConfig: CLIConfig{
				DatabaseURL: "new_value",
			},
			expectedApp: App{
				DatabaseURL: "new_value",
			},
		},
		{
			name: "does not apply zero values",
			originalApp: App{
				DatabaseURL: "old_value",
			},
			cliConfig: CLIConfig{},
			expectedApp: App{
				DatabaseURL: "old_value",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.cliConfig.applyTo(&tt.originalApp)

			if diff := cmp.Diff(tt.expectedApp, tt.originalApp); diff != "" {
				t.Error(diff)
			}
		})
	}
}
