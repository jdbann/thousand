package handlers

import (
	"context"
	"net/http"

	"emailaddress.horse/thousand/models"
	"emailaddress.horse/thousand/session"
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
				http.Redirect(w, r, "/user/new", http.StatusSeeOther)
				return
			}

			_, err := ug.GetUser(r.Context(), id)
			if err != nil {
				_ = s.ClearCurrentUserID(r, w)
				http.Redirect(w, r, "/user/new", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	})
}
