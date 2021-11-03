package browser

import (
	"context"
	"fmt"
	"testing"
)

func TestAppTitle(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	bt.Run(
		bt.Navigate("/"),
		bt.Text("h1").Equals("Thousand"),
	)
}

func TestShowVampires(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	if _, err := bt.Models().CreateVampire(context.Background(), "Gruffudd"); err != nil {
		t.Fatal(err)
	}

	bt.Run(
		bt.Navigate("/vampires"),
		bt.WaitVisible(`#vampires`),
		bt.Text(`#vampires`).Contains("Gruffudd"),
	)
}

func TestCreateVampire(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	newLinkSelector := `#newVampire a`
	nameFieldSelector := `#newVampire input[name="name"]`

	bt.Run(
		bt.Navigate("/vampires"),
		bt.WaitVisible(newLinkSelector),
		bt.Click(newLinkSelector),
		bt.WaitVisible(nameFieldSelector),
		bt.SendKeys(nameFieldSelector, "Gruffudd"),
		bt.Submit(nameFieldSelector),
		bt.Text(`#details`).Equals("Gruffudd"),
	)
}

func TestAddExperience(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	vampire, err := bt.Models().CreateVampire(context.Background(), "Gruffudd")
	if err != nil {
		t.Fatal(err)
	}

	newExperienceLinkSelector := `#memories a[href$="/experiences/new"]`
	experienceFieldSelector := `#memories input[name="description"]`
	expectedExperience := "I am Gruffudd, a Welsh farmer in the valleys of Pembroke; I am a recluse, fond of nature and withdrawn from the village."

	bt.Run(
		bt.Navigate(fmt.Sprintf("/vampires/%s", vampire.ID.String())),
		bt.WaitForTurbo(),
		bt.WaitVisible(newExperienceLinkSelector),
		bt.Click(newExperienceLinkSelector),
		bt.WaitVisible(experienceFieldSelector),
		bt.SendKeys(experienceFieldSelector, expectedExperience),
		bt.Submit(experienceFieldSelector),
		bt.Text(`#memories`).Contains(expectedExperience),
	)
}

func TestAddExperienceFormDismisses(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	vampire, err := bt.Models().CreateVampire(context.Background(), "Gruffudd")
	if err != nil {
		t.Fatal(err)
	}

	memory := vampire.Memories[0]

	newExperienceLinkSelector := fmt.Sprintf(`#memory-%s a[href$="/experiences/new"]`, memory.ID.String())
	experienceFieldSelector := fmt.Sprintf(`#memory-%s input[name="description"]`, memory.ID.String())

	bt.Run(
		bt.Navigate(fmt.Sprintf("/vampires/%s", vampire.ID.String())),
		bt.WaitForTurbo(),

		// Clicking outside the form should dismiss it
		bt.WaitVisible(newExperienceLinkSelector),
		bt.Click(newExperienceLinkSelector),
		bt.WaitVisible(experienceFieldSelector),
		bt.Click("h1"),
		bt.WaitNotPresent(experienceFieldSelector),

		// Clicking outside the form should not dismiss it if user has started
		// writing
		bt.WaitVisible(newExperienceLinkSelector),
		bt.Click(newExperienceLinkSelector),
		bt.WaitVisible(experienceFieldSelector),
		bt.SendKeys(experienceFieldSelector, "Sta"),
		bt.Click("h1"),
		bt.WaitVisible(experienceFieldSelector),

		// Clicking outside the form should dismiss it if user has cleared input
		bt.SendKeys(experienceFieldSelector, "\b\b\b"), // \b => Backspace
		bt.Click("h1"),
		bt.WaitNotPresent(experienceFieldSelector),
	)
}

func TestCannotAddFourExperiences(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	vampire, err := bt.Models().CreateVampire(context.Background(), "Gruffudd")
	if err != nil {
		t.Fatal(err)
	}

	memory := vampire.Memories[0]

	if _, err := bt.Models().AddExperience(context.Background(), vampire.ID, memory.ID, "Experience #1"); err != nil {
		t.Fatal(err)
	}

	if _, err := bt.Models().AddExperience(context.Background(), vampire.ID, memory.ID, "Experience #2"); err != nil {
		t.Fatal(err)
	}

	newExperienceLinkSelector := fmt.Sprintf(`#memory-%s a[href$="/experiences/new"]`, memory.ID.String())
	experienceFieldSelector := fmt.Sprintf(`#memory-%s input[name="description"]`, memory.ID.String())
	expectedExperience := "I am Gruffudd, a Welsh farmer in the valleys of Pembroke; I am a recluse, fond of nature and withdrawn from the village."

	bt.Run(
		bt.Navigate(fmt.Sprintf("/vampires/%s", vampire.ID.String())),
		bt.WaitForTurbo(),
		bt.WaitVisible(newExperienceLinkSelector),
		bt.Click(newExperienceLinkSelector),
		bt.WaitVisible(experienceFieldSelector),
		bt.SendKeys(experienceFieldSelector, expectedExperience),
		bt.Submit(experienceFieldSelector),
		bt.Text(fmt.Sprintf("#memory-%s", memory.ID.String())).Contains(expectedExperience),
		bt.Text(fmt.Sprintf("#memory-%s", memory.ID.String())).Not().Contains("New Experience"),
	)
}
