package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"emailaddress.horse/thousand/app"
	"emailaddress.horse/thousand/repository"
	"github.com/urfave/cli/v2"
)

func BuildCLIApp() *cli.App {
	var thousand *app.App

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
		},
		Before: func(c *cli.Context) error {
			repo, err := repository.New(repository.Options{
				DatabaseURL: databaseURL(c),
			})
			if err != nil {
				return err
			}

			thousand = app.NewApp(app.Options{
				Debug:      c.Bool("debug"),
				Repository: repo,
			})

			return nil
		},
		Action: func(c *cli.Context) error {
			port := os.Getenv("PORT")
			if port == "" {
				port = "4000"
			}

			addr := fmt.Sprintf(":%s", port)
			if err := thousand.Start(addr); err != nil {
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

					routes := thousand.Routes()

					sort.Slice(routes, func(i, j int) bool {
						if routes[i].Path == routes[j].Path {
							return methodOrder[routes[i].Method] < methodOrder[routes[j].Method]
						}

						return routes[i].Path < routes[j].Path
					})

					writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
					fmt.Fprintf(writer, "Method\tPath\tName\n")
					for _, route := range routes {
						fmt.Fprintf(writer, "%s\t%s\t%s\n", route.Method, route.Path, route.Name)
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
