package app

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func TxnTestConfig(t *testing.T) EnvConfigurer {
	return func(app *App) {
		TestConfig(t)(app)

		appDBTX, err := app.DB()
		if err != nil {
			t.Fatal(err)
		}

		sqlDB, ok := appDBTX.(*sql.DB)
		if !ok {
			t.Fatal(errors.New("app._db is not a *sql.DB"))
		}

		tx, err := sqlDB.BeginTx(context.Background(), &sql.TxOptions{})
		if err != nil {
			t.Fatal(err)
		}
		app._db = tx

		t.Cleanup(func() {
			if err := tx.Rollback(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
