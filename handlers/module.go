package handlers

import (
	"emailaddress.horse/thousand/health"
	"emailaddress.horse/thousand/middleware"
	"emailaddress.horse/thousand/repository"
	"emailaddress.horse/thousand/session"
	"emailaddress.horse/thousand/static"
	"emailaddress.horse/thousand/templates"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Invoke(fxRegister),
)

type RegisterParams struct {
	fx.In

	Health     *health.Health
	Logger     *zap.Logger
	Renderer   *templates.Renderer
	Repository *repository.Repository
	Router     chi.Router
	Store      *session.Store
}

func fxRegister(p RegisterParams) {
	Assets(p.Router, static.Assets)
	Health(p.Router, p.Logger, p.Health)

	Root(p.Router)

	NewSession(p.Router, p.Logger, p.Renderer)
	CreateSession(p.Router, p.Logger, p.Repository, p.Renderer, p.Store)
	DestroySession(p.Router, p.Logger, p.Store)

	NewUser(p.Router, p.Logger, p.Renderer)
	CreateUser(p.Router, p.Logger, p.Repository, p.Renderer, p.Store)

	p.Router.Group(func(r chi.Router) {
		middleware.EnsureLoggedIn(r, p.Store, p.Repository)

		NewCharacter(r, p.Logger, p.Renderer, p.Repository)
		CreateCharacter(r, p.Logger, p.Repository)

		NewExperience(r, p.Logger, p.Renderer, p.Repository)
		CreateExperience(r, p.Logger, p.Repository)

		NewMark(r, p.Logger, p.Renderer, p.Repository)
		CreateMark(r, p.Logger, p.Repository)

		NewResource(r, p.Logger, p.Renderer, p.Repository)
		CreateResource(r, p.Logger, p.Repository)

		NewSkill(r, p.Logger, p.Renderer, p.Repository)
		CreateSkill(r, p.Logger, p.Repository)

		ListVampires(r, p.Logger, p.Renderer, p.Repository)
		NewVampire(r, p.Logger, p.Renderer)
		CreateVampire(r, p.Logger, p.Repository)
		ShowVampire(r, p.Logger, p.Renderer, p.Repository)
	})
}
