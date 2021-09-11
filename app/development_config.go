package app

import "github.com/labstack/echo/v4/middleware"

// DevelopmentConfig sets up the app for a development environment.
var DevelopmentConfig Configurer = EnvConfigurer(developmentConfig)

func developmentConfig(app *App) {
	// Echo configuration values
	app.Debug = true

	// App configuration values
	app.DatabaseURL = "postgres://localhost:5432/thousand_development?sslmode=disable"

	// Injected middleware
	app.LoggerMiddleware = middleware.Logger()
}
