package handlers

import (
	"errors"
	"net/http"

	"emailaddress.horse/thousand/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func CreateCharacter(r *chi.Mux, l *zap.Logger, cc characterCreator) {
	r.Post("/vampires/{vampireID}/characters", func(w http.ResponseWriter, r *http.Request) {
		vampireID, err := uuid.Parse(chi.URLParam(r, "vampireID"))
		if err != nil {
			l.Error("failed to parse id as UUID", zap.Error(err))
			handleError(w, err)
			return
		}

		params := models.CreateCharacterParams{
			Name: r.FormValue("name"),
			Type: r.FormValue("type"),
		}

		_, err = cc.CreateCharacter(r.Context(), vampireID, params)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				err = NotFoundError.Cause(err)
			}

			l.Error("failed to create character", zap.Stringer("vampireID", vampireID), zap.Object("params", params), zap.Error(err))
			handleError(w, err)
			return
		}

		http.Redirect(w, r, "/vampires/"+vampireID.String(), http.StatusSeeOther)
	})
}
