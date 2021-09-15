package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"emailaddress.horse/thousand/app"
	"github.com/urfave/cli/v2"
)

var thousand app.App

func main() {
	// Setup an app.CLIConfig struct to receive config values from CLI flags for
	// application before any actions are performed.
	var cliConfig app.CLIConfig

	cli := &cli.App{
		Name:  "thousand",
		Usage: "I forget why I made this...",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "database-url",
				Usage:       "override the default DB connection",
				Destination: &cliConfig.DatabaseURL,
				EnvVars:     []string{"DATABASE_URL"},
			},
			&cli.StringFlag{
				Name:  "environment",
				Usage: "specify the environment to act as",
				Value: "development",
			},
		},
		Before: func(c *cli.Context) error {
			requestedEnv := c.String("environment")

			envConfig, err := app.ConfigFor(requestedEnv)
			if err != nil {
				if errors.Is(app.ErrUnrecognisedEnvironment, err) {
					return cli.Exit(fmt.Sprintf("Unrecognised environment: %q", requestedEnv), 0)
				}

				return err
			}

			thousand = *app.NewApp(envConfig, cliConfig)

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
		},
	}

	if err := cli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
