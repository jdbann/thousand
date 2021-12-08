package handlers

import (
	"errors"
	"net/http"

	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/templates"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func NewSkill(e *echo.Echo, vg vampireGetter) {
	e.GET("/vampires/:vampireID/skills/new", func(c echo.Context) error {
		vampireID, err := uuid.Parse(c.Param("vampireID"))
		if err != nil {
			return err
		}

		vampire, err := vg.GetVampire(c.Request().Context(), vampireID)
		if errors.Is(err, models.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Vampire could not be found").SetInternal(err)
		} else if err != nil {
			return err
		}

		data := templates.NewData().Add("vampire", vampire)
		return c.Render(http.StatusOK, "skills/new", data)
	}).Name = "new-skill"
}

func CreateSkill(e *echo.Echo, sc skillCreator) {
	e.POST("/vampires/:vampireID/skills", func(c echo.Context) error {
		vampireID, err := uuid.Parse(c.Param("vampireID"))
		if err != nil {
			return err
		}

		description := c.FormValue("description")

		if _, err := sc.CreateSkill(c.Request().Context(), vampireID, description); errors.Is(err, models.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Vampire could not be found").SetInternal(err)
		} else if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, "/vampires/"+vampireID.String())
	}).Name = "create-skill"
}
