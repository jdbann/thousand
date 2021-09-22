package app

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	"emailaddress.horse/thousand/app/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TestConfig sets up the app for running tests in a test environment by
// running all DB interactions in a transaction to prevent tests impacting on
// each other.
func TestConfig(t *testing.T) EnvConfigurer {
	return func(app *App) {
		// Apply base config
		BaseTestConfigurer(t.Logf)(app)

		// App configuration values
		// Manually check for DATABASE_URL as tests cannot infer config from flags
		// and environment variables through the CLI.
		if os.Getenv("DATABASE_URL") != "" {
			app.DatabaseURL = os.Getenv("DATABASE_URL")
		}
		app.DBConnector = _transactionConnector(t)
	}
}

// LiveTestConfig sets up the app for a test environment with an adapter for the
// usual (testing.T).Log and (testing.T).Logf which sends them to the app's
// default Logger.
var LiveTestConfig Configurer = EnvConfigurer(liveTestConfig)

func liveTestConfig(app *App) {
	// Apply base config
	BaseTestConfigurer(app.Logger.Debugf)(app)
}

// BaseTestConfigurer sets up the app for a test environment.
func BaseTestConfigurer(outf func(string, ...interface{})) EnvConfigurer {
	return func(app *App) {
		// Echo configuraton values
		app.Debug = true

		// App configuration values
		app.DatabaseURL = "postgres://localhost:5432/thousand_test?sslmode=disable"

		// Injected middleware
		app.LoggerMiddleware = _readableLogger(outf)
		app.HTTPErrorHandler = _logError(outf, app.HTTPErrorHandler)
	}
}

func _transactionConnector(t *testing.T) DBConnector {
	return func(databaseURL string) (models.DBTX, error) {
		conn, err := sql.Open("postgres", databaseURL)
		if err != nil {
			return nil, err
		}

		txn, err := conn.BeginTx(context.Background(), &sql.TxOptions{})
		if err != nil {
			return nil, err
		}

		t.Cleanup(func() {
			if err := txn.Rollback(); err != nil {
				t.Fatal(err)
			}
		})

		return txn, nil
	}
}

func _readableLogger(outf func(string, ...interface{})) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Path(), "/assets/")
		},

		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			outf("%d %-8s %s", values.Status, values.Method, values.URI)
			return nil
		},
		LogMethod: true,
		LogURI:    true,
		LogStatus: true,
	})
}

func _logError(outf func(string, ...interface{}), handler echo.HTTPErrorHandler) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		outf("%w", err)

		handler(err, c)
	}
}
