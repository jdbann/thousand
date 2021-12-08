package app

import (
	"net/http"

	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/static"
	"github.com/go-chi/chi/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (app *App) setupRoutes() {
	app.echo.Use(loggerMiddleware(app.logger))

	app.echo.Pre(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))

	app.echo.Pre(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
		Getter: middleware.MethodFromForm("_method"),
	}))

	r := chi.NewRouter()

	// Temporarily specify routes whilst we still route through Echo to avoid
	// routing conflicts between Echo and Chi
	app.echo.GET("/", echo.WrapHandler(r))
	app.echo.GET("/vampires", echo.WrapHandler(r))
	app.echo.GET("/vampires/new", echo.WrapHandler(r))

	handlers.Root(r)

	handlers.ListVampires(r, app.logger, app.renderer, app.repository)
	handlers.NewVampire(r, app.logger, app.renderer)

	app.echo.Group("/assets", static.Middleware())

	handlers.CreateVampire(app.echo, app.repository)
	handlers.ShowVampire(app.echo, app.repository)
	handlers.NewExperience(app.echo, app.repository)
	handlers.CreateExperience(app.echo, app.repository)
	handlers.NewSkill(app.echo, app.repository)
	handlers.CreateSkill(app.echo, app.repository)
	handlers.NewResource(app.echo, app.repository)
	handlers.CreateResource(app.echo, app.repository)
	handlers.NewCharacter(app.echo, app.repository)
	handlers.CreateCharacter(app.echo, app.repository)
	handlers.NewMark(app.echo, app.repository)
	handlers.CreateMark(app.echo, app.repository)
}
