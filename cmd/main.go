package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"emailaddress.horse/thousand/repository"
	"emailaddress.horse/thousand/server"
	"emailaddress.horse/thousand/static"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func BuildCLIApp() *cli.App {
	return &cli.App{
		Name:  "thousand",
		Usage: "I forget why I made this...",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Usage:   "enable debug mode for detailed logging",
				Value:   true,
				EnvVars: []string{"DEBUG"},
			},
			&cli.StringFlag{
				Name:    "database-url",
				Usage:   "override the default DB connection",
				EnvVars: []string{"DATABASE_URL"},
			},
			&cli.StringFlag{
				Name:  "environment",
				Usage: "specify the environment to act as",
				Value: "development",
			},
			&cli.IntFlag{
				Name:  "port",
				Usage: "port to run the server on",
				Value: 4000,
			},
		},
		Action: func(c *cli.Context) error {
			logger, err := buildLogger(c)
			if err != nil {
				return err
			}
			defer logger.Sync()

			repo, err := repository.New(repository.Options{
				DatabaseURL: databaseURL(c),
				Logger:      logger.Named("repository"),
			})
			if err != nil {
				return err
			}

			thousand := server.New(server.Options{
				Port:       c.Int("port"),
				Assets:     static.Assets,
				Debug:      c.Bool("debug"),
				Logger:     logger,
				Repository: repo,
			})
			if err := thousand.Start(); err != nil {
				return err
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "routes",
				Usage: "present a list of routes the app handles",
				Action: func(c *cli.Context) error {
					var methodOrder = map[string]int{
						"GET":     0,
						"POST":    1,
						"PUT":     2,
						"PATCH":   3,
						"DELETE":  4,
						"HEAD":    5,
						"CONNECT": 6,
						"OPTIONS": 7,
						"TRACE":   8,
					}

					routes := server.New(server.Options{}).Routes()

					sort.Slice(routes, func(i, j int) bool {
						if routes[i].Path == routes[j].Path {
							return methodOrder[routes[i].Method] < methodOrder[routes[j].Method]
						}

						return routes[i].Path < routes[j].Path
					})

					writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
					fmt.Fprintf(writer, "Method\tPath\n")
					for _, route := range routes {
						fmt.Fprintf(writer, "%s\t%s\n", route.Method, route.Path)
					}
					writer.Flush()

					return nil
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

func databaseURL(c *cli.Context) string {
	var databaseURL string
	switch c.String("environment") {
	case "development":
		databaseURL = "postgres://localhost:5432/thousand_development?sslmode=disable"
	case "test":
		databaseURL = "postgres://localhost:5432/thousand_test?sslmode=disable"
	}

	if c.String("database-url") != "" {
		databaseURL = c.String("database-url")
	}

	return databaseURL
}

func buildLogger(c *cli.Context) (*zap.Logger, error) {
	switch c.String("environment") {
	case "production":
		return zap.NewProduction()
	case "development":
		return zap.NewDevelopment()
	case "test":
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}
}
