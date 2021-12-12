package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Invoke(register),
)

type RegisterParams struct {
	fx.In

	Logger   *zap.Logger
	Registry *prometheus.Registry
	Router   chi.Router
}

func register(p RegisterParams) {
	RequestLogger(p.Router, p.Logger.Named("server"))
	MethodOverride(p.Router)
	RedirectSlashes(p.Router)
	CollectMetrics(p.Router, p.Registry)
}
