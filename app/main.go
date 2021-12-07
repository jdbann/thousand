package app

import (
	"emailaddress.horse/thousand/repository"
	"github.com/labstack/echo/v4"
)

// App is a configured instance of the application, ready to be served by a
// server or interacted with by CLI commands.
type App struct {
	*echo.Echo

	// Injected middleware
	LoggerMiddleware echo.MiddlewareFunc

	// Runtime values
	Repository *repository.Repository
}

type Options struct {
	Repository *repository.Repository
}

// NewApp configures an instance of the application with helpful defaults.
func NewApp(opts Options, configurers ...Configurer) *App {
	app := &App{Echo: echo.New()}

	for _, configurer := range configurers {
		configurer.applyTo(app)
	}

	app.Repository = opts.Repository

	app.setupRoutes()

	return app
}
