package handlers

import (
	"errors"
	"net/http"

	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/templates"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func NewExperience(e *echo.Echo, mg memoryGetter) {
	e.GET("/vampires/:vampireID/memories/:id/experiences/new", func(c echo.Context) error {
		vampireID, err := uuid.Parse(c.Param("vampireID"))
		if err != nil {
			return err
		}

		memoryID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return err
		}

		memory, err := mg.GetMemory(c.Request().Context(), vampireID, memoryID)
		if errors.Is(err, models.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Memory could not be found").SetInternal(err)
		} else if err != nil {
			return err
		}

		data := templates.NewData().Add("memory", memory)
		return c.Render(http.StatusOK, "experiences/new", data)
	}).Name = "new-experience"
}

func CreateExperience(e *echo.Echo, ec experienceCreator) {
	e.POST("/vampires/:vampireID/memories/:id/experiences", func(c echo.Context) error {
		vampireID, err := uuid.Parse(c.Param("vampireID"))
		if err != nil {
			return err
		}

		memoryID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return err
		}

		description := c.FormValue("description")

		_, err = ec.CreateExperience(c.Request().Context(), vampireID, memoryID, description)
		if errors.Is(err, models.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Memory could not be found").SetInternal(err)
		} else if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, e.Reverse("show-vampire", vampireID.String()))
	}).Name = "create-experience"
}
