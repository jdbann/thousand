package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v4"
	"github.com/urfave/cli/v2"
)

func createDatabase(c *cli.Context) error {
	dbName, err := getDatabaseName(thousand.DatabaseURL)
	if err != nil {
		return err
	}

	err = execOnPostgresDB(c.Context, fmt.Sprintf("CREATE database %s", dbName))
	if err != nil {
		return err
	}

	fmt.Printf("Created database: %q\n", dbName)
	return nil
}

func dropDatabase(c *cli.Context) error {
	dbName, err := getDatabaseName(thousand.DatabaseURL)
	if err != nil {
		return err
	}

	err = execOnPostgresDB(c.Context, fmt.Sprintf("DROP database %s", dbName))
	if err != nil {
		return err
	}

	fmt.Printf("Dropped database: %q\n", dbName)
	return nil
}

func getDatabaseName(dbURL string) (string, error) {
	url, err := url.Parse(dbURL)
	if err != nil {
		return "", err
	}

	dbName := url.Path
	if len(dbName) > 0 && dbName[:1] == "/" {
		dbName = dbName[1:]
	}

	return dbName, nil
}

func execOnPostgresDB(ctx context.Context, sqlString string) error {
	url, err := url.Parse(thousand.DatabaseURL)
	if err != nil {
		return err
	}

	url.Path = "postgres"

	db, err := pgx.Connect(ctx, url.String())
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sqlString)

	return err
}
