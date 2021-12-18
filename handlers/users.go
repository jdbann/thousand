package handlers

import (
	"context"
	"errors"
	"net/http"

	"emailaddress.horse/thousand/form"
	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/session"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type newUserRenderer interface {
	NewUser(http.ResponseWriter, *form.NewUserForm) error
}

func NewUser(r chi.Router, l *zap.Logger, t newUserRenderer) {
	r.Get("/user/new", func(w http.ResponseWriter, r *http.Request) {
		err := t.NewUser(w, form.NewUser("", ""))
		if err != nil {
			l.Error("failed to render", zap.Error(err))
			handleError(w, err)
		}
	})
}

type userCreator interface {
	CreateUser(context.Context, *form.NewUserForm) (models.User, error)
}

func CreateUser(r chi.Router, l *zap.Logger, uc userCreator, t newUserRenderer, s *session.Store) {
	r.Post("/user", func(w http.ResponseWriter, r *http.Request) {
		form := form.NewUser(
			r.FormValue("email"),
			r.FormValue("password"),
		)

		if !form.Valid() {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err := t.NewUser(w, form)
			if err != nil {
				l.Error("failed to render", zap.Error(err))
				handleError(w, err)
			}
			return
		}

		user, err := uc.CreateUser(r.Context(), form)
		if errors.Is(err, models.ErrEmailAlreadyInUse) {
			form.Email.Message = "Email already in use."

			w.WriteHeader(http.StatusUnprocessableEntity)
			err := t.NewUser(w, form)
			if err != nil {
				l.Error("failed to render", zap.Error(err))
				handleError(w, err)
			}
			return
		} else if err != nil {
			l.Error("failed to create user", zap.Object("params", form), zap.Error(err))
			handleError(w, err)
			return
		}

		s.SetCurrentUserID(r, w, user.ID)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
