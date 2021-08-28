package app

import (
	"net/http"
	"strconv"

	"emailaddress.horse/thousand/app/models"
	"emailaddress.horse/thousand/static"
	"emailaddress.horse/thousand/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// App is a configured instance of the application, ready to be served by a
// server or interacted with by CLI commands.
type App struct {
	*echo.Echo
	Character *models.Character
}

// NewApp configures an instance of the application with helpful defaults.
func NewApp() *App {
	app := &App{
		Echo:      echo.New(),
		Character: &models.Character{},
	}

	app.configure()
	app.setupRoutes()

	return app
}

func (app *App) configure() {
	app.Debug = true
	app.Renderer = templates.NewRenderer()
}

func (app *App) setupRoutes() {
	app.Use(middleware.Logger())
	app.Use(static.Middleware())

	app.GET("/", app.root).Name = "root"

	app.POST("/details", app.createDetails).Name = "create-details"

	app.POST("/memories/:id/experiences", app.createExperience).Name = "create-experience"
}

func (app *App) root(c echo.Context) error {
	return c.Render(http.StatusOK, "index", app)
}

func (app *App) createDetails(c echo.Context) error {
	var details = new(models.Details)
	if err := c.Bind(details); err != nil {
		return err
	}

	app.Character.Details = details

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) createExperience(c echo.Context) error {
	memoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	if err := app.Character.AddExperience(memoryID, c.FormValue("experience")); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}
