package browser

import (
	"context"
	"net/http/httptest"
	"testing"

	"emailaddress.horse/thousand/app"
	"github.com/chromedp/chromedp"
)

func TestAppTitle(t *testing.T) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	app := &app.App{}
	ts := httptest.NewServer(app.Routes())
	defer ts.Close()

	var title string
	err := chromedp.Run(ctx,
		chromedp.Navigate(ts.URL),
		chromedp.Text("h1", &title, chromedp.ByQuery),
	)
	if err != nil {
		t.Fatal(err)
	}

	if title != "Thousand" {
		t.Errorf("expected %q; got %q", "Thousand", title)
	}
}
