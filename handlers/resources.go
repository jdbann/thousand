package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"emailaddress.horse/thousand/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type newResourceRenderer interface {
	NewResource(http.ResponseWriter, models.Vampire) error
}

func NewResource(r chi.Router, l *zap.Logger, t newResourceRenderer, vg vampireGetter) {
	r.Get("/vampires/{vampireID}/resources/new", func(w http.ResponseWriter, r *http.Request) {
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

		err = t.NewResource(w, vampire)
		if err != nil {
			l.Error("failed to render", zap.Error(err))
			handleError(w, err)
		}
	})
}

func CreateResource(r chi.Router, l *zap.Logger, rc resourceCreator) {
	r.Post("/vampires/{vampireID}/resources", func(w http.ResponseWriter, r *http.Request) {
		vampireID, err := uuid.Parse(chi.URLParam(r, "vampireID"))
		if err != nil {
			l.Error("failed to parse id as UUID", zap.Error(err))
			handleError(w, err)
			return
		}

		params := models.CreateResourceParams{
			Description: r.FormValue("description"),
		}

		if params.Stationary, err = strconv.ParseBool(r.FormValue("stationary")); err != nil {
			l.Error("failed to parse stationary param as bool", zap.Error(err))
			handleError(w, err)
			return
		}

		_, err = rc.CreateResource(r.Context(), vampireID, params)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				err = NotFoundError.Cause(err)
			}

			l.Error("failed to create resource", zap.Stringer("vampireID", vampireID), zap.Object("params", params), zap.Error(err))
			handleError(w, err)
			return
		}

		http.Redirect(w, r, "/vampires/"+vampireID.String(), http.StatusSeeOther)
	})
}
