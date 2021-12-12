package repository

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(fxNew),
)

type Params struct {
	fx.In

	DatabaseURL string `name:"databaseURL"`
	Logger      *zap.Logger
}

func fxNew(params Params) (*Repository, error) {
	opts := Options{
		DatabaseURL: params.DatabaseURL,
		Logger:      params.Logger.Named("repository"),
	}
	return New(opts)
}
