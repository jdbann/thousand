package handlers

import (
	"errors"
	"net/http"

	"emailaddress.horse/thousand/app/models"
	"emailaddress.horse/thousand/templates"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ListVampires(e *echo.Echo, vg vampiresGetter) {
	e.GET("/vampires", func(c echo.Context) error {
		vampires, err := vg.GetVampires(c.Request().Context())
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "vampires/index", templates.NewData().Add("vampires", vampires))
	}).Name = "list-vampires"
}

func NewVampire(e *echo.Echo) {
	e.GET("/vampires/new", func(c echo.Context) error {
		return c.Render(http.StatusOK, "vampires/new", templates.NewData())
	}).Name = "new-vampire"
}

func CreateVampire(e *echo.Echo, vc vampireCreator) {
	e.POST("/vampires", func(c echo.Context) error {
		name := c.FormValue("name")

		vampire, err := vc.CreateVampire(c.Request().Context(), name)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, e.Reverse("show-vampire", vampire.ID.String()))
	}).Name = "create-vampire"
}

func ShowVampire(e *echo.Echo, vg vampireGetter) {
	e.GET("/vampires/:id", func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return err
		}

		vampire, err := vg.GetVampire(c.Request().Context(), id)
		if errors.Is(err, models.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Vampire could not be found").SetInternal(err)
		} else if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "vampires/show", templates.NewData().Add("vampire", vampire))
	}).Name = "show-vampire"
}
