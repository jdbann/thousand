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
