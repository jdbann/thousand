package app

import (
	"net/http"

	"emailaddress.horse/thousand/repository"
	"emailaddress.horse/thousand/templates"
	"github.com/labstack/echo/v4"
)

// App is a configured instance of the application, ready to be served by a
// server or interacted with by CLI commands.
type App struct {
	echo *echo.Echo

	// Injected middleware
	loggerMiddleware echo.MiddlewareFunc

	// Runtime values
	repository *repository.Repository
}

type Options struct {
	Debug      bool
	Repository *repository.Repository
}

// NewApp configures an instance of the application with helpful defaults.
func NewApp(opts Options) *App {
	echo := echo.New()
	echo.Renderer = templates.NewRenderer(echo)
	echo.Debug = opts.Debug

	app := &App{
		echo:       echo,
		repository: opts.Repository,
	}

	app.setupRoutes()

	return app
}

func (a *App) Start(addr string) error {
	return a.echo.Start(addr)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.echo.ServeHTTP(w, r)
}

func (a *App) Routes() []*echo.Route {
	return a.echo.Routes()
}
