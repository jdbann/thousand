package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"emailaddress.horse/thousand/app"
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
				Name:        "database-url",
				Usage:       "override the default DB connection",
				Value:       thousand.DatabaseURL,
				Destination: &cliConfig.DatabaseURL,
			},
		},
		Before: func(c *cli.Context) error {
			*thousand = *app.NewApp(app.DevelopmentConfig, cliConfig)

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
		},
	}

	if err := cli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
