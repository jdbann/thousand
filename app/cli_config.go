package app

// Ensure CLIConfig satisfies the Configurer interface.
var _ Configurer = CLIConfig{}

// CLIConfig is a temporary holder for CLI configuration values which can then
// be applied as a Configurer to the app.App configuration.
type CLIConfig struct {
	DatabaseURL string
}

func (c CLIConfig) applyTo(app *App) {
	app.DatabaseURL = c.DatabaseURL
}
