package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"emailaddress.horse/thousand/health"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type checker interface {
	Check(context.Context) (health.Result, bool)
}

func Health(r chi.Router, l *zap.Logger, h checker) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		result, ok := h.Check(r.Context())
		if !ok {
			w.WriteHeader(http.StatusBadGateway)
		}

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")

		err := enc.Encode(result)
		if err != nil {
			l.Error("error encoding health response", zap.Error(err))
			handleError(w, err)
		}
	})
}
