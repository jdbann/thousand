package main

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/lib/pq"
	"github.com/urfave/cli/v2"
)

func createDatabase(c *cli.Context) error {
	dbName, err := getDatabaseName(thousand.DatabaseURL)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("CREATE database %s", pq.QuoteIdentifier(dbName))

	err = execOnPostgresDB(thousand.DatabaseURL, sql)
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

	sql := fmt.Sprintf("DROP database %s", pq.QuoteIdentifier(dbName))

	err = execOnPostgresDB(thousand.DatabaseURL, sql)
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

func execOnPostgresDB(dbURL, sqlString string) error {
	url, err := url.Parse(dbURL)
	if err != nil {
		return err
	}

	url.Path = "postgres"

	db, err := sql.Open("postgres", url.String())
	if err != nil {
		return err
	}

	_, err = db.Exec(sqlString)

	return err
}
