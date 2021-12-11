package app

import (
	"net/http"
)

func MethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			method := r.FormValue("_method")

			if method != "" {
				r.Method = method
			}
		}

		next.ServeHTTP(w, r)
	})
}
