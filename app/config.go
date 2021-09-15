package app

import (
	"errors"

	"emailaddress.horse/thousand/templates"
)

var (
	// ErrUnrecognisedEnvironment is returned when trying to retrieve the config
	// for an environment which has not been defined.
	ErrUnrecognisedEnvironment = errors.New("unrecognised environment")

	_environmentMap = map[string]Configurer{
		"development": DevelopmentConfig,
		"test":        LiveTestConfig,
	}
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

// ConfigFor returns the correct Configurer for the requested environment, or an
// error indicating that the requested environment cannot be found.
func ConfigFor(environment string) (Configurer, error) {
	configurer, ok := _environmentMap[environment]
	if ok {
		return configurer, nil
	}

	return nil, ErrUnrecognisedEnvironment
}
