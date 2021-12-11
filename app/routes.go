package app

import (
	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/static"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (app *App) setupRoutes() {
	app.echo.Pre(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
		Getter: middleware.MethodFromForm("_method"),
	}))

	r := chi.NewRouter()

	r.Use(RequestLogger(app.logger.Named("server")))
	r.Use(chimiddleware.RedirectSlashes)

	// Temporarily specify routes whilst we still route through Echo to avoid
	// routing conflicts between Echo and Chi
	app.echo.Any("*", echo.WrapHandler(r))

	handlers.Root(r)

	handlers.ListVampires(r, app.logger, app.renderer, app.repository)
	handlers.NewVampire(r, app.logger, app.renderer)
	handlers.CreateVampire(r, app.logger, app.repository)
	handlers.ShowVampire(r, app.logger, app.renderer, app.repository)
	handlers.NewExperience(r, app.logger, app.renderer, app.repository)
	handlers.CreateExperience(r, app.logger, app.repository)
	handlers.NewSkill(r, app.logger, app.renderer, app.repository)
	handlers.CreateSkill(r, app.logger, app.repository)
	handlers.NewResource(r, app.logger, app.renderer, app.repository)
	handlers.CreateResource(r, app.logger, app.repository)
	handlers.NewCharacter(r, app.logger, app.renderer, app.repository)
	handlers.CreateCharacter(r, app.logger, app.repository)
	handlers.NewMark(r, app.logger, app.renderer, app.repository)
	handlers.CreateMark(r, app.logger, app.repository)

	app.echo.Group("/assets", static.Middleware())
}
