package app

import (
	"emailaddress.horse/thousand/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/labstack/echo/v4"
)

func (app *App) setupRoutes() {
	r := chi.NewRouter()

	r.Use(RequestLogger(app.logger.Named("server")))
	r.Use(MethodOverride)
	r.Use(middleware.RedirectSlashes)

	handlers.Assets(r, app.assets)

	handlers.Root(r)

	handlers.NewCharacter(r, app.logger, app.renderer, app.repository)
	handlers.CreateCharacter(r, app.logger, app.repository)

	handlers.NewExperience(r, app.logger, app.renderer, app.repository)
	handlers.CreateExperience(r, app.logger, app.repository)

	handlers.NewMark(r, app.logger, app.renderer, app.repository)
	handlers.CreateMark(r, app.logger, app.repository)

	handlers.NewResource(r, app.logger, app.renderer, app.repository)
	handlers.CreateResource(r, app.logger, app.repository)

	handlers.NewSkill(r, app.logger, app.renderer, app.repository)
	handlers.CreateSkill(r, app.logger, app.repository)

	handlers.ListVampires(r, app.logger, app.renderer, app.repository)
	handlers.NewVampire(r, app.logger, app.renderer)
	handlers.CreateVampire(r, app.logger, app.repository)
	handlers.ShowVampire(r, app.logger, app.renderer, app.repository)

	app.echo.Any("*", echo.WrapHandler(r))
}
