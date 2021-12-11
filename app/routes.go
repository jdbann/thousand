package app

import (
	"emailaddress.horse/thousand/handlers"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *App) setupRoutes() {
	app.mux.Use(RequestLogger(app.logger.Named("server")))
	app.mux.Use(MethodOverride)
	app.mux.Use(middleware.RedirectSlashes)

	handlers.Assets(app.mux, app.assets)

	handlers.Root(app.mux)

	handlers.NewCharacter(app.mux, app.logger, app.renderer, app.repository)
	handlers.CreateCharacter(app.mux, app.logger, app.repository)

	handlers.NewExperience(app.mux, app.logger, app.renderer, app.repository)
	handlers.CreateExperience(app.mux, app.logger, app.repository)

	handlers.NewMark(app.mux, app.logger, app.renderer, app.repository)
	handlers.CreateMark(app.mux, app.logger, app.repository)

	handlers.NewResource(app.mux, app.logger, app.renderer, app.repository)
	handlers.CreateResource(app.mux, app.logger, app.repository)

	handlers.NewSkill(app.mux, app.logger, app.renderer, app.repository)
	handlers.CreateSkill(app.mux, app.logger, app.repository)

	handlers.ListVampires(app.mux, app.logger, app.renderer, app.repository)
	handlers.NewVampire(app.mux, app.logger, app.renderer)
	handlers.CreateVampire(app.mux, app.logger, app.repository)
	handlers.ShowVampire(app.mux, app.logger, app.renderer, app.repository)
}
