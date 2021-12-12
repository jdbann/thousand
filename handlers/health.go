package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type pinger interface {
	Ping(context.Context) error
}

func Health(r chi.Router, l *zap.Logger, p pinger, now func() time.Time) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response := healthResponse{
			Details: []healthState{
				{
					Name:      "repository",
					Status:    "ok",
					Timestamp: now(),
				},
			},
			Status: "ok",
		}

		if err := p.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			response.Status = "failed"
			response.Details[0].Status = "failed"
			response.Details[0].Error = err.Error()
		}

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")

		err := enc.Encode(response)
		if err != nil {
			l.Error("error encoding health response", zap.Error(err))
			handleError(w, err)
		}
	})
}

type healthResponse struct {
	Details []healthState `json:"details"`
	Status  string        `json:"status"`
}

type healthState struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
