package app

import (
	"emailaddress.horse/thousand/static"
	"github.com/labstack/echo/v4/middleware"
)

func (app *App) setupRoutes() {
	app.Use(middleware.Logger())
	app.Use(static.Middleware())

	app.Pre(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
		Getter: middleware.MethodFromForm("_method"),
	}))

	app.GET("/", app.root).Name = "root"

	// Details
	app.POST("/details", app.createDetails).Name = "create-details"

	// Memories
	app.POST("/memories/:id/experiences", app.createExperience).Name = "create-experience"

	// Skills
	app.POST("/skills", app.createSkill).Name = "create-skill"
	app.PATCH("/skills/:id", app.updateSkill).Name = "update-skill"
}
