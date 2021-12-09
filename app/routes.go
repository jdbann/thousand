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
	app.echo.POST("/vampires", echo.WrapHandler(r))
	app.echo.GET("/vampires/new", echo.WrapHandler(r))
	app.echo.GET("/vampires/:id", echo.WrapHandler(r))
	app.echo.GET("/vampires/:id/memories/:memoryID/experiences/new", echo.WrapHandler(r))
	app.echo.POST("/vampires/:id/memories/:memoryID/experiences", echo.WrapHandler(r))
	app.echo.GET("/vampires/:id/skills/new", echo.WrapHandler(r))
	app.echo.POST("/vampires/:id/skills", echo.WrapHandler(r))

	handlers.Root(r)

	handlers.ListVampires(r, app.logger, app.renderer, app.repository)
	handlers.NewVampire(r, app.logger, app.renderer)
	handlers.CreateVampire(r, app.logger, app.repository)
	handlers.ShowVampire(r, app.logger, app.renderer, app.repository)
	handlers.NewExperience(r, app.logger, app.renderer, app.repository)
	handlers.CreateExperience(r, app.logger, app.repository)
	handlers.NewSkill(r, app.logger, app.renderer, app.repository)
	handlers.CreateSkill(r, app.logger, app.repository)

	app.echo.Group("/assets", static.Middleware())

	handlers.NewResource(app.echo, app.repository)
	handlers.CreateResource(app.echo, app.repository)
	handlers.NewCharacter(app.echo, app.repository)
	handlers.CreateCharacter(app.echo, app.repository)
	handlers.NewMark(app.echo, app.repository)
	handlers.CreateMark(app.echo, app.repository)
}
