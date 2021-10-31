package app

import (
	"net/http"

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
	app.GET("/vampires", app.listVampires).Name = "list-vampires"
	app.GET("/vampires/new", app.newVampire).Name = "new-vampire"
	app.POST("/vampires", app.createVampire).Name = "create-vampire"
	app.GET("/vampires/:id", app.showVampire).Name = "show-vampire"
	app.GET("/vampires/:vampireID/memories/:id/experiences/new", app.newExperience).Name = "new-experience"
	app.POST("/vampires/:vampireID/memories/:id/experiences", app.createExperience).Name = "create-experience"

	// Details
	app.POST("/details", app.createDetails).Name = "create-details"

	// Memories
	app.DELETE("/memories/:id", app.deleteMemory).Name = "delete-memory"
	app.POST("/memories/:id/experiences", app.oldCreateExperience).Name = "old-create-experience"

	// Skills
	app.POST("/skills", app.createSkill).Name = "create-skill"
	app.PATCH("/skills/:id", app.updateSkill).Name = "update-skill"

	// Resources
	app.POST("/resources", app.createResource).Name = "create-resource"
	app.PATCH("/resources/:id", app.updateResource).Name = "update-resource"

	// Characters
	app.POST("/characters", app.createCharacter).Name = "create-character"
	app.PATCH("/characters/:id", app.updateCharacter).Name = "update-character"
	app.POST("/characters/:id/descriptor", app.createDescriptor).Name = "create-descriptor"

	// Marks
	app.POST("/marks", app.createMark).Name = "create-mark"
}
