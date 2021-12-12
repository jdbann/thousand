package browser

import (
	"context"
	"testing"

	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/middleware"
	"emailaddress.horse/thousand/repository"
	"emailaddress.horse/thousand/server"
	"emailaddress.horse/thousand/templates"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type params struct {
	fx.Out

	DatabaseURL string `name:"databaseURL"`
	Port        int    `name:"port"`

	T *testing.T
}

var testModule = fx.Options(
	fx.Provide(fx.Annotate(chi.NewMux, fx.As(new(chi.Router)))),

	middleware.Module,

	handlers.Module,
	server.Module,
	templates.Module,

	fx.Provide(func(t *testing.T) *zap.Logger {
		return newLogger(t)
	}),

	fx.Provide(fxTestRepository),

	fx.NopLogger,
)

type trParams struct {
	fx.In

	DatabaseURL string `name:"databaseURL"`

	Logger *zap.Logger
}

func fxTestRepository(t *testing.T, params trParams) *repository.Repository {
	repo, err := repository.New(repository.Options{
		DatabaseURL: params.DatabaseURL,
		Logger:      params.Logger.Named("repository"),
	})
	if err != nil {
		t.Fatal(err)
	}

	repo, tx, err := repo.WithTx(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := tx.Rollback(context.Background()); err != nil {
			t.Fatal(err)
		}
	})

	return repo
}
