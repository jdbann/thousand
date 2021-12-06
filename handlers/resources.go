package handlers

import (
	"context"
	"errors"
	"net/http"

	"emailaddress.horse/thousand/app/models"
	"emailaddress.horse/thousand/templates"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func NewResource(e *echo.Echo, vg vampireGetter) {
	e.GET("/vampires/:vampireID/resources/new", func(c echo.Context) error {
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
		return c.Render(http.StatusOK, "resources/new", data)
	}).Name = "new-resource"
}

type resourceCreator interface {
	CreateResource(context.Context, uuid.UUID, models.CreateResourceParams) (models.Resource, error)
}

func CreateResource(e *echo.Echo, rc resourceCreator) {
	e.POST("/vampires/:vampireID/resources", func(c echo.Context) error {
		vampireID, err := uuid.Parse(c.Param("vampireID"))
		if err != nil {
			return err
		}

		var params models.CreateResourceParams

		if err := c.Bind(&params); err != nil {
			return err
		}

		if _, err := rc.CreateResource(c.Request().Context(), vampireID, params); errors.Is(err, models.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Vampire could not be found").SetInternal(err)
		} else if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, e.Reverse("show-vampire", vampireID.String()))
	}).Name = "create-resource"
}
