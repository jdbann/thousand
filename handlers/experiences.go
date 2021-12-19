package handlers

import (
	"errors"
	"net/http"

	"emailaddress.horse/thousand/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type newExperienceRenderer interface {
	NewExperience(http.ResponseWriter, *http.Request, models.Memory) error
}

func NewExperience(r chi.Router, l *zap.Logger, t newExperienceRenderer, mg memoryGetter) {
	r.Get("/vampires/{vampireID}/memories/{id}/experiences/new", func(w http.ResponseWriter, r *http.Request) {
		vampireID, err := uuid.Parse(chi.URLParam(r, "vampireID"))
		if err != nil {
			l.Error("failed to parse id as UUID", zap.Error(err))
			handleError(w, err)
			return
		}

		memoryID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			l.Error("failed to parse id as UUID", zap.Error(err))
			handleError(w, err)
			return
		}

		memory, err := mg.GetMemory(r.Context(), vampireID, memoryID)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				err = NotFoundError.Cause(err)
			}

			l.Error("failed to find memory", zap.Stringer("vampireID", vampireID), zap.Stringer("memoryID", memoryID), zap.Error(err))
			handleError(w, err)
			return
		}

		err = t.NewExperience(w, r, memory)
		if err != nil {
			l.Error("failed to render", zap.Error(err))
			handleError(w, err)
		}
	})
}

func CreateExperience(r chi.Router, l *zap.Logger, ec experienceCreator) {
	r.Post("/vampires/{vampireID}/memories/{id}/experiences", func(w http.ResponseWriter, r *http.Request) {
		vampireID, err := uuid.Parse(chi.URLParam(r, "vampireID"))
		if err != nil {
			l.Error("failed to parse id as UUID", zap.Error(err))
			handleError(w, err)
			return
		}

		memoryID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			l.Error("failed to parse id as UUID", zap.Error(err))
			handleError(w, err)
			return
		}

		description := r.FormValue("description")

		_, err = ec.CreateExperience(r.Context(), vampireID, memoryID, description)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				err = NotFoundError.Cause(err)
			}

			l.Error("failed to create experience", zap.Stringer("vampireID", vampireID), zap.Stringer("memoryID", memoryID), zap.Error(err))
			handleError(w, err)
			return
		}

		http.Redirect(w, r, "/vampires/"+vampireID.String(), http.StatusSeeOther)
	})
}
