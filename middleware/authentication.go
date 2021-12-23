package middleware

import (
	"context"
	"errors"
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
				http.Redirect(w, r, "/session/new", http.StatusSeeOther)
				return
			}

			user, err := ug.GetUser(r.Context(), id)
			if err != nil {
				_ = s.ClearCurrentUserID(w, r)
				http.Redirect(w, r, "/session/new", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, RequestWithCurrentUser(r, user))
		})
	})
}

type contextKey string

const (
	currentUserContextKey contextKey = "currentUser"
)

func RequestWithCurrentUser(r *http.Request, user models.User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), currentUserContextKey, user))
}

func MaybeCurrentUser(ctx context.Context) (models.User, bool) {
	user, ok := ctx.Value(currentUserContextKey).(models.User)
	return user, ok
}

func CurrentUser(ctx context.Context) models.User {
	user, ok := MaybeCurrentUser(ctx)
	if !ok {
		panic(errors.New("current user has not been set on context"))
	}

	return user
}
