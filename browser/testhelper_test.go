package browser

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/jdbann/browsertest"
	"go.uber.org/fx"

	"emailaddress.horse/thousand/form"
	"emailaddress.horse/thousand/repository"
	"emailaddress.horse/thousand/server"
)

var screenshotDir string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	screenshotDir = path.Join(path.Dir(filename), "..", "tmp", "screenshots")
}

type BrowserTest struct {
	browsertest.Test
	repo *repository.Repository
}

func NewBrowserTest(t *testing.T) *BrowserTest {
	databaseURL := "postgres://localhost:5432/thousand_test?sslmode=disable"
	if url, ok := os.LookupEnv("DATABASE_URL"); ok {
		databaseURL = url
	}

	var s *server.Server
	var repo *repository.Repository

	a := fx.New(
		fx.Supply(params{
			DatabaseURL: databaseURL,
			Port:        0,
			SecretKey:   "secret",
			T:           t,
		}),

		testModule,

		fx.Populate(&s, &repo),
	)

	go func() {
		if err := a.Start(context.Background()); err != nil {
			panic(err)
		}
	}()

	t.Cleanup(func() {
		if err := a.Stop(context.Background()); err != nil {
			panic(err)
		}
	})

	url := func() string {
		wait := time.NewTimer(500 * time.Millisecond)
		for {
			select {
			case <-wait.C:
				panic("error getting address for integration test server")
			default:
				if s != nil && s.URL() != "" {
					return "http://" + s.URL()
				}
			}
		}
	}()

	return &BrowserTest{
		browsertest.NewTest(t, url),
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

func (bt *BrowserTest) Authenticate() browsertest.Action {
	return bt.ActionFunc(
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := bt.Repository().CreateUser(context.Background(), form.NewUser("john@bannister.com", "password"))
			if err != nil {
				bt.Fatal(err)
			}

			return chromedp.Run(ctx, chromedp.Tasks{
				bt.Navigate("/session/new"),
				bt.SendKeys(`#newSession input[name="email"]`, "john@bannister.com"),
				bt.SendKeys(`#newSession input[name="password"]`, "password"),
				bt.Submit(`#newSession button[type="submit"]`),
				bt.Text(`#flashes`).Equals("Welcome back!"),
			})
		}),
		"[Authenticate]",
	)
}
