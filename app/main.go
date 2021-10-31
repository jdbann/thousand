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

	// Temporary data store
	Vampire *models.Vampire
}

// DBConnector is a function that returns a connection to the provided URL or
// any appropriate errors.
type DBConnector func(databaseURL string) (models.DBTX, error)

// NewApp configures an instance of the application with helpful defaults.
func NewApp(configurers ...Configurer) *App {
	app := &App{
		Echo: echo.New(),
		Vampire: &models.Vampire{
			Memories: []*models.OldMemory{
				{ID: 1},
				{ID: 2},
				{ID: 3},
				{ID: 4},
				{ID: 5},
			},
		},
	}

	for _, configurer := range configurers {
		configurer.applyTo(app)
	}

	app.setupRoutes()

	dbtx, err := app.DBConnector(app.DatabaseURL)
	if err != nil {
		app.Logger.Fatal(err)
	}

	app.Models = models.NewModels(dbtx)

	return app
}
