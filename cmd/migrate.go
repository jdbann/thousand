package main

import (
	"database/sql"
	"errors"

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

		sqlDB, ok := db.(*sql.DB)
		if !ok {
			return errors.New("App's configured DB is not a valid *sql.DB")
		}

		return goose.Create(sqlDB, migrationsPath, c.Args().Get(0), "sql")
	}
}

func runMigrations(thousand *app.App) cli.ActionFunc {
	return func(c *cli.Context) error {
		db, err := thousand.DB()
		if err != nil {
			return err
		}

		sqlDB, ok := db.(*sql.DB)
		if !ok {
			return errors.New("App's configured DB is not a valid *sql.DB")
		}

		return goose.Up(sqlDB, migrationsPath)
	}
}
