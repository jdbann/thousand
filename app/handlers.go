package app

import (
	"net/http"
	"strconv"

	"emailaddress.horse/thousand/app/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type templateData struct {
	*App
	Data dataMap
}

type dataMap map[string]interface{}

func (app *App) data(data dataMap) templateData {
	return templateData{app, data}
}

func (app *App) root(c echo.Context) error {
	return c.Render(http.StatusOK, "root", app)
}

func (app *App) showVampire(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	vampire, err := app.Models.GetVampire(c.Request().Context(), id)
	if err != nil {
		return err
	}

	data := app.data(dataMap{"vampire": vampire})
	return c.Render(http.StatusOK, "vampires/show", data)
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

func (app *App) oldCreateExperience(c echo.Context) error {
	memoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	if err := app.Vampire.AddExperience(memoryID, c.FormValue("experience")); err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) newExperience(c echo.Context) error {
	vampireID, err := uuid.Parse(c.Param("vampireID"))
	if err != nil {
		return err
	}

	memoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	memory, err := app.Models.GetMemory(c.Request().Context(), vampireID, memoryID)
	if err != nil {
		return err
	}

	data := app.data(dataMap{"memory": memory})
	return c.Render(http.StatusOK, "experiences/new", data)
}

func (app *App) createExperience(c echo.Context) error {
	vampireID, err := uuid.Parse(c.Param("vampireID"))
	if err != nil {
		return err
	}

	memoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	description := c.FormValue("description")

	if _, err := app.Models.AddExperience(c.Request().Context(), vampireID, memoryID, description); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, app.Reverse("show-vampire", vampireID.String()))
}

func (app *App) newSkill(c echo.Context) error {
	vampireID, err := uuid.Parse(c.Param("vampireID"))
	if err != nil {
		return err
	}

	vampire, err := app.Models.GetVampire(c.Request().Context(), vampireID)
	if err != nil {
		return err
	}

	data := app.data(dataMap{"vampire": vampire})
	return c.Render(http.StatusOK, "skills/new", data)
}

func (app *App) createSkill(c echo.Context) error {
	vampireID, err := uuid.Parse(c.Param("vampireID"))
	if err != nil {
		return err
	}

	description := c.FormValue("description")

	if _, err := app.Models.AddSkill(c.Request().Context(), vampireID, description); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, app.Reverse("show-vampire", vampireID.String()))
}

func (app *App) newResource(c echo.Context) error {
	vampireID, err := uuid.Parse(c.Param("vampireID"))
	if err != nil {
		return err
	}

	vampire, err := app.Models.GetVampire(c.Request().Context(), vampireID)
	if err != nil {
		return err
	}

	data := app.data(dataMap{"vampire": vampire})
	return c.Render(http.StatusOK, "resources/new", data)
}

func (app *App) createResource(c echo.Context) error {
	vampireID, err := uuid.Parse(c.Param("vampireID"))
	if err != nil {
		return err
	}

	var params models.AddResourceParams

	if err := c.Bind(&params); err != nil {
		return err
	}

	if _, err := app.Models.AddResource(c.Request().Context(), vampireID, params); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, app.Reverse("show-vampire", vampireID.String()))
}

func (app *App) newCharacter(c echo.Context) error {
	vampireID, err := uuid.Parse(c.Param("vampireID"))
	if err != nil {
		return err
	}

	vampire, err := app.Models.GetVampire(c.Request().Context(), vampireID)
	if err != nil {
		return err
	}

	data := app.data(dataMap{"vampire": vampire})
	return c.Render(http.StatusOK, "characters/new", data)
}

func (app *App) createCharacter(c echo.Context) error {
	vampireID, err := uuid.Parse(c.Param("vampireID"))
	if err != nil {
		return err
	}

	var params models.AddCharacterParams

	if err := c.Bind(&params); err != nil {
		return err
	}

	if _, err := app.Models.AddCharacter(c.Request().Context(), vampireID, params); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, app.Reverse("show-vampire", vampireID.String()))
}

func (app *App) oldCreateSkill(c echo.Context) error {
	skill := new(models.OldSkill)
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

func (app *App) oldCreateResource(c echo.Context) error {
	resource := new(models.OldResource)
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

func (app *App) oldCreateCharacter(c echo.Context) error {
	character := new(models.OldCharacter)
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

func (app *App) oldCreateMark(c echo.Context) error {
	mark := new(models.OldMark)
	if err := c.Bind(mark); err != nil {
		return err
	}

	app.Vampire.AddMark(mark)

	return c.Redirect(http.StatusFound, "/")
}

func (app *App) newMark(c echo.Context) error {
	vampireID, err := uuid.Parse(c.Param("vampireID"))
	if err != nil {
		return err
	}

	vampire, err := app.Models.GetVampire(c.Request().Context(), vampireID)
	if err != nil {
		return err
	}

	data := app.data(dataMap{"vampire": vampire})
	return c.Render(http.StatusOK, "marks/new", data)
}

func (app *App) createMark(c echo.Context) error {
	vampireID, err := uuid.Parse(c.Param("vampireID"))
	if err != nil {
		return err
	}

	description := c.FormValue("description")

	if _, err := app.Models.AddMark(c.Request().Context(), vampireID, description); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, app.Reverse("show-vampire", vampireID.String()))
}
