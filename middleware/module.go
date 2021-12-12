package middleware

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Invoke(register),
)

type RegisterParams struct {
	fx.In

	Logger *zap.Logger
	Router chi.Router
}

func register(s RegisterParams) {
	RequestLogger(s.Router, s.Logger.Named("server"))
	MethodOverride(s.Router)
	RedirectSlashes(s.Router)
}
