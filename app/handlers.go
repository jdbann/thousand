package app

import (
	"net/http"
	"strconv"

	"emailaddress.horse/thousand/app/models"
	"github.com/labstack/echo/v4"
)

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
