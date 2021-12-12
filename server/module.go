package server

import (
	"emailaddress.horse/thousand/repository"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(fxNew),
)

type Params struct {
	fx.In

	Host string `name:"host" optional:"true"`
	Port int    `name:"port"`

	Logger     *zap.Logger
	Repository *repository.Repository
	Router     chi.Router
}

func fxNew(lc fx.Lifecycle, params Params) *Server {
	server := New(Options{
		Host:   params.Host,
		Logger: params.Logger.Named("server"),
		Port:   params.Port,
		Router: params.Router,
	})

	lc.Append(fx.Hook{
		OnStart: server.Listen,
	})

	lc.Append(fx.Hook{
		OnStart: server.Start,
		OnStop:  server.Stop,
	})

	return server
}
