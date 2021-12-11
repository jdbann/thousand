package handlers

import (
	"errors"
	"net/http"

	"emailaddress.horse/thousand/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type newCharacterRenderer interface {
	NewCharacter(http.ResponseWriter, models.Vampire) error
}

func NewCharacter(r *chi.Mux, l *zap.Logger, t newCharacterRenderer, vg vampireGetter) {
	r.Get("/vampires/{vampireID}/characters/new", func(w http.ResponseWriter, r *http.Request) {
		vampireID, err := uuid.Parse(chi.URLParam(r, "vampireID"))
		if err != nil {
			l.Error("failed to parse id as UUID", zap.Error(err))
			handleError(w, err)
			return
		}

		vampire, err := vg.GetVampire(r.Context(), vampireID)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				err = NotFoundError.Cause(err)
			}

			l.Error("failed to find vampire", zap.Stringer("vampireID", vampireID), zap.Error(err))
			handleError(w, err)
			return
		}

		err = t.NewCharacter(w, vampire)
		if err != nil {
			l.Error("failed to render", zap.Error(err))
			handleError(w, err)
		}
	})
}

func CreateCharacter(e *echo.Echo, cc characterCreator) {
	e.POST("/vampires/:vampireID/characters", func(c echo.Context) error {
		vampireID, err := uuid.Parse(c.Param("vampireID"))
		if err != nil {
			return err
		}

		var params models.CreateCharacterParams

		if err := c.Bind(&params); err != nil {
			return err
		}

		if _, err := cc.CreateCharacter(c.Request().Context(), vampireID, params); errors.Is(err, models.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Vampire could not be found").SetInternal(err)
		} else if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, "/vampires/"+vampireID.String())
	}).Name = "create-character"
}
