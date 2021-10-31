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
	"emailaddress.horse/thousand/app/models"
)

var screenshotDir string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	screenshotDir = path.Join(path.Dir(filename), "..", "tmp", "screenshots")
}

type BrowserTest struct {
	browsertest.Test
	app *app.App
}

func NewBrowserTest(t *testing.T) *BrowserTest {
	app := app.NewApp(app.TestConfig(t))
	ts := httptest.NewServer(app)
	t.Cleanup(func() {
		ts.Close()
	})

	return &BrowserTest{
		browsertest.NewTest(t, ts.URL),
		app,
	}
}

func (bt *BrowserTest) Models() *models.Models {
	return bt.app.Models
}

func (bt *BrowserTest) WaitForTurbo() browsertest.Action {
	return bt.ActionFunc(
		chromedp.ActionFunc(func(ctx context.Context) error {
			var result bool

			if err := chromedp.Run(ctx, chromedp.Poll(`window.Turbo != undefined`, &result)); err != nil {
				return err
			}

			if result == false {
				// I don't think this should be possible
				bt.Fatal("[WaitForTurbo] Attempt to wait returned false")
			}

			return nil
		}),
		"[WaitForTurbo]",
	)
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
