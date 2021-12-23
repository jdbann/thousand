package handlers

import (
	"context"
	"net/http"

	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/session"
	"emailaddress.horse/thousand/templates"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type userGetter interface {
	GetUser(context.Context, uuid.UUID) (models.User, error)
}

func EnsureLoggedIn(r chi.Router, s *session.Store, ug userGetter) {
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, ok := s.GetCurrentUserID(r)
			if !ok {
				http.Redirect(w, r, "/session/new", http.StatusSeeOther)
				return
			}

			user, err := ug.GetUser(r.Context(), id)
			if err != nil {
				_ = s.ClearCurrentUserID(w, r)
				http.Redirect(w, r, "/session/new", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, templates.RequestWithCurrentUser(r, user))
		})
	})
}
