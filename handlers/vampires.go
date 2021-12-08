package handlers

import (
	"errors"
	"net/http"

	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/templates"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type showVampiresRenderer interface {
	ShowVampires(http.ResponseWriter, []models.Vampire) error
}

func ListVampires(r *chi.Mux, l *zap.Logger, t showVampiresRenderer, vg vampiresGetter) {
	r.Get("/vampires", func(w http.ResponseWriter, r *http.Request) {
		vampires, err := vg.GetVampires(r.Context())
		if err != nil {
			l.Error("failed to load vampires", zap.Error(err))
			handleError(w, err)
			return
		}

		err = t.ShowVampires(w, vampires)
		if err != nil {
			l.Error("failed to render", zap.Error(err))
			handleError(w, err)
		}
	})
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
