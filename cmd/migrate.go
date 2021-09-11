package main

import (
	"emailaddress.horse/thousand/app"
	"github.com/pressly/goose"
	"github.com/urfave/cli/v2"
)

const (
	migrationsPath = "./db/migrations"
)

func createMigration(thousand *app.App) cli.ActionFunc {
	return func(c *cli.Context) error {
		db, err := thousand.DB()
		if err != nil {
			return err
		}

		return goose.Create(db, migrationsPath, c.Args().Get(0), "sql")
	}
}

func runMigrations(thousand *app.App) cli.ActionFunc {
	return func(c *cli.Context) error {
		db, err := thousand.DB()
		if err != nil {
			return err
		}

		return goose.Up(db, migrationsPath)
	}
}
