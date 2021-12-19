package handlers

import (
	"errors"
	"net/http"

	"emailaddress.horse/thousand/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type newSkillRenderer interface {
	NewSkill(http.ResponseWriter, *http.Request, models.Vampire) error
}

func NewSkill(r chi.Router, l *zap.Logger, t newSkillRenderer, vg vampireGetter) {
	r.Get("/vampires/{vampireID}/skills/new", func(w http.ResponseWriter, r *http.Request) {
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

		err = t.NewSkill(w, r, vampire)
		if err != nil {
			l.Error("failed to render", zap.Error(err))
			handleError(w, err)
		}
	})
}

func CreateSkill(r chi.Router, l *zap.Logger, sc skillCreator) {
	r.Post("/vampires/{vampireID}/skills", func(w http.ResponseWriter, r *http.Request) {
		vampireID, err := uuid.Parse(chi.URLParam(r, "vampireID"))
		if err != nil {
			l.Error("failed to parse id as UUID", zap.Error(err))
			handleError(w, err)
			return
		}

		description := r.FormValue("description")

		_, err = sc.CreateSkill(r.Context(), vampireID, description)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				err = NotFoundError.Cause(err)
			}

			l.Error("failed to create experience", zap.Stringer("vampireID", vampireID), zap.Error(err))
			handleError(w, err)
			return
		}

		http.Redirect(w, r, "/vampires/"+vampireID.String(), http.StatusSeeOther)
	})
}
