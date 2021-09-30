package browser

import (
	"net/http/httptest"
	"testing"

	"github.com/jdbann/browsertest"

	"emailaddress.horse/thousand/app"
	"emailaddress.horse/thousand/app/models"
)

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
