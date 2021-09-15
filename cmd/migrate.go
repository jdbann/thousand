package main

import (
	"github.com/pressly/goose"
	"github.com/urfave/cli/v2"
)

const (
	migrationsPath = "./db/migrations"
)

func createMigration(c *cli.Context) error {
	return goose.Create(nil, migrationsPath, c.Args().Get(0), "sql")
}
