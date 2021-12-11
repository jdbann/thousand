package app

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"strconv"
	"time"

	"emailaddress.horse/thousand/repository"
	"emailaddress.horse/thousand/templates"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// App is a configured instance of the application, ready to be served by a
// server or interacted with by CLI commands.
type App struct {
	address    string
	assets     fs.FS
	logger     *zap.Logger
	mux        *chi.Mux
	renderer   *templates.Renderer
	repository *repository.Repository
	server     server
}

type server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

type Options struct {
	Assets     fs.FS
	Debug      bool
	Host       string
	Logger     *zap.Logger
	Mux        *chi.Mux
	Port       int
	Repository *repository.Repository
	Server     server
}

// NewApp configures an instance of the application with helpful defaults.
func NewApp(opts Options) *App {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}

	if opts.Mux == nil {
		opts.Mux = chi.NewMux()
	}

	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))

	if opts.Server == nil {
		opts.Server = &http.Server{
			Addr:    address,
			Handler: opts.Mux,
		}
	}

	app := &App{
		address:    address,
		assets:     opts.Assets,
		logger:     opts.Logger,
		mux:        opts.Mux,
		repository: opts.Repository,
		renderer:   templates.NewRenderer(),
		server:     opts.Server,
	}

	return app
}

func (a *App) Start() error {
	a.setupRoutes()

	a.logger.Info("starting", zap.String("address", a.address))
	if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

func (a *App) Stop() error {
	a.logger.Info("stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}

	return nil
}

type Route struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

func (a *App) Routes() []*Route {
	a.setupRoutes()

	routes := []*Route{}

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		routes = append(routes, &Route{
			Method: method,
			Path:   route,
		})
		return nil
	}

	if err := chi.Walk(a.mux, walkFunc); err != nil {
		panic(err)
	}

	return routes
}
