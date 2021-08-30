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

func (app *App) createSkill(c echo.Context) error {
	skill := new(models.Skill)
	if err := c.Bind(skill); err != nil {
		return err
	}

	app.Character.AddSkill(skill)

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) updateSkill(c echo.Context) error {
	skillID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	skill, err := app.Character.FindSkill(skillID)
	if err != nil {
		return err
	}

	if err := c.Bind(skill); err != nil {
		return err
	}

	if err := app.Character.UpdateSkill(skill); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) createResource(c echo.Context) error {
	resource := new(models.Resource)
	if err := c.Bind(resource); err != nil {
		return err
	}

	app.Character.AddResource(resource)

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) updateResource(c echo.Context) error {
	resourceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	resource, err := app.Character.FindResource(resourceID)
	if err != nil {
		return err
	}

	if err := c.Bind(resource); err != nil {
		return err
	}

	if err := app.Character.UpdateResource(resource); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}
