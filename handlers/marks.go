package handlers

import (
	"errors"
	"net/http"

	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/templates"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func NewMark(e *echo.Echo, vg vampireGetter) {
	e.GET("/vampires/:vampireID/marks/new", func(c echo.Context) error {
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
		return c.Render(http.StatusOK, "marks/new", data)
	}).Name = "new-mark"
}

func CreateMark(e *echo.Echo, cm markCreator) {
	e.POST("/vampires/:vampireID/marks", func(c echo.Context) error {
		vampireID, err := uuid.Parse(c.Param("vampireID"))
		if err != nil {
			return err
		}

		description := c.FormValue("description")

		if _, err := cm.CreateMark(c.Request().Context(), vampireID, description); errors.Is(err, models.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Vampire could not be found").SetInternal(err)
		} else if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, e.Reverse("show-vampire", vampireID.String()))
	}).Name = "create-mark"
}
