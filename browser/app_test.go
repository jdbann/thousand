package browser

import (
	"strings"
	"testing"
)

func TestAppTitle(t *testing.T) {
	bt := NewBrowserTest(t)

	var title string

	bt.Run(
		bt.Navigate("/"),
		bt.Text("h1", &title),
	)

	if title != "Thousand" {
		t.Errorf("expected %q; got %q", "Thousand", title)
	}
}

func TestSetName(t *testing.T) {
	bt := NewBrowserTest(t)

	var name string

	bt.Run(
		bt.Navigate("/"),
		bt.WaitVisible(`input[name="name"]`),
		bt.SendKeys(`input[name="name"]`, "Gruffudd"),
		bt.Submit(`input[name="name"]`),
		bt.Text(`#details`, &name),
	)

	if strings.TrimSpace(name) != "Gruffudd" {
		t.Errorf("expected %q; got %q", "Gruffudd", name)
	}
}
