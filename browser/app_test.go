package browser

import (
	"context"
	"fmt"
	"testing"

	"emailaddress.horse/thousand/form"
)

func TestVampireFlow(t *testing.T) {
	bt := NewBrowserTest(t)

	bt.Run(
		// Register
		bt.Navigate("/"),
		bt.Click(`#newSession a[href="/user/new"]`),
		bt.Text("header .brand").Equals("Thousand"),
		bt.SendKeys(`#newUser input[name="email"]`, "john@bannister.com"),
		bt.SendKeys(`#newUser input[name="password"]`, "password"),
		bt.Submit(`#newUser button[type="submit"]`),
		bt.Text(`#flashes`).Equals("Thank you for signing up!"),

		// Create vampire
		bt.Click(`#newVampire a`),
		bt.WaitVisible(`#newVampire input[name="name"]`),
		bt.SendKeys(`#newVampire input[name="name"]`, "Gruffudd"),
		bt.Submit(`#newVampire input[name="name"]`),
		bt.Text(`#details`).Equals("Gruffudd"),

		// Add experience
		bt.Click(`#memories a`),
		bt.WaitVisible(`#memories input[name="description"]`),
		bt.SendKeys(`#memories input[name="description"]`, "I am Gruffudd, a Welsh farmer in the valleys of Pembroke; I am a recluse, fond of nature and withdrawn from the village."),
		bt.Submit(`#memories input[name="description"]`),
		bt.Text(`#memories`).Contains("I am Gruffudd, a Welsh farmer in the valleys of Pembroke; I am a recluse, fond of nature and withdrawn from the village."),

		// Add skill
		bt.Click(`#skills a`),
		bt.WaitVisible(`#skills input[name="description"]`),
		bt.SendKeys(`#skills input[name="description"]`, "Navigating forests"),
		bt.Submit(`#skills input[name="description"]`),
		bt.Text(`#skills`).Contains("Navigating forests"),

		// Add resources
		bt.Click(`#resources a`),
		bt.WaitVisible(`#resources input[name="description"]`),
		bt.SendKeys(`#resources input[name="description"]`, "Calweddyn Farm, rich but challenging soils"),
		bt.Click(`#resources input[name="stationary"]`),
		bt.Submit(`#resources input[name="description"]`),
		bt.Text(`#resources`).Contains("Calweddyn Farm, rich but challenging soils"),
		bt.Text(`#resources`).Contains("Stationary"),

		// Add character
		bt.Click(`#characters a`),
		bt.WaitVisible(`#characters input[name="name"]`),
		bt.SendKeys(`#characters input[name="name"]`, "Lord Othian, English gentry visiting a cathedral in St. Davids."),
		bt.SendKeys(`#characters select[name="type"]`, "Immortal"),
		bt.Submit(`#characters input[name="name"]`),
		bt.Text(`#characters`).Contains("Lord Othian, English gentry visiting a cathedral in St. Davids."),
		bt.Text(`#characters`).Contains("Immortal"),

		// Add mark
		bt.Click(`#marks a`),
		bt.WaitVisible(`#marks input[name="description"]`),
		bt.SendKeys(`#marks input[name="description"]`, "Muddy footprints, muddy handprints, muddy sheets - I leave a trail of dirt wherever I travel."),
		bt.Submit(`#marks input[name="description"]`),
		bt.Text(`#marks`).Contains("Muddy footprints, muddy handprints, muddy sheets - I leave a trail of dirt wherever I travel."),
	)
}

func TestLogInAndOut(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	_, err := bt.Repository().CreateUser(context.Background(), form.NewUser("john@bannister.com", "password"))
	if err != nil {
		t.Fatal(err)
	}

	bt.Run(
		bt.Navigate("/session/new"),
		bt.SendKeys(`#newSession input[name="email"]`, "john@bannister.com"),
		bt.SendKeys(`#newSession input[name="password"]`, "password"),
		bt.Submit(`#newSession button[type="submit"]`),
		bt.Text(`#flashes`).Equals("Welcome back!"),

		bt.Click(`#destroySession button[type="submit"]`),
		bt.WaitNotPresent(`#destroySession`),
		bt.Text(`h1`).Equals("Log in"),
	)
}

