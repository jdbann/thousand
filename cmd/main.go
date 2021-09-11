package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"text/tabwriter"

	"emailaddress.horse/thousand/app"
	"github.com/lib/pq"
	"github.com/urfave/cli/v2"
)

func main() {
	// Initially setup the app in the development environment to allow presenting
	// development config values as default flag values.
	thousand := app.NewApp(app.DevelopmentConfig)

	// Setup an app.CLIConfig struct to receive config values from CLI flags for
	// application before any actions are performed.
	var cliConfig app.CLIConfig

	cli := &cli.App{
		Name:  "thousand",
		Usage: "I forget why I made this...",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "environment",
				Aliases: []string{"env"},
				Usage:   "override the default environment",
				Value:   "development",
			},
			&cli.StringFlag{
				Name:        "database_url",
				Usage:       "override the default DB connection",
				Value:       thousand.DatabaseURL,
				Destination: &cliConfig.DatabaseURL,
			},
		},
		Before: func(c *cli.Context) error {
			env := c.String("environment")

			envConfig, err := app.ConfigFor(env)
			if err != nil {
				if errors.Is(app.ErrUnrecognisedEnvironment, err) {
					return cli.Exit(fmt.Sprintf("Unrecognised environment: %q", env), 0)
				}
			}

			thousand = app.NewApp(cliConfig, envConfig)
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
						Name:  "create",
						Usage: "create the database",
						Action: func(c *cli.Context) error {
							url, err := url.Parse(thousand.DatabaseURL)
							if err != nil {
								return err
							}

							dbName := url.Path
							if len(dbName) > 0 && dbName[:1] == "/" {
								dbName = dbName[1:]
							}

							url.Path = "postgres"

							db, err := sql.Open("postgres", url.String())
							if err != nil {
								return err
							}

							_, err = db.Exec(fmt.Sprintf("CREATE database %s", pq.QuoteIdentifier(dbName)))
							if err != nil {
								return err
							}

							fmt.Printf("Created database: %q\n", dbName)
							return nil
						},
					},
					{
						Name:  "drop",
						Usage: "drop the database",
						Action: func(c *cli.Context) error {
							url, err := url.Parse(thousand.DatabaseURL)
							if err != nil {
								return err
							}

							dbName := url.Path
							if len(dbName) > 0 && dbName[:1] == "/" {
								dbName = dbName[1:]
							}

							url.Path = "postgres"

							db, err := sql.Open("postgres", url.String())
							if err != nil {
								return err
							}

							_, err = db.Exec(fmt.Sprintf("DROP database %s", pq.QuoteIdentifier(dbName)))
							if err != nil {
								return err
							}

							fmt.Printf("Dropped database: %q\n", dbName)
							return nil
						},
					},
				},
			},
		},
	}

	if err := cli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
