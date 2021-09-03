package app

import (
	"emailaddress.horse/thousand/app/models"
	"github.com/labstack/echo/v4"
)

// App is a configured instance of the application, ready to be served by a
// server or interacted with by CLI commands.
type App struct {
	*echo.Echo

	Logger echo.MiddlewareFunc

	Vampire *models.Vampire
}

// NewApp configures an instance of the application with helpful defaults.
func NewApp(configurer Configurer) *App {
	app := &App{
		Echo:    echo.New(),
		Vampire: &models.Vampire{},
	}

	configurer.apply(app)

	app.setupRoutes()

	return app
}
