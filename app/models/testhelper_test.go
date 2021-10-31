package models

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jackc/pgx/v4"
)

var (
	// Memories when blank have no matchable attributes so need to be ignored
	ignoreWholeVampireFields = cmpopts.IgnoreFields(WholeVampire{}, "Memories")

	// ID and timestamps are defined by the DB and do not require matching
	ignoreNewVampireFields = cmpopts.IgnoreFields(NewVampire{}, "ID", "CreatedAt", "UpdatedAt")

	// ID, MemoryID and timestamps are defined by the DB and do not require matching
	ignoreNewExperienceFields = cmpopts.IgnoreFields(NewExperience{}, "ID", "MemoryID", "CreatedAt", "UpdatedAt")
)

type testModels struct {
	*Models
	tx pgx.Tx
	t  *testing.T
}

func newTestModels(t *testing.T) testModels {
	// TODO: Derive DB connection from App config to prevent duplication - cannot
	// be done until config is extracted from app package due to circular
	// dependencies.
	var databaseURL = "postgres://localhost:5432/thousand_test?sslmode=disable"

	if os.Getenv("DATABASE_URL") != "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}

	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := tx.Rollback(context.Background()); err != nil {
			t.Fatalf("Error attempting to rollback - DB may have unexpected contents: %s", err)
		}
	})

	return testModels{NewModels(tx), tx, t}
}

func (tc testModels) WithSavepoint(query func(*Models) error) error {
	savepointTx, err := tc.tx.Begin(context.Background())
	if err != nil {
		return err
	}

	if err := query(NewModels(savepointTx)); err != nil {
		if err := savepointTx.Rollback(context.Background()); err != nil {
			tc.t.Fatal(err)
		}

		return err
	}

	// Commit the savepoint TX to simulate individual transactions in the DB
	return savepointTx.Commit(context.Background())
}
