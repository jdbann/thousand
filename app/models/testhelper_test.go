package models

import (
	"context"
	"database/sql"
	"os"
	"testing"
)

func newTestModels(t *testing.T) *Models {
	// TODO: Derive DB connection from App config to prevent duplication - cannot
	// be done until config is extracted from app package due to circular
	// dependencies.
	var databaseURL = "postgres://localhost:5432/thousand_test?sslmode=disable"

	if os.Getenv("DATABASE_URL") != "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}

	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	txn, err := conn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := txn.Rollback(); err != nil {
			t.Fatal(err)
		}
	})

	return NewModels(txn)
}
