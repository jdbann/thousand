package cmd

import (
	"database/sql"

	"emailaddress.horse/thousand/db"
	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v2"
)

func init() {
	goose.SetBaseFS(db.Migrations)
}

const (
	migrationsPath = "./db/migrations"
)

func createMigration(c *cli.Context) error {
	return goose.Create(nil, migrationsPath, c.Args().Get(0), "sql")
}

func runMigrations(c *cli.Context) error {
	conn, err := sql.Open("pgx", thousand.DatabaseURL)
	if err != nil {
		return err
	}

	return goose.Up(conn, db.FSMigrationsPath)
}

func rollbackMigrations(c *cli.Context) error {
	conn, err := sql.Open("pgx", thousand.DatabaseURL)
	if err != nil {
		return err
	}

	return goose.Down(conn, db.FSMigrationsPath)
}

func migrationsStatus(c *cli.Context) error {
	conn, err := sql.Open("pgx", thousand.DatabaseURL)
	if err != nil {
		return err
	}

	return goose.Status(conn, db.FSMigrationsPath)
}
