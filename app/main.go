package app

import (
	"database/sql"

	"emailaddress.horse/thousand/app/models"
	"github.com/labstack/echo/v4"
)

// App is a configured instance of the application, ready to be served by a
// server or interacted with by CLI commands.
type App struct {
	*echo.Echo

	// Config values
	DatabaseURL string

	// Injected middleware
	LoggerMiddleware echo.MiddlewareFunc

	// Runtime values
	_db *sql.DB

	// Temporary data store
	Vampire *models.Vampire
}

// NewApp configures an instance of the application with helpful defaults.
func NewApp(configurers ...Configurer) *App {
	app := &App{
		Echo: echo.New(),
		Vampire: &models.Vampire{
			Memories: []*models.Memory{
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

	return app
}

// DB returns a live connection to the database. If a connection cannot be
// created, an error is returned.
func (app *App) DB() (*sql.DB, error) {
	if app._db == nil {
		_db, err := sql.Open("postgres", app.DatabaseURL)
		if err != nil {
			return nil, err
		}

		app._db = _db
	}

	return app._db, nil
}
