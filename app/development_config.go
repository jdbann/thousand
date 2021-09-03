package app

import "github.com/labstack/echo/v4/middleware"

// DevelopmentConfig sets up the app for a development environment.
func DevelopmentConfig(app *App) {
	app.Debug = true
	app.Logger = middleware.Logger()
}
