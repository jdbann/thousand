package templates

import (
	"context"
	"net/http"

	"emailaddress.horse/thousand/models"
)

const (
	currentUserContextKey = "currentUser"
)

func RequestWithCurrentUser(r *http.Request, user models.User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), currentUserContextKey, user))
}
