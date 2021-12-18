// TODO: Move this into the handlers package
package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RedirectSlashes(r chi.Router) {
	r.Use(middleware.RedirectSlashes)
}

func MethodOverride(r chi.Router) {
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				method := r.FormValue("_method")

				if method != "" {
					r.Method = method
				}
			}

			next.ServeHTTP(w, r)
		})
	})
}
