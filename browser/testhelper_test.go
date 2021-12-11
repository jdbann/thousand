package browser

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/jdbann/browsertest"

	"emailaddress.horse/thousand/repository"
	"emailaddress.horse/thousand/server"
	"emailaddress.horse/thousand/static"
)

var screenshotDir string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	screenshotDir = path.Join(path.Dir(filename), "..", "tmp", "screenshots")
}

type BrowserTest struct {
	browsertest.Test
	server *server.Server
	repo   *repository.Repository
}

func NewBrowserTest(t *testing.T) *BrowserTest {
	databaseURL := "postgres://localhost:5432/thousand_test?sslmode=disable"
	if os.Getenv("DATABASE_URL") != "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}

	logger := newLogger(t)

	repo, err := repository.New(repository.Options{
		DatabaseURL: databaseURL,
		Logger:      logger,
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
			t.Fatal(err)
		}
	})

	ts := newIntegrationTS()

	server := server.New(server.Options{
		Assets:     static.Assets,
		Logger:     logger,
		Mux:        ts.mux,
		Repository: repo,
		Server:     ts,
	})

	go func() {
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()

	t.Cleanup(func() {
		if err := server.Stop(); err != nil {
			panic(err)
		}
	})

	return &BrowserTest{
		browsertest.NewTest(t, ts.URL()),
		server,
		repo,
	}
}

func (bt *BrowserTest) Repository() *repository.Repository {
	return bt.repo
}

func (bt *BrowserTest) WaitForTurbo() browsertest.Action {
	return bt.Test.Poll(`window.Turbo != undefined`)
}

func (bt *BrowserTest) CaptureScreenshot(name string) browsertest.Action {
	var buf []byte

	filename := fmt.Sprintf("%s-%s.jpeg", bt.Name(), name)

	return bt.ActionFunc(
		chromedp.ActionFunc(func(ctx context.Context) error {
			if err := chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf)); err != nil {
				return err
			}

			bt.Cleanup(func() {
				if err := os.MkdirAll(screenshotDir, 0755); err != nil {
					bt.Fatal(err)
				}

				if err := os.WriteFile(path.Join(screenshotDir, filename), buf, 0644); err != nil {
					bt.Fatal(err)
				}

				bt.Logf("[CaptureScreenshot] Captured %s\n", filename)
			})

			return nil
		}),
		fmt.Sprintf("[CaptureScreenshot] Capturing %s", filename),
	)
}
