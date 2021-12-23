package handlers

import (
	"errors"
	"net/http"

	"emailaddress.horse/thousand/middleware"
	"emailaddress.horse/thousand/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type showVampiresRenderer interface {
	ShowVampires(http.ResponseWriter, *http.Request, []models.Vampire) error
}

func ListVampires(r chi.Router, l *zap.Logger, t showVampiresRenderer, vg vampiresGetter) {
	r.Get("/vampires", func(w http.ResponseWriter, r *http.Request) {
		vampires, err := vg.GetVampires(r.Context())
		if err != nil {
			l.Error("failed to load vampires", zap.Error(err))
			handleError(w, err)
			return
		}

		err = t.ShowVampires(w, r, vampires)
		if err != nil {
			l.Error("failed to render", zap.Error(err))
			handleError(w, err)
		}
	})
}

type newVampireRenderer interface {
	NewVampire(http.ResponseWriter, *http.Request) error
}

func NewVampire(r chi.Router, l *zap.Logger, t newVampireRenderer) {
	r.Get("/vampires/new", func(w http.ResponseWriter, r *http.Request) {
		err := t.NewVampire(w, r)
		if err != nil {
			l.Error("failed to render", zap.Error(err))
			handleError(w, err)
		}
	})
}

func CreateVampire(r chi.Router, l *zap.Logger, vc vampireCreator) {
	r.Post("/vampires", func(w http.ResponseWriter, r *http.Request) {
		user := middleware.CurrentUser(r.Context())
		name := r.FormValue("name")

		vampire, err := vc.CreateVampire(r.Context(), user.ID, name)
		if err != nil {
			l.Error("failed to create vampire", zap.Error(err))
			handleError(w, err)
			return
		}

		http.Redirect(w, r, "/vampires/"+vampire.ID.String(), http.StatusSeeOther)
	})
}

type showVampireRenderer interface {
	ShowVampire(http.ResponseWriter, *http.Request, models.Vampire) error
}

func ShowVampire(r chi.Router, l *zap.Logger, t showVampireRenderer, vg vampireGetter) {
	r.Get("/vampires/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			l.Error("failed to parse id as UUID", zap.Error(err))
			handleError(w, err)
			return
		}

		vampire, err := vg.GetVampire(r.Context(), id)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				err = NotFoundError.Cause(err)
			}

			l.Error("failed to find vampire", zap.Stringer("id", id), zap.Error(err))
			handleError(w, err)
			return
		}

		err = t.ShowVampire(w, r, vampire)
		if err != nil {
			l.Error("failed to render", zap.Error(err))
			handleError(w, err)
		}
	})
}
