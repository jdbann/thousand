package app

import (
	"emailaddress.horse/thousand/templates"
)

// Configurer holds or generates configuration values which can be applied to an
// instance of App.
type Configurer interface {
	applyTo(*App)
}

// applyTo initially applies the BaseConfig to the provided App and then its own
// environment-specific configuration.
func (envConfig EnvConfigurer) applyTo(app *App) {
	baseConfig(app)
	envConfig(app)
}

// EnvConfigurer takes an App as a param and is expected to apply environment
// specific configuration to the app.
type EnvConfigurer func(*App)

// Ensure baseConfig satisfies the Configurer interface when typed as an
// EnvConfigurer.
var _ Configurer = EnvConfigurer(baseConfig)

// baseConfig sets up common configuration values for all environments.
func baseConfig(app *App) {
	app.Renderer = templates.NewRenderer()
}
