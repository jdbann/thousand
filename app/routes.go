package app

import (
	"emailaddress.horse/thousand/static"
	"github.com/labstack/echo/v4/middleware"
)

func (app *App) setupRoutes() {
	app.Use(app.Logger)
	app.Use(static.Middleware())

	app.Pre(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
		Getter: middleware.MethodFromForm("_method"),
	}))

	app.GET("/", app.root).Name = "root"

	// Details
	app.POST("/details", app.createDetails).Name = "create-details"

	// Memories
	app.DELETE("/memories/:id", app.deleteMemory).Name = "delete-memory"
	app.POST("/memories/:id/experiences", app.createExperience).Name = "create-experience"

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
