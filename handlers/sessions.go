package handlers

import (
	"context"
	"net/http"

	"emailaddress.horse/thousand/form"
	"emailaddress.horse/thousand/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type newSessionRenderer interface {
	NewSession(http.ResponseWriter, *http.Request, *form.NewSessionForm) error
}

func NewSession(r chi.Router, l *zap.Logger, t newSessionRenderer) {
	r.Get("/session/new", func(w http.ResponseWriter, r *http.Request) {
		err := t.NewSession(w, r, form.NewSession("", ""))
		if err != nil {
			l.Error("failed to render", zap.Error(err))
			handleError(w, err)
		}
	})
}

type userAuthenticator interface {
	AuthenticateUser(context.Context, *form.NewSessionForm) (models.User, error)
}

func CreateSession(r chi.Router, l *zap.Logger, ua userAuthenticator, t newSessionRenderer, s sessionSetter) {
	r.Post("/session", func(w http.ResponseWriter, r *http.Request) {
		form := form.NewSession(
			r.FormValue("email"),
			r.FormValue("password"),
		)

		if !form.Valid() {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err := t.NewSession(w, r, form)
			if err != nil {
				l.Error("failed to render", zap.Error(err))
				handleError(w, err)
			}
			return
		}

		user, err := ua.AuthenticateUser(r.Context(), form)
		if err != nil {
			l.Error("failed to authenticate user", zap.Error(err))
			handleError(w, err)
			return
		}

		if user.ID == (uuid.UUID{}) {
			form.Email.Message = "No user found with this email address and password."

			w.WriteHeader(http.StatusUnprocessableEntity)
			err := t.NewSession(w, r, form)
			if err != nil {
				l.Error("failed to render", zap.Error(err))
				handleError(w, err)
			}
			return
		}

		if err := s.SetCurrentUserID(r, w, user.ID); err != nil {
			l.Error("failed to set user id in session", zap.Error(err))
			handleError(w, err)
			return
		}

		if err := s.SetFlash(r, w, "Welcome back!"); err != nil {
			l.Error("failed to set flash", zap.Error(err))
			handleError(w, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

type currentUserIDClearer interface {
	ClearCurrentUserID(http.ResponseWriter, *http.Request) error
}

func DestroySession(r chi.Router, l *zap.Logger, s currentUserIDClearer) {
	r.Delete("/session", func(w http.ResponseWriter, r *http.Request) {
		if err := s.ClearCurrentUserID(w, r); err != nil {
			l.Error("failed to clear user id in session", zap.Error(err))
			handleError(w, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
