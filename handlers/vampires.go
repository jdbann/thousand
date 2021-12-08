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

type newVampireRenderer interface {
	NewVampire(http.ResponseWriter) error
}

func NewVampire(r *chi.Mux, l *zap.Logger, t newVampireRenderer) {
	r.Get("/vampires/new", func(w http.ResponseWriter, r *http.Request) {
		err := t.NewVampire(w)
		if err != nil {
			l.Error("failed to render", zap.Error(err))
			handleError(w, err)
		}
	})
}

func CreateVampire(r *chi.Mux, l *zap.Logger, vc vampireCreator) {
	r.Post("/vampires", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")

		vampire, err := vc.CreateVampire(r.Context(), name)
		if err != nil {
			l.Error("failed to create vampire", zap.Error(err))
			handleError(w, err)
			return
		}

		http.Redirect(w, r, "/vampires/"+vampire.ID.String(), http.StatusSeeOther)
	})
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
