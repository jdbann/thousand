package app

import (
	"net/http"

	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/static"
	"github.com/labstack/echo/v4/middleware"
)

func (app *App) setupRoutes() {
	app.Use(app.LoggerMiddleware)

	app.Pre(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))

	app.Pre(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
		Getter: middleware.MethodFromForm("_method"),
	}))

	handlers.Root(app.Echo)

	app.Group("/assets", static.Middleware())

	// Vampires
	handlers.ListVampires(app.Echo, app.Repository)
	handlers.NewVampire(app.Echo)
	handlers.CreateVampire(app.Echo, app.Repository)
	handlers.ShowVampire(app.Echo, app.Repository)
	handlers.NewExperience(app.Echo, app.Repository)
	handlers.CreateExperience(app.Echo, app.Repository)
	handlers.NewSkill(app.Echo, app.Repository)
	handlers.CreateSkill(app.Echo, app.Repository)
	handlers.NewResource(app.Echo, app.Repository)
	handlers.CreateResource(app.Echo, app.Repository)
	handlers.NewCharacter(app.Echo, app.Repository)
	handlers.CreateCharacter(app.Echo, app.Repository)
	handlers.NewMark(app.Echo, app.Repository)
	handlers.CreateMark(app.Echo, app.Repository)
}
