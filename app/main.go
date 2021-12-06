package app

import (
	"emailaddress.horse/thousand/app/models"
	"github.com/labstack/echo/v4"
)

// App is a configured instance of the application, ready to be served by a
// server or interacted with by CLI commands.
type App struct {
	*echo.Echo

	// Config values
	DatabaseURL string
	DBConnector DBConnector

	// Injected middleware
	LoggerMiddleware echo.MiddlewareFunc

	// Runtime values
	Models *models.Models
}

// DBConnector is a function that returns a connection to the provided URL or
// any appropriate errors.
type DBConnector func(databaseURL string) (models.DBTX, error)

// NewApp configures an instance of the application with helpful defaults.
func NewApp(configurers ...Configurer) *App {
	app := &App{Echo: echo.New()}

	for _, configurer := range configurers {
		configurer.applyTo(app)
	}

	dbtx, err := app.DBConnector(app.DatabaseURL)
	if err != nil {
		app.Logger.Fatal(err)
	}

	app.Models = models.NewModels(dbtx)

	app.setupRoutes()

	return app
}
