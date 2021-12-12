package logger

import (
	"context"
	"errors"
	"syscall"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(fxNew),
	fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger.Named("fx")}
	}),
)

type Params struct {
	fx.In

	LogFormat string `name:"logFormat" optional:"true"`
}

func fxNew(lc fx.Lifecycle, params Params) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	switch params.LogFormat {
	case "prod":
		logger, err = zap.NewProduction()
	case "dev":
		logger, err = zap.NewDevelopment()
	default:
		logger, err = zap.NewNop(), nil
	}
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			if err := logger.Sync(); !errors.Is(err, syscall.ENOTTY) {
				return err
			}

			return nil
		},
	})

	return logger, nil
}
