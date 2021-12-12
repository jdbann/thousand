package handlers

import (
	"time"

	"emailaddress.horse/thousand/repository"
	"emailaddress.horse/thousand/static"
	"emailaddress.horse/thousand/templates"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Invoke(register),
)

type RegisterParams struct {
	fx.In

	Logger     *zap.Logger
	Renderer   *templates.Renderer
	Repository *repository.Repository
	Router     chi.Router
}

func register(p RegisterParams) {
	Assets(p.Router, static.Assets)
	Health(p.Router, p.Logger, p.Repository, time.Now)

	Root(p.Router)

	NewCharacter(p.Router, p.Logger, p.Renderer, p.Repository)
	CreateCharacter(p.Router, p.Logger, p.Repository)

	NewExperience(p.Router, p.Logger, p.Renderer, p.Repository)
	CreateExperience(p.Router, p.Logger, p.Repository)

	NewMark(p.Router, p.Logger, p.Renderer, p.Repository)
	CreateMark(p.Router, p.Logger, p.Repository)

	NewResource(p.Router, p.Logger, p.Renderer, p.Repository)
	CreateResource(p.Router, p.Logger, p.Repository)

	NewSkill(p.Router, p.Logger, p.Renderer, p.Repository)
	CreateSkill(p.Router, p.Logger, p.Repository)

	ListVampires(p.Router, p.Logger, p.Renderer, p.Repository)
	NewVampire(p.Router, p.Logger, p.Renderer)
	CreateVampire(p.Router, p.Logger, p.Repository)
	ShowVampire(p.Router, p.Logger, p.Renderer, p.Repository)
}
