package handlers

import (
	"errors"
	"net/http"

	"emailaddress.horse/thousand/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type newMarkRenderer interface {
	NewMark(http.ResponseWriter, *http.Request, models.Vampire) error
}

func NewMark(r chi.Router, l *zap.Logger, t newMarkRenderer, vg vampireGetter) {
	r.Get("/vampires/{vampireID}/marks/new", func(w http.ResponseWriter, r *http.Request) {
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

		err = t.NewMark(w, r, vampire)
		if err != nil {
			l.Error("failed to render", zap.Error(err))
			handleError(w, err)
		}
	})
}

func CreateMark(r chi.Router, l *zap.Logger, cm markCreator) {
	r.Post("/vampires/{vampireID}/marks", func(w http.ResponseWriter, r *http.Request) {
		vampireID, err := uuid.Parse(chi.URLParam(r, "vampireID"))
		if err != nil {
			l.Error("failed to parse id as UUID", zap.Error(err))
			handleError(w, err)
			return
		}

		description := r.FormValue("description")

		_, err = cm.CreateMark(r.Context(), vampireID, description)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				err = NotFoundError.Cause(err)
			}

			l.Error("failed to create mark", zap.Stringer("vampireID", vampireID), zap.String("description", description), zap.Error(err))
			handleError(w, err)
			return
		}

		http.Redirect(w, r, "/vampires/"+vampireID.String(), http.StatusSeeOther)
	})
}
