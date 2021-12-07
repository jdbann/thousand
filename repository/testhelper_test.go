package repository_test

import (
	"context"
	"os"
	"testing"

	"emailaddress.horse/thousand/repository"
	"github.com/jackc/pgx/v4"
)

type testRepository struct {
	*repository.Repository
	tx pgx.Tx
	t  *testing.T
}

func newTestRepository(t *testing.T) testRepository {
	// TODO: Derive DB connection from App config to prevent duplication - cannot
	// be done until config is extracted from app package due to circular
	// dependencies.
	var databaseURL = "postgres://localhost:5432/thousand_test?sslmode=disable"

	if os.Getenv("DATABASE_URL") != "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}

	repo, err := repository.New(repository.Options{
		DatabaseURL: databaseURL,
	})
	if err != nil {
		t.Fatal(err)
	}

	repo, tx, err := repo.WithTx(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := tx.Rollback(context.Background()); err != nil {
			t.Fatalf("Error attempting to rollback - DB may have unexpected contents: %s", err)
		}
	})

	return testRepository{repo, tx, t}
}

func (tc testRepository) WithSavepoint(query func(*repository.Repository) error) error {
	spRepo, spTx, err := tc.Repository.WithSavepoint(context.Background())
	if err != nil {
		return err
	}

	if err := query(spRepo); err != nil {
		if err := spTx.Rollback(context.Background()); err != nil {
			tc.t.Fatal(err)
		}

		return err
	}

	// Commit the savepoint TX to simulate individual transactions in the DB
	return spTx.Commit(context.Background())
}
