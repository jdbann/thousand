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

	app.GET("/", app.root).Name = "root"

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
	app.GET("/vampires/:vampireID/marks/new", app.newMark).Name = "new-mark"
	app.POST("/vampires/:vampireID/marks", app.createMark).Name = "create-mark"

	// Details
	app.POST("/details", app.createDetails).Name = "create-details"

	// Memories
	app.DELETE("/memories/:id", app.deleteMemory).Name = "delete-memory"
	app.POST("/memories/:id/experiences", app.oldCreateExperience).Name = "old-create-experience"

	// Skills
	app.POST("/skills", app.oldCreateSkill).Name = "old-create-skill"
	app.PATCH("/skills/:id", app.updateSkill).Name = "update-skill"

	// Resources
	app.POST("/resources", app.oldCreateResource).Name = "old-create-resource"
	app.PATCH("/resources/:id", app.updateResource).Name = "update-resource"

	// Characters
	app.POST("/characters", app.oldCreateCharacter).Name = "old-create-character"
	app.PATCH("/characters/:id", app.updateCharacter).Name = "update-character"
	app.POST("/characters/:id/descriptor", app.createDescriptor).Name = "create-descriptor"

	// Marks
	app.POST("/marks", app.oldCreateMark).Name = "old-create-mark"
}
