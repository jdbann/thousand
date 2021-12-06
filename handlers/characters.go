package handlers

import (
	"errors"
	"net/http"

	"emailaddress.horse/thousand/app/models"
	"emailaddress.horse/thousand/templates"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func NewCharacter(e *echo.Echo, vg vampireGetter) {
	e.GET("/vampires/:vampireID/characters/new", func(c echo.Context) error {
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
		return c.Render(http.StatusOK, "characters/new", data)
	}).Name = "new-character"
}
