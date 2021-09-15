package browser

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"emailaddress.horse/thousand/app"
	"github.com/chromedp/chromedp"
)

type BrowserTest struct {
	*testing.T
	ctx     context.Context
	timeout time.Duration
	baseURL string
	app     *app.App
}

type BrowserAction struct {
	chromedp.Action
	msg string
}

func NewBrowserTest(t *testing.T) *BrowserTest {
	app := app.NewApp(app.TxnTestConfig(t))
	ts := httptest.NewServer(app)
	t.Cleanup(func() {
		ts.Close()
	})

	return &BrowserTest{
		t,
		context.Background(),
		time.Second * 2,
		ts.URL,
		app,
	}
}

func (bt *BrowserTest) Run(actions ...BrowserAction) {
	ctx, cancel := chromedp.NewContext(bt.ctx)
	defer cancel()

	bt.executeAction(ctx, actions[0])

	for _, action := range actions[1:] {
		ctx, cancel := context.WithTimeout(ctx, bt.timeout)
		defer cancel()

		bt.executeAction(ctx, action)
	}
}

func (bt *BrowserTest) executeAction(ctx context.Context, action BrowserAction) {
	bt.Log(action.msg)
	if err := chromedp.Run(ctx, action); err != nil {
		bt.Fatalf("%s: %q", action.msg, err)
	}
}

func (bt *BrowserTest) Navigate(url string) BrowserAction {
	return BrowserAction{
		chromedp.Navigate(bt.baseURL + url),
		fmt.Sprintf("[Navigate] %q", url),
	}
}

func (bt *BrowserTest) Click(sel interface{}, opts ...chromedp.QueryOption) BrowserAction {
	return BrowserAction{
		chromedp.Click(sel, opts...),
		fmt.Sprintf("[Click] %v", sel),
	}
}

func (bt *BrowserTest) Text(sel interface{}, text *string, opts ...chromedp.QueryOption) BrowserAction {
	return BrowserAction{
		chromedp.Text(sel, text, opts...),
		fmt.Sprintf("[Text] %v", sel),
	}
}

func (bt *BrowserTest) InnerHTML(sel interface{}, text *string, opts ...chromedp.QueryOption) BrowserAction {
	return BrowserAction{
		chromedp.InnerHTML(sel, text, opts...),
		fmt.Sprintf("[InnerHTML] %v", sel),
	}
}

func (bt *BrowserTest) WaitVisible(sel interface{}, opts ...chromedp.QueryOption) BrowserAction {
	return BrowserAction{
		chromedp.WaitVisible(sel, opts...),
		fmt.Sprintf("[WaitVisible] %v", sel),
	}
}

func (bt *BrowserTest) WaitReady(sel interface{}, opts ...chromedp.QueryOption) BrowserAction {
	return BrowserAction{
		chromedp.WaitReady(sel, opts...),
		fmt.Sprintf("[WaitReady] %v", sel),
	}
}

func (bt *BrowserTest) SendKeys(sel interface{}, v string, opts ...chromedp.QueryOption) BrowserAction {
	return BrowserAction{
		chromedp.SendKeys(sel, v, opts...),
		fmt.Sprintf("[SendKeys] %v %q", sel, v),
	}
}

func (bt *BrowserTest) Submit(sel interface{}, opts ...chromedp.QueryOption) BrowserAction {
	return BrowserAction{
		chromedp.Submit(sel, opts...),
		fmt.Sprintf("[Submit] %v", sel),
	}
}
