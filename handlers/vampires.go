package handlers

import (
	"context"
	"net/http"

	"emailaddress.horse/thousand/app/models"
	"emailaddress.horse/thousand/templates"
	"github.com/labstack/echo/v4"
)

type vampiresGetter interface {
	GetVampires(context.Context) ([]models.Vampire, error)
}

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

type vampireCreator interface {
	CreateVampire(context.Context, string) (models.Vampire, error)
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

func ShowVampire(e *echo.Echo) {
	e.GET("/vampires/:id", func(c echo.Context) error {
		return nil
	}).Name = "show-vampire"
}
