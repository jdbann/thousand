package app

import (
	"emailaddress.horse/thousand/app/models"
	"emailaddress.horse/thousand/templates"
	"github.com/labstack/echo/v4"
)

// App is a configured instance of the application, ready to be served by a
// server or interacted with by CLI commands.
type App struct {
	*echo.Echo
	Vampire *models.Vampire
}

// NewApp configures an instance of the application with helpful defaults.
func NewApp() *App {
	app := &App{
		Echo:    echo.New(),
		Vampire: &models.Vampire{},
	}

	app.configure()
	app.setupRoutes()

	return app
}

func (app *App) configure() {
	app.Debug = true
	app.Renderer = templates.NewRenderer()
}
