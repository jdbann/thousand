package handlers

import (
	"context"
	"errors"
	"net/http"

	"emailaddress.horse/thousand/form"
	"emailaddress.horse/thousand/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type newUserRenderer interface {
	NewUser(http.ResponseWriter, *http.Request, *form.NewUserForm) error
}

func NewUser(r chi.Router, l *zap.Logger, t newUserRenderer) {
	r.Get("/user/new", func(w http.ResponseWriter, r *http.Request) {
		err := t.NewUser(w, r, form.NewUser("", ""))
		if err != nil {
			l.Error("failed to render", zap.Error(err))
			handleError(w, err)
		}
	})
}

type userCreator interface {
	CreateUser(context.Context, *form.NewUserForm) (models.User, error)
}

type sessionSetter interface {
	SetCurrentUserID(*http.Request, http.ResponseWriter, uuid.UUID) error
	SetFlash(*http.Request, http.ResponseWriter, string) error
}

func CreateUser(r chi.Router, l *zap.Logger, uc userCreator, t newUserRenderer, s sessionSetter) {
	r.Post("/user", func(w http.ResponseWriter, r *http.Request) {
		form := form.NewUser(
			r.FormValue("email"),
			r.FormValue("password"),
		)

		if !form.Valid() {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err := t.NewUser(w, r, form)
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
			err := t.NewUser(w, r, form)
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

		if err := s.SetCurrentUserID(r, w, user.ID); err != nil {
			l.Error("failed to set new user id in session", zap.Error(err))
			handleError(w, err)
			return
		}

		if err := s.SetFlash(r, w, "Thank you for signing up!"); err != nil {
			l.Error("failed to set flash", zap.Error(err))
			handleError(w, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
