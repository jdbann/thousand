package main

import (
	"database/sql"

	"github.com/pressly/goose"
	"github.com/urfave/cli/v2"
)

const (
	migrationsPath = "./db/migrations"
)

func createMigration(c *cli.Context) error {
	return goose.Create(nil, migrationsPath, c.Args().Get(0), "sql")
}

func runMigrations(c *cli.Context) error {
	db, err := sql.Open("postgres", thousand.DatabaseURL)
	if err != nil {
		return err
	}

	return goose.Up(db, migrationsPath)
}

func rollbackMigrations(c *cli.Context) error {
	db, err := sql.Open("postgres", thousand.DatabaseURL)
	if err != nil {
		return err
	}

	return goose.Down(db, migrationsPath)
}
