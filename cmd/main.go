package cmd

import (
	"context"

	"emailaddress.horse/thousand/handlers"
	"emailaddress.horse/thousand/health"
	"emailaddress.horse/thousand/logger"
	"emailaddress.horse/thousand/middleware"
	"emailaddress.horse/thousand/registry"
	"emailaddress.horse/thousand/repository"
	"emailaddress.horse/thousand/server"
	"emailaddress.horse/thousand/templates"
	"github.com/go-chi/chi/v5"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func BuildCLIApp() *cli.App {
	return &cli.App{
		Name:  "thousand",
		Usage: "I forget why I made this...",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "database-url",
				Usage:   "override the default DB connection",
				Value:   "postgres://localhost:5432/thousand_development?sslmode=disable",
				EnvVars: []string{"DATABASE_URL"},
			},
			&cli.StringFlag{
				Name:    "log-format",
				Usage:   "format for logs: `prod` uses structured logging; `dev` uses readable logging",
				Value:   "dev",
				EnvVars: []string{"LOG_FORMAT"},
			},
			&cli.IntFlag{
				Name:    "port",
				Usage:   "port to run the server on",
				Value:   4000,
				EnvVars: []string{"PORT"},
			},
		},
		Action: func(c *cli.Context) error {
			a := fx.New(
				fx.Supply(
					struct {
						fx.Out

						DatabaseURL string `name:"databaseURL"`
						Host        string `name:"host" optional:"true"`
						LogFormat   string `name:"logFormat" optional:"true"`
						Port        int    `name:"port"`
					}{
						DatabaseURL: c.String("database-url"),
						LogFormat:   c.String("log-format"),
						Port:        c.Int("port"),
					},
				),

				fx.Provide(fx.Annotate(chi.NewMux, fx.As(new(chi.Router)))),
				middleware.Module,

				handlers.Module,
				health.Module,
				logger.Module,
				registry.Module,
				repository.Module,
				server.Module,
				templates.Module,

				fx.Invoke(func(s *server.Server) {}),
			)

			a.Run()

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "routes",
				Usage: "present a list of routes the app handles",
				Action: func(c *cli.Context) error {
					a := fx.New(
						fx.Provide(fx.Annotate(chi.NewMux, fx.As(new(chi.Router)))),

						fx.Provide(func() *health.Health { return nil }),
						fx.Provide(func() *repository.Repository { return nil }),
						fx.Provide(func() *templates.Renderer { return nil }),
						fx.Provide(func() *zap.Logger { return nil }),

						handlers.Module,

						fx.NopLogger,

						fx.Invoke(server.PrintRoutes),
					)

					return a.Start(context.Background())
				},
			},
			{
				Name:  "db",
				Usage: "setup and migrate the database",
				Subcommands: []*cli.Command{
					{
						Name:   "create",
						Usage:  "create the database",
						Action: createDatabase,
					},
					{
						Name:   "drop",
						Usage:  "drop the database",
						Action: dropDatabase,
					},
				},
			},
			{
				Name:  "migrate",
				Usage: "manage migrations",
				Subcommands: []*cli.Command{
					{
						Name:   "create",
						Usage:  "create new SQL migration",
						Action: createMigration,
					},
					{
						Name:   "run",
						Usage:  "run pending migrations",
						Action: runMigrations,
					},
					{
						Name:   "rollback",
						Usage:  "rollback latest migration",
						Action: rollbackMigrations,
					},
					{
						Name:   "status",
						Usage:  "report current status of migrations",
						Action: migrationsStatus,
					},
				},
			},
		},
	}
}
