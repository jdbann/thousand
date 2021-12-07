package app

import "github.com/labstack/echo/v4/middleware"

// DevelopmentConfig wraps the developmentConfig function as the appropriate
// EnvConfigurer type.
var DevelopmentConfig EnvConfigurer = developmentConfig

// developmentConfig sets up the app for a development environment.
func developmentConfig(app *App) {
	// Echo configuraton values
	app.Debug = true

	// Injected middleware
	app.LoggerMiddleware = middleware.Logger()
}
