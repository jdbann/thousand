package app

import (
	"context"
	"net/http"
	"strconv"

	"emailaddress.horse/thousand/app/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (app *App) root(c echo.Context) error {
	return c.Render(http.StatusOK, "root", app)
}

func (app *App) listVampires(c echo.Context) error {
	return c.Render(http.StatusOK, "vampires/index", app)
}

func (app *App) createVampire(c echo.Context) error {
	var vampire = new(models.NewVampire)
	if err := c.Bind(vampire); err != nil {
		return err
	}

	queries, err := app.Queries()
	if err != nil {
		return err
	}

	v, err := queries.CreateVampire(context.Background(), vampire.Name)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, app.Reverse("show-vampire", v.ID))
}

func (app *App) showVampire(c echo.Context) error {
	vampireID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	queries, err := app.Queries()
	if err != nil {
		return err
	}

	vampire, err := queries.GetVampire(context.Background(), vampireID)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "vampires/show", vampire)
}

func (app *App) createDetails(c echo.Context) error {
	var details = new(models.Details)
	if err := c.Bind(details); err != nil {
		return err
	}

	app.Vampire.Details = details

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) deleteMemory(c echo.Context) error {
	memoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	if err := app.Vampire.ForgetMemory(memoryID); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) createExperience(c echo.Context) error {
	memoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	if err := app.Vampire.AddExperience(memoryID, c.FormValue("experience")); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) createSkill(c echo.Context) error {
	skill := new(models.Skill)
	if err := c.Bind(skill); err != nil {
		return err
	}

	app.Vampire.AddSkill(skill)

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) updateSkill(c echo.Context) error {
	skillID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	skill, err := app.Vampire.FindSkill(skillID)
	if err != nil {
		return err
	}

	if err := c.Bind(skill); err != nil {
		return err
	}

	if err := app.Vampire.UpdateSkill(skill); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) createResource(c echo.Context) error {
	resource := new(models.Resource)
	if err := c.Bind(resource); err != nil {
		return err
	}

	app.Vampire.AddResource(resource)

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) updateResource(c echo.Context) error {
	resourceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	resource, err := app.Vampire.FindResource(resourceID)
	if err != nil {
		return err
	}

	if err := c.Bind(resource); err != nil {
		return err
	}

	if err := app.Vampire.UpdateResource(resource); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) createCharacter(c echo.Context) error {
	character := new(models.Character)
	if err := c.Bind(character); err != nil {
		return err
	}

	app.Vampire.AddCharacter(character)

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) updateCharacter(c echo.Context) error {
	characterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	character, err := app.Vampire.FindCharacter(characterID)
	if err != nil {
		return err
	}

	if err := c.Bind(character); err != nil {
		return err
	}

	if err := app.Vampire.UpdateCharacter(character); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) createDescriptor(c echo.Context) error {
	characterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	if err := app.Vampire.AddDescriptor(characterID, c.FormValue("descriptor")); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) createMark(c echo.Context) error {
	mark := new(models.Mark)
	if err := c.Bind(mark); err != nil {
		return err
	}

	app.Vampire.AddMark(mark)

	return c.Redirect(http.StatusFound, "/")
}
