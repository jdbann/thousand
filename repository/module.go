package repository

import (
	"emailaddress.horse/thousand/health"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(fxNew),
)

type Params struct {
	fx.In

	DatabaseURL string `name:"databaseURL"`

	Health *health.Health
	Logger *zap.Logger
}

func fxNew(params Params) (*Repository, error) {
	opts := Options{
		DatabaseURL: params.DatabaseURL,
		Logger:      params.Logger.Named("repository"),
	}

	repo, err := New(opts)
	if err != nil {
		return nil, err
	}

	params.Health.Register(&health.Component{
		Name:  "repository",
		Check: repo.Ping,
	})

	return repo, nil
}
