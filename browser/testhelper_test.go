package browser

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/jdbann/browsertest"

	"emailaddress.horse/thousand/app"
	"emailaddress.horse/thousand/repository"
)

var screenshotDir string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	screenshotDir = path.Join(path.Dir(filename), "..", "tmp", "screenshots")
}

type BrowserTest struct {
	browsertest.Test
	app  *app.App
	repo *repository.Repository
}

func NewBrowserTest(t *testing.T) *BrowserTest {
	databaseURL := "postgres://localhost:5432/thousand_test?sslmode=disable"
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
			t.Fatal(err)
		}
	})

	app := app.NewApp(app.Options{
		Repository: repo,
	})

	ts := httptest.NewServer(app)
	t.Cleanup(func() {
		ts.Close()
	})

	return &BrowserTest{
		browsertest.NewTest(t, ts.URL),
		app,
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
