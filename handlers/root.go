package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Root(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/vampires", http.StatusSeeOther)
	})
}
