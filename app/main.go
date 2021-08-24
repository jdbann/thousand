package app

import (
	"html/template"
	"log"
	"net/http"

	"emailaddress.horse/thousand/app/models"
	"emailaddress.horse/thousand/templates"
	"github.com/gin-gonic/gin"
)

// App is a configured instance of the application, ready to be served by a
// server or interacted with by CLI commands.
type App struct {
	Character *models.Character
}

// NewApp configures an instance of the application with helpful defaults.
func NewApp() *App {
	return &App{
		Character: &models.Character{},
	}
}

// Engine returns the configured set of routes for the app to be used by an HTTP
// server.
func (app *App) Engine() http.Handler {
	r := gin.Default()

	templates, err := template.ParseFS(templates.Templates, "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	r.SetHTMLTemplate(templates)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", app.Character)
	})

	r.POST("/details", func(c *gin.Context) {
		if err := c.ShouldBind(app.Character); err != nil {
			_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		c.Redirect(http.StatusFound, "/")
	})

	return r
}
