package app

import (
	"emailaddress.horse/thousand/templates"
)

// Configurer takes an App as a param with the expectation of applying specific
// configuration to the app.
type Configurer func(*App)

// apply initially applies the BaseConfig to the provided App and then its own
// environment-specific configuration.
func (envConfig Configurer) apply(app *App) {
	baseConfig(app)
	envConfig(app)
}

// baseConfig sets up common configuration values for all environments.
func baseConfig(app *App) {
	app.Renderer = templates.NewRenderer()
}
