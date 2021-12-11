package app

import (
	"net/http"

	"emailaddress.horse/thousand/repository"
	"emailaddress.horse/thousand/templates"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// App is a configured instance of the application, ready to be served by a
// server or interacted with by CLI commands.
type App struct {
	echo       *echo.Echo
	logger     *zap.Logger
	repository *repository.Repository
	renderer   *templates.Renderer
}

type Options struct {
	Debug      bool
	Logger     *zap.Logger
	Repository *repository.Repository
}

// NewApp configures an instance of the application with helpful defaults.
func NewApp(opts Options) *App {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}

	echo := echo.New()
	echo.Debug = opts.Debug

	app := &App{
		echo:       echo,
		logger:     opts.Logger,
		repository: opts.Repository,
		renderer:   templates.NewRenderer(),
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
