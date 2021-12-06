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
	handlers.ListVampires(app.Echo, app.Models)
	handlers.NewVampire(app.Echo)
	handlers.CreateVampire(app.Echo, app.Models)
	handlers.ShowVampire(app.Echo, app.Models)
	handlers.NewExperience(app.Echo, app.Models)
	handlers.CreateExperience(app.Echo, app.Models)
	handlers.NewSkill(app.Echo, app.Models)
	handlers.CreateSkill(app.Echo, app.Models)
	handlers.NewResource(app.Echo, app.Models)
	handlers.CreateResource(app.Echo, app.Models)
	handlers.NewCharacter(app.Echo, app.Models)
	handlers.CreateCharacter(app.Echo, app.Models)
	handlers.NewMark(app.Echo, app.Models)
	handlers.CreateMark(app.Echo, app.Models)
}