func TestAddExperienceFormDismisses(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	user, err := bt.Repository().CreateUser(context.Background(), form.NewUser("john@bannister.com", "password"))
	if err != nil {
		t.Fatal(err)
	}

	vampire, err := bt.Repository().CreateVampire(context.Background(), user.ID, "Gruffudd")
	if err != nil {
		t.Fatal(err)
	}

	memory := vampire.Memories[0]

	newExperienceLinkSelector := fmt.Sprintf(`#memory-%s a[href$="/experiences/new"]`, memory.ID.String())
	experienceFieldSelector := fmt.Sprintf(`#memory-%s input[name="description"]`, memory.ID.String())

	bt.Run(
		bt.AuthenticateAs("john@bannister.com", "password"),
		bt.Navigate(fmt.Sprintf("/vampires/%s", vampire.ID.String())),

		// Clicking outside the form should dismiss it
		bt.WaitVisible(newExperienceLinkSelector),
		bt.Click(newExperienceLinkSelector),
		bt.WaitVisible(experienceFieldSelector),
		bt.Click("body"),
		bt.WaitNotPresent(experienceFieldSelector),

		// Clicking outside the form should not dismiss it if user has started
		// writing
		bt.WaitVisible(newExperienceLinkSelector),
		bt.Click(newExperienceLinkSelector),
		bt.WaitVisible(experienceFieldSelector),
		bt.SendKeys(experienceFieldSelector, "Sta"),
		bt.Click("body"),
		bt.WaitVisible(experienceFieldSelector),

		// Clicking outside the form should dismiss it if user has cleared input
		bt.SendKeys(experienceFieldSelector, "\b\b\b"), // \b => Backspace
		bt.Click("body"),
		bt.WaitNotPresent(experienceFieldSelector),
	)
}

func TestCannotAddFourExperiences(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	user, err := bt.Repository().CreateUser(context.Background(), form.NewUser("john@bannister.com", "password"))
	if err != nil {
		t.Fatal(err)
	}

	vampire, err := bt.Repository().CreateVampire(context.Background(), user.ID, "Gruffudd")
	if err != nil {
		t.Fatal(err)
	}

	memory := vampire.Memories[0]

	if _, err := bt.Repository().CreateExperience(context.Background(), vampire.ID, memory.ID, "Experience #1"); err != nil {
		t.Fatal(err)
	}

	if _, err := bt.Repository().CreateExperience(context.Background(), vampire.ID, memory.ID, "Experience #2"); err != nil {
		t.Fatal(err)
	}

	newExperienceLinkSelector := fmt.Sprintf(`#memory-%s a[href$="/experiences/new"]`, memory.ID.String())
	experienceFieldSelector := fmt.Sprintf(`#memory-%s input[name="description"]`, memory.ID.String())
	expectedExperience := "I am Gruffudd, a Welsh farmer in the valleys of Pembroke; I am a recluse, fond of nature and withdrawn from the village."

	bt.Run(
		bt.AuthenticateAs("john@bannister.com", "password"),
		bt.Navigate(fmt.Sprintf("/vampires/%s", vampire.ID.String())),

		bt.WaitVisible(newExperienceLinkSelector),
		bt.Click(newExperienceLinkSelector),
		bt.WaitVisible(experienceFieldSelector),
		bt.SendKeys(experienceFieldSelector, expectedExperience),
		bt.Submit(experienceFieldSelector),
		bt.Text(fmt.Sprintf("#memory-%s", memory.ID.String())).Contains(expectedExperience),
		bt.Text(fmt.Sprintf("#memory-%s", memory.ID.String())).Not().Contains("New Experience"),
	)
}
